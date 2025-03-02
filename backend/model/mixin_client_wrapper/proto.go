package mixin_client_wrapper

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	defaultMaxMixinRetry = 3

	mixinAssetAmountCacheTTL   = 30 * time.Minute
	mixinAssetAmountCacheDelay = time.Hour

	MAX_UTXO_NUM = 255
)

type TransferOneRequest struct {
	RequestId string
	AssetId   string

	Member string
	Amount decimal.Decimal
	Memo   string
}

func NewMemberAmount(member []string, amount decimal.Decimal) MemberAmount {
	return MemberAmount{
		Member: member,
		Amount: amount,
	}
}

type MemberAmount struct {
	Member []string
	Amount decimal.Decimal
}

type TransferManyRequest struct {
	RequestId string
	AssetId   string

	MemberAmount []MemberAmount
	Memo         string
}

type InscriptionTransferRequest struct {
	RequestId   string
	AssetId     string
	Inscription string

	Memo   string
	Member string
}
