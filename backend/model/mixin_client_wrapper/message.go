package mixin_client_wrapper

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/rand"
)

func (m *MixinClientWrapper) SendCardWithRetry(ctx context.Context, receiptId string, card *mixin.AppCardMessage) (err error) {
	card.AppID = m.ClientID

	cardBytes, err := json.Marshal(card)
	if err != nil {
		return
	}

	cardBase64code := base64.StdEncoding.EncodeToString(cardBytes)
	messageRequest := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(m.ClientID, receiptId),
		RecipientID:    receiptId,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryAppCard,
		Data:           cardBase64code,
	}

	return m.sendMessageV2(ctx, messageRequest)
}

func (m *MixinClientWrapper) SendMessageWithRetry(ctx context.Context, receiptId string, text string) (err error) {
	for i := 0; i < defaultMaxMixinRetry; i++ {
		if err = m.sendMessage(ctx, receiptId, text); err != nil {
			log.Error().Err(err).Msg("send message failed, retrying...")
			time.Sleep(time.Second << i)
			continue
		} else {
			return nil
		}
	}
	return err
}

func (m *MixinClientWrapper) sendMessage(ctx context.Context, receiptId string, text string) (err error) {
	req := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(m.Client.ClientID, receiptId),
		RecipientID:    receiptId,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryPlainText,
		Data:           base64.StdEncoding.EncodeToString([]byte(text)),
	}

	err = m.Client.SendMessage(ctx, req)
	if err != nil {
		// try create conversation
		_, err = m.Client.CreateContactConversation(ctx, req.RecipientID)
		if err != nil {
			return err
		}
		err = m.Client.SendMessage(ctx, req)
		if err != nil {
			return err
		}

		return err
	}
	return nil
}

func (m *MixinClientWrapper) sendMessageV2(ctx context.Context, message *mixin.MessageRequest) error {
	var baseDelay = time.Second
	var maxRetries = 3
	sendMessageToUser := func() error {
		err := m.Client.SendMessage(ctx, message)
		if err != nil {
			// try create conversation
			_, err = m.Client.CreateContactConversation(ctx, message.RecipientID)
			if err != nil {
				return err
			}
			err = m.Client.SendMessage(ctx, message)
			if err != nil {
				return err
			}
		}
		return nil
	}

	var err error
	for retry := 0; retry < maxRetries; retry++ {
		err = sendMessageToUser()
		if err == nil {
			break
		}
		delay := baseDelay * time.Duration(1<<retry)
		time.Sleep(delay)
	}

	return err
}

var cardColorList = []string{
	"#7983C2", "#8F7AC5", "#C5595A", "#C97B46", "#76A048", "#3D98D0",
	"#5979F0", "#8A64D0", "#B76753", "#AA8A46", "#9CAD23", "#6BC0CE",
	"#6C89D3", "#AA66C3", "#C8697D", "#C49B4B", "#5FB05F", "#52A98B",
	"#75A2CB", "#A75C96", "#9B6D77", "#A49373", "#6AB48F", "#93B289",
}

func RandomCardColor() string {
	return cardColorList[rand.Intn(len(cardColorList))]
}

// // TODO ENCRYPTED MESSAGE
// func (m *MixinClientWrapper) SendEncryptedMessagesWithRetry(ctx context.Context, receiptIds []string, text string) (err error) {
// 	for i := 0; i < defaultMaxMixinRetry; i++ {
// 		if err = m.SendMessages(ctx, receiptIds, text); err != nil {
// 			time.Sleep(time.Second << i)
// 			continue
// 		} else {
// 			return nil
// 		}
// 	}
// 	return err
// }

// func (m *MixinClientWrapper) SendMessages(ctx context.Context, receiptId []string, text string) (err error) {
// 	var req []*mixin.MessageRequest
// 	for _, r := range receiptId {
// 		req = append(req, &mixin.MessageRequest{
// 			ConversationID: mixin.UniqueConversationID(m.Client.ClientID, r),
// 			RecipientID:    r,
// 			MessageID:      mixin.RandomTraceID(),
// 			Category:       mixin.MessageCategoryPlainText,
// 			Data:           base64.StdEncoding.EncodeToString([]byte(text)),
// 		})
// 	}

// 	err = m.Client.SendMessages(ctx, req)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
