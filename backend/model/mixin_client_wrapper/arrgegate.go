package mixin_client_wrapper

import (
	"context"
	"donate/utils"
	"time"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/fox-one/mixin-sdk-go/v2/mixinnet"
	"github.com/shopspring/decimal"

	"github.com/rs/zerolog/log"
)

const (
	AGGREGRATE_UTXO_MEMO = "aggregate_utxo"
)

// 主动聚合utxos 至 utxo 数量不超过 255 个
func (m *MixinClientWrapper) SyncArrgegateUtxos(ctx context.Context, assetId string) (utxos []*mixin.SafeUtxo, err error) {
	m.transferMutex.Lock()
	defer m.transferMutex.Unlock()
	utxos = make([]*mixin.SafeUtxo, 0)
	for {
		requestId := utils.RandomTraceID()
		utxos, err = m.Client.SafeListUtxos(ctx, mixin.SafeListUtxoOption{
			Asset:     assetId,
			State:     mixin.SafeUtxoStateUnspent,
			Threshold: 1,
		})

		if err != nil {
			log.Error().Err(err).Msg("list utxos failed")
			continue
		}

		if len(utxos) <= MAX_UTXO_NUM {
			// 主动聚合完成
			break
		}

		// 将utxos分割成 255 个大小的数组
		utxoSlice := make([]*mixin.SafeUtxo, 0, MAX_UTXO_NUM)
		utxoSliceAmount := decimal.Zero
		for i := 0; i < len(utxos); i++ {
			if utxos[i].InscriptionHash.HasValue() {
				continue
			}
			utxoSlice = append(utxoSlice, utxos[i])
			utxoSliceAmount = utxoSliceAmount.Add(utxos[i].Amount)
		}

		// 2: build transaction
		b := mixin.NewSafeTransactionBuilder(utxoSlice)
		b.Memo = AGGREGRATE_UTXO_MEMO
		var tx *mixinnet.Transaction
		tx, err = m.Client.MakeTransaction(ctx, b, []*mixin.TransactionOutput{
			{
				Address: mixin.RequireNewMixAddress([]string{m.Client.ClientID}, 1),
				Amount:  utxoSliceAmount,
			},
		})
		if err != nil {
			return nil, err
		}
		var raw string
		raw, err = tx.Dump()
		if err != nil {
			return
		}

		// 3. create transaction
		var request *mixin.SafeTransactionRequest
		request, err = m.Client.SafeCreateTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
			RequestID:      requestId,
			RawTransaction: raw,
		})
		if err != nil {
			return
		}

		// 4. sign transaction
		err = mixin.SafeSignTransaction(
			tx,
			m.SpendKey,
			request.Views,
			0,
		)
		if err != nil {
			return
		}

		var signedRaw string
		signedRaw, err = tx.Dump()
		if err != nil {
			return
		}

		// 5. submit transaction
		_, err = m.Client.SafeSubmitTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
			RequestID:      requestId,
			RawTransaction: signedRaw,
		})
		if err != nil {
			return
		}

		// 重试读取交易状态
		const defaultMaxRetryTimes = 3
		retryTimes := 0
		for {
			if retryTimes >= defaultMaxRetryTimes {
				break
			}

			retryTimes++
			time.Sleep(time.Second * time.Duration(retryTimes))
			_, err = m.Client.SafeReadTransactionRequest(ctx, requestId)
			if err != nil {
				continue
			} else {
				break
			}
		}

		// 等待 250ms
		time.Sleep(time.Second >> 2)
	}
	return
}
