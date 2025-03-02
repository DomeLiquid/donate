package mixin_client_wrapper

// func (m *MixinClientWrapper) GetUserMixinAssetAmount(ctx context.Context, accessToken, assetId string) decimal.Decimal {
// 	key := fmt.Sprintf("%s-%s", accessToken, assetId)
// 	val, err := m.userMixinAssetAmountCache.Do(key, func() (interface{}, error) {
// 		// 检查限流器是否可以立即获取许可
// 		// if !m.rateLimiter.Allow() {
// 		// 	return decimal.Zero, nil
// 		// }

// 		snapshots, err := mixin.ReadSafeSnapshots(ctx, accessToken, assetId, time.Time{}, "", 500)
// 		if err != nil {
// 			return decimal.Zero, err
// 		}

// 		var amount decimal.Decimal
// 		for _, snapshot := range snapshots {
// 			if snapshot.InscriptionHash != nil && snapshot.InscriptionHash.HasValue() {
// 				continue
// 			}
// 			amount = amount.Add(snapshot.Amount)
// 		}

// 		return amount, nil
// 	})

// 	if err != nil {
// 		return decimal.Zero
// 	}

// 	if v, ok := val.(decimal.Decimal); !ok {
// 		return decimal.Zero
// 	} else {
// 		return v
// 	}
// }
