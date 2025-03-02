package mixin_client_wrapper

import (
	"context"
	"errors"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/samber/lo"
)

func (m *MixinClientWrapper) InscriptionTransfer(ctx context.Context, req *InscriptionTransferRequest) error {
	var utxos []*mixin.SafeUtxo
	utxos, _ = m.Client.SafeListUtxos(ctx, mixin.SafeListUtxoOption{
		Asset:     req.AssetId,
		State:     mixin.SafeUtxoStateUnspent,
		Threshold: 1,
		Limit:     500,
	})

	utxos = lo.Filter(utxos, func(utxo *mixin.SafeUtxo, _ int) bool {
		return utxo.InscriptionHash.String() == req.Inscription
	})

	if len(utxos) == 0 {
		return errors.New("inscription not found")
	}
	if len(utxos) > 1 {
		return errors.New("multiple inscriptions found")
	}

	b := mixin.NewSafeTransactionBuilder(utxos)
	b.Memo = req.Memo

	txOutout := &mixin.TransactionOutput{
		Address: mixin.RequireNewMixAddress([]string{req.Member}, 1),
		Amount:  utxos[0].Amount,
	}

	tx, err := m.Client.MakeTransaction(ctx, b, []*mixin.TransactionOutput{txOutout})
	if err != nil {
		return err
	}

	raw, err := tx.Dump()
	if err != nil {
		return err
	}

	// 3. create transaction
	request, err := m.Client.SafeCreateTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      req.RequestId,
		RawTransaction: raw,
	})
	if err != nil {
		return err
	}
	// 4. sign transaction
	err = mixin.SafeSignTransaction(
		tx,
		m.SpendKey,
		request.Views,
		0,
	)
	if err != nil {
		return err
	}
	signedRaw, err := tx.Dump()
	if err != nil {
		return err
	}

	// 5. submit transaction
	_, err = m.Client.SafeSubmitTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID:      req.RequestId,
		RawTransaction: signedRaw,
	})
	if err != nil {
		return err
	}

	// 6. read transaction
	_, err = m.Client.SafeReadTransactionRequest(ctx, req.RequestId)
	if err != nil {
		return err
	}
	return nil
}
