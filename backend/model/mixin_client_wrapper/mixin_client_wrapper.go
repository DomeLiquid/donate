package mixin_client_wrapper

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"donate/config"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/fox-one/mixin-sdk-go/v2/mixinnet"

	"github.com/lixvyang/go-utils/cacheflight"

	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

var (
	ErrConfigNil = errors.New("config is nil")
)

type MixinClientWrapper struct {
	*mixin.Client
	User *mixin.User

	transferMutex             sync.Mutex
	SpendKey                  mixinnet.Key
	userMixinAssetAmountCache *cacheflight.Group
	rateLimiter               *rate.Limiter
}

func NewMixinClientWrapper(config *config.MixinConfig) (*MixinClientWrapper, error) {
	if config == nil {
		return nil, ErrConfigNil
	}

	client, err := mixin.NewFromKeystore(&mixin.Keystore{
		SessionID:         config.SessionID,
		ServerPublicKey:   config.ServerPublicKey,
		ClientID:          config.ClientID,
		PrivateKey:        config.PrivateKey,
		PinToken:          config.PinToken,
		AppID:             config.AppID,
		SessionPrivateKey: config.SessionPrivateKey,
	})
	if err != nil {
		return nil, err
	}
	user, err := client.UserMe(context.Background())
	if err != nil {
		return nil, err
	}
	spendKeyStr := config.SpendKey
	if spendKeyStr == "" {
		spendKeyStr = os.Getenv("SPEND_KEY")
	}

	spendKey, err := mixinnet.ParseKeyWithPub(spendKeyStr, user.SpendPublicKey)
	if err != nil {
		return nil, err
	}

	return &MixinClientWrapper{
		Client:                    client,
		User:                      user,
		SpendKey:                  spendKey,
		transferMutex:             sync.Mutex{},
		userMixinAssetAmountCache: cacheflight.New(mixinAssetAmountCacheTTL, mixinAssetAmountCacheDelay),
		rateLimiter:               rate.NewLimiter(rate.Every(time.Second), 20),
	}, nil
}

func (m *MixinClientWrapper) CreateSubbot(ctx context.Context) error {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}
	sessionPrivateKey := hex.EncodeToString(privateKey)
	sessionPublicKey := hex.EncodeToString(publicKey)

	_, keystore, err := m.Client.CreateUser(ctx, privateKey, "Inscription Mgr")
	if err != nil {
		log.Error().Err(err).Msg("create subbot failed")
		return err
	}
	subClient, err := mixin.NewFromKeystore(keystore)
	if err != nil {
		log.Error().Err(err).Msg("create subbot failed")
		return err
	}

	pin := mixinnet.GenerateKey(rand.Reader)
	err = subClient.ModifyPin(context.TODO(), "", pin.Public().String())
	if err != nil {
		log.Error().Err(err).Msg("modify pin failed")
		return err
	}
	spendKey := mixinnet.GenerateKey(rand.Reader)
	user, err := subClient.SafeMigrate(ctx, spendKey.String(), pin.String())
	if err != nil {
		log.Error().Err(err).Msg("safe migrate failed")
		return err
	}

	log.Info().Msg("create subbot success")
	bts, _ := json.Marshal(user)
	keystoreJson, _ := json.Marshal(keystore)
	log.Info().
		Any("session private key", sessionPrivateKey).
		Any("session public key", sessionPublicKey).
		Any("keystore", keystoreJson).
		Any("user info", bts).
		Any("pin", pin.Public().String()).
		Any("spend public key", spendKey.Public().String()).
		Any("spend private key", spendKey.String()).
		Msg("subbot info")

	return nil
}
