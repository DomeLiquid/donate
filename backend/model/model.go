package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	MixinUID       string    `json:"-" gorm:"type:varchar(36);column:mixin_uid"` // mixin id
	IdentityNumber string    `json:"identityNumber" gorm:"type:varchar(255);column:identity_number"`
	FullName       string    `json:"fullName" gorm:"type:varchar(255);column:full_name"`
	AvatarUrl      string    `json:"avatarUrl" gorm:"type:varchar(255);column:avatar_url"`
	Biography      string    `json:"biography" gorm:"type:text;column:biography"`
	MixinCreatedAt time.Time `json:"-" gorm:"autoCreateTime;column:mixin_created_at"`
	CreatedAt      time.Time `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	UpdatedAt      time.Time `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at"`
}

type Project struct {
	PID            string    `json:"pid" gorm:"primaryKey;type:varchar(36);column:pid"`
	Title          string    `json:"title" gorm:"type:varchar(255);column:title"`
	Description    string    `json:"description,omitempty" gorm:"type:text;column:description"`
	ImgUrl         string    `json:"imgUrl,omitempty" gorm:"type:varchar(255);column:img_url"`
	Link           string    `json:"link,omitempty" gorm:"type:varchar(255);column:link"`
	IdentityNumber string    `json:"identityNumber" gorm:"type:varchar(255);column:identity_number"`
	MixinUID       string    `json:"-" gorm:"type:varchar(36);column:mixin_uid"` // mixin id
	DonateCnt      int64     `json:"donateCnt" gorm:"column:donate_cnt"`         // 被捐赠次数
	CreatedAt      time.Time `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
}

type DonateAction struct {
	ID             string          `json:"id" gorm:"primaryKey;type:varchar(36);column:id"`
	PID            string          `json:"pid" gorm:"primaryKey;type:varchar(36);column:pid"`
	IdentityNumber string          `json:"identityNumber" gorm:"type:varchar(255);column:identity_number"`
	AssetID        string          `json:"assetId" gorm:"type:varchar(36);column:asset_id"`
	Amount         decimal.Decimal `json:"amount" gorm:"type:decimal(64,8);column:amount"`
	CreatedAt      time.Time       `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
}

type Asset struct {
	AssetID      string          `json:"assetId" gorm:"column:asset_id;primaryKey;type:varchar(36)"`
	ChainID      string          `json:"chainId" gorm:"column:chain_id;type:varchar(36)"`
	ChainSymbol  string          `json:"chainSymbol" gorm:"column:chain_symbol;type:varchar(255)"`
	ChainIconURL string          `json:"chainIconUrl" gorm:"column:chain_icon_url;type:varchar(255)"`
	Symbol       string          `json:"symbol,omitempty" gorm:"column:symbol;type:varchar(255)"`
	Name         string          `json:"name,omitempty" gorm:"column:name;type:varchar(255)"`
	IconURL      string          `json:"iconUrl,omitempty" gorm:"column:icon_url;type:varchar(255)"`
	PriceUSD     decimal.Decimal `json:"priceUsd,omitempty" gorm:"column:priceUsd;type:decimal(64,8)"`
}

type Snapshot struct {
	SnapshotId string          `gorm:"column:snapshot_id;primaryKey" json:"snapshotId"`
	RequestId  string          `gorm:"column:request_id;index;type:varchar(36)" json:"requestId"`
	UserId     string          `gorm:"column:user_id;index;type:varchar(36)" json:"userId"`
	AssetId    string          `gorm:"column:asset_id;index;type:varchar(36)" json:"assetId"`
	Amount     decimal.Decimal `gorm:"column:amount;type:decimal(64,8)" json:"amount"`
	Memo       string          `gorm:"column:memo;type:varchar(512)" json:"memo"`
	CreatedAt  int64           `gorm:"column:created_at;type:int;not null" json:"createdAt"`
}
