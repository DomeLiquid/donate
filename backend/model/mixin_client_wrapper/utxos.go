package mixin_client_wrapper

import (
	"context"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/shopspring/decimal"
)

func (m *MixinClientWrapper) ListUtxos(ctx context.Context) ([]*mixin.SafeUtxo, error) {
	utxos, _ := m.Client.SafeListUtxos(ctx, mixin.SafeListUtxoOption{
		// Order:     "DESC",
		// Asset:     "15f4e203-70e8-39d0-88a4-4ed8ee462dfd",
		State:     mixin.SafeUtxoStateUnspent,
		Threshold: 1,
		Limit:     30,
	})

	return utxos, nil
}

func (m *MixinClientWrapper) ListAssets(ctx context.Context) ([]*mixin.SafeAsset, error) {
	assets, err := m.Client.SafeReadAssets(ctx)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (m *MixinClientWrapper) GetAsset(ctx context.Context, assetId string) (*mixin.SafeAsset, error) {
	asset, err := m.Client.SafeReadAsset(ctx, assetId)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

// 获取某种资产的总数额
// !!! TODO
func (m *MixinClientWrapper) GetAssetTotalAmount(ctx context.Context, assetId string) (decimal.Decimal, error) {
	asset, err := m.Client.SafeReadAsset(ctx, assetId)
	if err != nil {
		return decimal.Zero, err
	}
	return asset.PriceBTC, nil
}
