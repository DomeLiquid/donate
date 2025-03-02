package router

import (
	"context"
	"donate/model"
	"donate/model/mixin_client_wrapper"
	"donate/pkg/thread"
	"donate/router/middleware"
	"donate/utils"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	ErrInternalServerError = errors.New("internal server error")

	ErrReadMixinSafeSnapshotFailed = errors.New("read mixin safe snapshot failed")
	ErrGetSnapshotCountFailed      = errors.New("get snapshot count failed")
	ErrGetLastestSnapshotFailed    = errors.New("get lastest snapshot failed")
	ErrGetSnapshotByIdFailed       = errors.New("get snapshot by id failed")
)

type sortSnapshot []*mixin.SafeSnapshot

func (s sortSnapshot) Len() int {
	return len(s)
}

func (s sortSnapshot) Less(i, j int) bool {
	return s[i].CreatedAt.Before(s[j].CreatedAt)
}

func (s sortSnapshot) Swap(i, j int) {
	if s[i] == nil || s[j] == nil {
		return
	}
	s[i], s[j] = s[j], s[i]
}

func (s *Service) RunMixinLoop(ctx context.Context) {
	thread.GoSafe(func() {
		ticker := s.clock.Ticker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Error().Err(ctx.Err()).Msg("cancel cron server")
				return
			case <-ticker.C:
				err := s.handleMixinSnapshotInput(ctx)
				if err != nil {
					log.Error().Err(err).Msg("cron handle snapshot input failed")
				}
			}
		}
	})
}

func (s *Service) handleMixinSnapshotInput(ctx context.Context) error {
	now := s.clock.Now()

	snapshots, err := s.mixinClient.ReadSafeSnapshots(ctx, "", now.Add(-time.Hour), "DESC", 500)
	if err != nil {
		return ErrReadMixinSafeSnapshotFailed
	}

	count, err := s.store.GetSnapshotCount(ctx)
	if err != nil {
		return ErrGetSnapshotCountFailed
	}
	sort.Sort(sortSnapshot(snapshots))

	startIndex := 0
	if count > 0 {
		lastestSnapshot, err := s.store.GetLastestSnapshot(ctx)
		if err != nil {
			return ErrGetLastestSnapshotFailed
		}

		for i, snapshot := range snapshots {
			if snapshot.RequestID == lastestSnapshot.RequestId {
				startIndex = i + 1 // 从最新快照的下一个开始处理
				break
			}
		}
	}

	for i := startIndex; i < len(snapshots); i++ {
		snapshot := snapshots[i]
		if !snapshot.Amount.IsPositive() {
			continue
		}

		// 聚合 utxo 和 memo 为空 忽略
		if snapshot.Memo == "" {
			continue
		}

		err = s.handleMixinInput(ctx, snapshot)
		if err != nil {
			log.Error().Any("snapshot", snapshot).Err(err).Msg("handle mixin input failed")
			continue
		}
	}

	return nil
}

func (s *Service) handleMixinInput(ctx context.Context, snapshot *mixin.SafeSnapshot) (err error) {
	logger := log.Logger.With().Str(middleware.DefaultXid, middleware.GenReqId()).Logger()

	err = s.store.InsertSnapshot(ctx, &model.Snapshot{
		SnapshotId: snapshot.SnapshotID,
		RequestId:  snapshot.RequestID,
		UserId:     snapshot.OpponentID,
		AssetId:    snapshot.AssetID,
		Memo:       snapshot.Memo,
		CreatedAt:  snapshot.CreatedAt.Unix(),
		Amount:     snapshot.Amount,
	})
	if err != nil {
		logger.Error().Err(err).Msg("insert snapshot failed")
		return err
	}

	// 解析meme 获取 pid,然后将资产转给 pid 对应的用户
	// 解析 meme 失败则退还给用户
	snapshotMemoHex, err := hex.DecodeString(snapshot.Memo)
	if err != nil {
		return err
	}

	pid, err := uuid.FromString(string(snapshotMemoHex))
	if err != nil || pid == uuid.Nil {
		return errors.New("invalid pid")
	}

	project, err := s.store.GetProject(ctx, pid.String())
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	refundToUser := func() error {
		err = s.mixinClient.TransferOneWithRetry(ctx, &mixin_client_wrapper.TransferOneRequest{
			RequestId: utils.GenUuidFromStrings(snapshot.RequestID, "donate-refund"),
			AssetId:   snapshot.AssetID,
			Amount:    snapshot.Amount,
			Member:    snapshot.OpponentID,
			Memo:      fmt.Sprintf("Donate failed"),
		})
		if err != nil {
			return err
		}
		return nil
	}

	if err == gorm.ErrRecordNotFound {
		logger.Error().Err(err).Msg("project not found")
		return refundToUser()
	}

	recipientUser, err := s.mixinClient.ReadUser(ctx, snapshot.OpponentID)
	if err != nil {
		logger.Error().Err(err).Msg("read user failed")
		return refundToUser()
	}

	// 查找用户,如果没有则写入 有就更新

	// 检查用户在本地库是否存在,存在更新,不存在创建
	duser, err := s.store.GetUserByIdentityNumber(ctx, recipientUser.IdentityNumber)
	switch err {
	case nil:
		// 更新
		duser = &model.User{
			IdentityNumber: recipientUser.IdentityNumber,
			FullName:       recipientUser.FullName,
			MixinUID:       recipientUser.UserID,
			AvatarUrl:      recipientUser.AvatarURL,
			Biography:      recipientUser.Biography,
			MixinCreatedAt: recipientUser.CreatedAt,
			UpdatedAt:      time.Now(),
		}
		if err = s.store.UpdateUserBymuid(ctx, recipientUser.UserID, duser); err != nil {
			logger.Error().Err(err).Msg("failed to update user")
		}
	case gorm.ErrRecordNotFound:
		// 创建
		duser = &model.User{
			IdentityNumber: recipientUser.IdentityNumber,
			FullName:       recipientUser.FullName,
			MixinUID:       recipientUser.UserID,
			AvatarUrl:      recipientUser.AvatarURL,
			Biography:      recipientUser.Biography,
			MixinCreatedAt: recipientUser.CreatedAt,
			CreatedAt:      time.Now(),
		}
		if err := s.store.AddUser(ctx, duser); err != nil {
			logger.Error().Err(err).Msg("failed to create user")
		}
	default:
		logger.Error().Err(err).Msg("failed to get user")
	}

	// donate cnt ++
	_ = s.store.DonateActionStore.AddDonateAction(ctx, &model.DonateAction{
		ID:             utils.GenUuidFromStrings(snapshot.RequestID, "donate"),
		PID:            pid.String(),
		Amount:         snapshot.Amount,
		IdentityNumber: recipientUser.IdentityNumber,
		AssetID:        snapshot.AssetID,
		CreatedAt:      snapshot.CreatedAt,
	})

	_ = s.store.ProjectStore.IncrProjectDonateCnt(ctx, pid.String())

	donateMsg := fmt.Sprintf("User %s has donated to you for project %s.",
		recipientUser.IdentityNumber,
		project.PID)
	err = s.mixinClient.SendMessageWithRetry(ctx, project.MixinUID, donateMsg)
	if err != nil {
		logger.Error().Err(err).Msg("send donate msg error")
	}

	// 转给 pid 对应的用户
	err = s.mixinClient.TransferOneWithRetry(ctx, &mixin_client_wrapper.TransferOneRequest{
		RequestId: utils.GenUuidFromStrings(snapshot.RequestID, "donate-transfer"),
		AssetId:   snapshot.AssetID,
		Amount:    snapshot.Amount,
		Member:    project.MixinUID,
		Memo:      fmt.Sprintf("Donate for you"),
	})
	if err != nil {
		return err
	}

	return nil
}
