package model

import (
	"context"

	"github.com/fox-one/pkg/store2"
	"gorm.io/gorm"
)

func init() {
	store2.RegisterMigrate(func(db *store2.DB) error {
		var tx *gorm.DB
		tx = db.Update().Model(&User{})
		if err := tx.AutoMigrate(&User{}); err != nil {
			return err
		}

		tx = db.Update().Model(&Project{})
		if err := tx.AutoMigrate(&Project{}); err != nil {
			return err
		}

		tx = db.Update().Model(&DonateAction{})
		if err := tx.AutoMigrate(&DonateAction{}); err != nil {
			return err
		}

		tx = db.Update().Model(&Asset{})
		if err := tx.AutoMigrate(&Asset{}); err != nil {
			return err
		}

		tx = db.Update().Model(&Snapshot{})
		if err := tx.AutoMigrate(&Snapshot{}); err != nil {
			return err
		}
		return nil
	})
}

type UserStore interface {
	// 创建用户
	AddUser(ctx context.Context, user *User) error
	// 查询所有用户
	ListUsers(ctx context.Context) ([]*User, error)
	// 根据 did 查询用户
	GetUserByDID(ctx context.Context, did string) (*User, error)
	// 根据 uid 查询用户
	GetUserByUID(ctx context.Context, uid string) (*User, error)
	GetUserByIdentityNumber(ctx context.Context, ident string) (*User, error)
	// 更新用户信息
	UpdateUserBymuid(ctx context.Context, muid string, user *User) error
}

type ProjectStore interface {
	// 创建项目
	AddProject(ctx context.Context, item *Project) error
	// 删除项目
	DeleteProject(ctx context.Context, id string) error
	// 查询所有的项目
	ListProjects(ctx context.Context, limit, offset int64) ([]*Project, error)
	// 根据 identity_number 查询项目
	GetProjectsByIdentityNumber(ctx context.Context, ident string, limit, offset int64) ([]*Project, error)
	// 根据 id 查询项目
	GetProject(ctx context.Context, pid string) (*Project, error)
	// Incr project donate cnt
	IncrProjectDonateCnt(ctx context.Context, pid string) error
}

type DonateActionStore interface {
	// 添加捐赠记录
	AddDonateAction(ctx context.Context, action *DonateAction) error
	// 查询某个 did 的所有捐赠记录
	QueryDonateActionsByIdentityNumber(ctx context.Context, ident string) ([]*DonateAction, error)
	// 查询某个 项目 的被捐赠记录
	QueryDonateActionsByPID(ctx context.Context, pid string) ([]*DonateAction, error)
}

type AssetStore interface {
	// 列出资产
	ListAssets(ctx context.Context) ([]*Asset, error)
	// 获取某个资产
	GetAsset(ctx context.Context, id string) (*Asset, error)
	// 添加资产
	AddAsset(ctx context.Context, asset *Asset) error
}

type SnapshotStore interface {
	UpsertSnapshot(ctx context.Context, snapshot *Snapshot) error
	GetSnapshotCount(ctx context.Context) (int64, error)
	InsertSnapshot(ctx context.Context, snapshot *Snapshot) error
	GetSnapshotById(ctx context.Context, snapshotId string) (*Snapshot, error)
	GetLastestSnapshot(ctx context.Context) (*Snapshot, error)
}
