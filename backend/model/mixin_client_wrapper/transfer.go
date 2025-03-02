package mixin_client_wrapper

import (
	"context"
	"donate/utils"
	"sort"
	"strconv"
	"time"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

// func (m *MixinClientWrapper) TransferOneWithStore(ctx context.Context, logger core.Log, clk clock.Clock, mixinTransactionStore core.MixinTransactionStore, payment *core.Payment, req *TransferOneRequest) error {
// 	var err error
// 	err = mixinTransactionStore.CreateMixinTransaction(ctx, &core.MixinTransaction{
// 		RequestId: req.RequestId,
// 		PaymentId: payment.RequestId,
// 		Uid:       payment.Uid,
// 		Status:    core.MixinTransactionStatusPending,
// 		Memo:      req.Memo,

// 		CreatedAt: clk.Now().Unix(),
// 		UpdatedAt: clk.Now().Unix(),
// 	})
// 	if err != nil {
// 		logger.Error().Err(err).Msg("create mixin transaction failed")
// 		return err
// 	}

// 	err = m.TransferOneWithRetry(ctx, req)
// 	if err != nil {
// 		err1 := mixinTransactionStore.UpdateMixinTransactionStatus(ctx, req.RequestId, core.MixinTransactionStatusFailed)
// 		if err1 != nil {
// 			logger.Error().Err(err1).Msg("update mixin transaction failed")
// 		}
// 		logger.Error().Err(err).Msg("transfer one failed")
// 		return err
// 	}

// 	err = mixinTransactionStore.UpdateMixinTransactionStatus(ctx, req.RequestId, core.MixinTransactionStatusConfirmed)
// 	if err != nil {
// 		logger.Error().Err(err).Msg("update mixin transaction status failed")
// 		return err
// 	}

// 	return nil
// }

func (m *MixinClientWrapper) TransferManyWithRetry(ctx context.Context, requestId string, assetId string, memberAmounts []MemberAmount, memo string) error {
	if len(memberAmounts) < MAX_UTXO_NUM {
		req := &TransferManyRequest{
			RequestId:    requestId,
			AssetId:      assetId,
			MemberAmount: memberAmounts,
			Memo:         memo,
		}
		return m.transferManyWithRetry(ctx, req)
	} else {
		transferMany := buildTransferMany(memberAmounts)
		for index, transfer := range transferMany {
			req := &TransferManyRequest{
				RequestId:    utils.GenUuidFromStrings(requestId, strconv.Itoa(index)),
				AssetId:      assetId,
				MemberAmount: transfer,
				Memo:         memo,
			}

			err := m.transferManyWithRetry(ctx, req)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MixinClientWrapper) transferManyWithRetry(ctx context.Context, req *TransferManyRequest) error {
	var err error
	for i := 0; i < defaultMaxMixinRetry; i++ {
		if _, err = m.transferMany(ctx, req); err != nil {
			log.Error().Err(err).Msg("send transfer many failed, retrying...")
			time.Sleep(time.Second << i)
			continue
		} else {
			return nil
		}
	}
	return err
}

func (m *MixinClientWrapper) TransferOneWithRetry(ctx context.Context, req *TransferOneRequest) error {
	var err error
	for i := 0; i < defaultMaxMixinRetry; i++ {
		if _, err = m.transferOne(ctx, req); err != nil {
			log.Error().Err(err).Msg("send transfer one failed, retrying...")
			time.Sleep(time.Second << i)
			continue
		} else {
			return nil
		}
	}
	return err
}

func (m *MixinClientWrapper) transferOne(ctx context.Context, req *TransferOneRequest) (*mixin.SafeTransactionRequest, error) {
	var err error
	var utxos []*mixin.SafeUtxo

	utxos, err = m.SyncArrgegateUtxos(ctx, req.AssetId)
	if err != nil {
		return nil, err
	}

	m.transferMutex.Lock()
	defer m.transferMutex.Unlock()

	for i := 0; i < 3 && len(utxos) == 0; i++ {
		utxos, _ = m.SafeListUtxos(ctx, mixin.SafeListUtxoOption{
			Asset:     req.AssetId,
			State:     mixin.SafeUtxoStateUnspent,
			Threshold: 1,
			Limit:     500,
		})
		if len(utxos) > 0 {
			break
		}
		time.Sleep(time.Second << 1)
	}

	if len(utxos) == 0 {
		return nil, ErrNotEnoughUtxos
	}

	for i := 0; i < len(utxos); i++ {
		if utxos[i].InscriptionHash.HasValue() {
			utxos = append(utxos[:i], utxos[i+1:]...)
			i--
		}
	}

	// 1: select utxos
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Amount.LessThanOrEqual(utxos[j].Amount)
	})
	var useAmount decimal.Decimal
	var useUtxos []*mixin.SafeUtxo

	for _, utxo := range utxos {
		useAmount = useAmount.Add(utxo.Amount)
		useUtxos = append(useUtxos, utxo)

		if len(useUtxos) > MAX_UTXO_NUM {
			useUtxos = useUtxos[1:]
			useAmount = useAmount.Sub(utxos[0].Amount)
		}

		if useAmount.GreaterThanOrEqual(req.Amount) {
			break
		}
	}

	if useAmount.LessThan(req.Amount) {
		return nil, ErrNotEnoughUtxos
	}

	// 2: build transaction
	b := mixin.NewSafeTransactionBuilder(useUtxos)
	b.Memo = req.Memo

	txOutout := &mixin.TransactionOutput{
		Address: mixin.RequireNewMixAddress([]string{req.Member}, 1),
		Amount:  req.Amount,
	}

	tx, err := m.Client.MakeTransaction(ctx, b, []*mixin.TransactionOutput{txOutout})
	if err != nil {
		return nil, err
	}

	raw, err := tx.Dump()
	if err != nil {
		return nil, err
	}

	// 3. create transaction
	request, err := m.Client.SafeCreateTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      req.RequestId,
		RawTransaction: raw,
	})
	if err != nil {
		return nil, err
	}
	// 4. sign transaction
	err = mixin.SafeSignTransaction(
		tx,
		m.SpendKey,
		request.Views,
		0,
	)
	if err != nil {
		return nil, err
	}
	signedRaw, err := tx.Dump()
	if err != nil {
		return nil, err
	}

	// 5. submit transaction
	_, err = m.Client.SafeSubmitTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      req.RequestId,
		RawTransaction: signedRaw,
	})
	if err != nil {
		return nil, err
	}

	// 6. read transaction
	req1, err := m.Client.SafeReadTransactionRequest(ctx, req.RequestId)
	if err != nil {
		return nil, err
	}
	return req1, nil
}

func (m *MixinClientWrapper) transferMany(ctx context.Context, req *TransferManyRequest) (*mixin.SafeTransactionRequest, error) {
	var utxos []*mixin.SafeUtxo
	var err error
	utxos, err = m.SyncArrgegateUtxos(ctx, req.AssetId)
	if err != nil {
		return nil, err
	}

	m.transferMutex.Lock()
	defer m.transferMutex.Unlock()

	totalAmount := decimal.Zero
	lo.ForEach(req.MemberAmount, func(item MemberAmount, _ int) {
		totalAmount = totalAmount.Add(item.Amount)
	})
	retryCount := 0
	for len(utxos) == 0 && retryCount < 3 {
		// 1. 将utxos聚合
		utxos, _ = m.Client.SafeListUtxos(ctx, mixin.SafeListUtxoOption{
			Asset:     req.AssetId,
			State:     mixin.SafeUtxoStateUnspent,
			Threshold: 1,
			Limit:     500,
		})
		time.Sleep(time.Second << retryCount)
		retryCount++
	}
	if len(utxos) == 0 {
		return nil, ErrNotEnoughUtxos
	}

	// 1: select utxos
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Amount.LessThanOrEqual(utxos[j].Amount)
	})

	var useAmount decimal.Decimal
	var useUtxos []*mixin.SafeUtxo
	for _, utxo := range utxos {
		useAmount = useAmount.Add(utxo.Amount)
		useUtxos = append(useUtxos, utxo)

		if len(useUtxos) > MAX_UTXO_NUM {
			useUtxos = useUtxos[1:]
			useAmount = useAmount.Sub(utxos[0].Amount)
		}

		if useAmount.GreaterThanOrEqual(totalAmount) {
			break
		}
	}

	if useAmount.LessThan(totalAmount) {
		return nil, ErrNotEnoughUtxos
	}

	// 2: build transaction
	b := mixin.NewSafeTransactionBuilder(useUtxos)
	b.Memo = req.Memo

	txOutout := make([]*mixin.TransactionOutput, len(req.MemberAmount))
	for i := 0; i < len(req.MemberAmount); i++ {
		txOutout[i] = &mixin.TransactionOutput{
			Address: mixin.RequireNewMixAddress(req.MemberAmount[i].Member, byte(len(req.MemberAmount[i].Member))),
			Amount:  req.MemberAmount[i].Amount,
		}
	}

	tx, err := m.Client.MakeTransaction(ctx, b, txOutout)
	if err != nil {
		return nil, err
	}

	raw, err := tx.Dump()
	if err != nil {
		return nil, err
	}

	// 3. create transaction
	request, err := m.Client.SafeCreateTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      req.RequestId,
		RawTransaction: raw,
	})
	if err != nil {
		return nil, err
	}
	// 4. sign transaction
	err = mixin.SafeSignTransaction(
		tx,
		m.SpendKey,
		request.Views,
		0,
	)
	if err != nil {
		return nil, err
	}
	signedRaw, err := tx.Dump()
	if err != nil {
		return nil, err
	}

	// 5. submit transaction
	_, err = m.Client.SafeSubmitTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      req.RequestId,
		RawTransaction: signedRaw,
	})
	if err != nil {
		return nil, err
	}

	// 6. read transaction
	req1, err := m.Client.SafeReadTransactionRequest(ctx, req.RequestId)
	if err != nil {
		return nil, err
	}
	return req1, nil
}

// 一个功能函数，将一个数组中的多个元素切分成 n个数组，每个数组长度最多不超过255个
func buildTransferMany(memberAmounts []MemberAmount) [][]MemberAmount {
	result := make([][]MemberAmount, (len(memberAmounts)+MAX_UTXO_NUM-1)/MAX_UTXO_NUM)
	for i := 0; i < len(result); i++ {
		start := i * MAX_UTXO_NUM
		end := (i + 1) * MAX_UTXO_NUM
		if end > len(memberAmounts) {
			end = len(memberAmounts)
		}
		result[i] = memberAmounts[start:end]
	}
	return result
}
