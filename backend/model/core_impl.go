package model

import (
	"context"

	"github.com/fox-one/pkg/store2"
	"gorm.io/gorm"
)

func NewStore(db *store2.DB) Store {
	return Store{
		UserStore:         NewUserStore(db),
		ProjectStore:      NewProjectStore(db),
		DonateActionStore: NewDonateActionStore(db),
		AssetStore:        NewAssetStore(db),
		SnapshotStore:     NewSnapshotStore(db),
	}
}

type Store struct {
	UserStore
	ProjectStore
	DonateActionStore
	AssetStore
	SnapshotStore
}

func NewAssetStore(db *store2.DB) AssetStore {
	return &assetStore{db: db}
}

type assetStore struct {
	db *store2.DB
}

// 列出资产
func (a *assetStore) ListAssets(ctx context.Context) (assets []*Asset, err error) {
	err = a.db.View().Find(&assets).Error
	return assets, err
}

// 获取某个资产
func (a *assetStore) GetAsset(ctx context.Context, id string) (asset *Asset, err error) {
	asset = &Asset{}
	err = a.db.View().Where("id = ?", id).First(asset).Error
	return
}

func (a *assetStore) AddAsset(ctx context.Context, asset *Asset) (err error) {
	return a.db.Update().Create(asset).Error
}

type store struct {
	db *store2.DB
}

// User 实现
type userStore struct {
	*store
}

func NewUserStore(db *store2.DB) UserStore {
	return &userStore{&store{db: db}}
}

func (s *userStore) AddUser(ctx context.Context, user *User) error {
	return s.db.Create(user).Error
}

func (s *userStore) ListUsers(ctx context.Context) ([]*User, error) {
	var users []*User
	err := s.db.Find(&users).Error
	return users, err
}

func (s *userStore) GetUserByDID(ctx context.Context, did string) (*User, error) {
	var user User
	err := s.db.Where("did = ?", did).First(&user).Error
	return &user, err
}

func (s *userStore) GetUserByUID(ctx context.Context, mixin_uid string) (*User, error) {
	var user User
	err := s.db.Where("mixin_uid = ?", mixin_uid).First(&user).Error
	return &user, err
}

func (s *userStore) GetUserByIdentityNumber(ctx context.Context, ident string) (*User, error) {
	var user User
	err := s.db.Where("identity_number = ?", ident).First(&user).Error
	return &user, err
}

func (s *userStore) UpdateUserBymuid(ctx context.Context, mixin_uid string, user *User) error {
	// 根据muid更新用户信息
	return s.db.Where("mixin_uid = ?", mixin_uid).Updates(user).Error
}

// DonateItem 实现
type projectStore struct {
	*store
}

func NewProjectStore(db *store2.DB) ProjectStore {
	return &projectStore{&store{db: db}}
}

func (s *projectStore) AddProject(ctx context.Context, project *Project) error {
	return s.db.Create(project).Error
}

func (s *projectStore) DeleteProject(ctx context.Context, id string) error {
	return s.db.Delete(&Project{}, "id = ?", id).Error
}

func (s *projectStore) ListProjects(ctx context.Context, limit, offset int64) (projects []*Project, err error) {
	err = s.db.Order("donate_cnt DESC").Limit(int(limit)).Offset(int(offset)).Find(&projects).Error
	return
}

func (s *projectStore) GetProject(ctx context.Context, id string) (project *Project, err error) {
	err = s.db.Where("pid = ?", id).First(&project).Error
	return
}
func (s *projectStore) GetProjectsByIdentityNumber(ctx context.Context, ident string, limit, offset int64) (projects []*Project, err error) {
	err = s.db.Where("identity_number = ?", ident).Order("donate_cnt DESC").Limit(int(limit)).Offset(int(offset)).Find(&projects).Error
	return
}

func (s *projectStore) IncrProjectDonateCnt(ctx context.Context, pid string) error {
	return s.db.Model(&Project{}).Where("pid = ?", pid).Update("donate_cnt", gorm.Expr("donate_cnt + 1")).Error
}

// DonateAction 实现
type donateActionStore struct {
	*store
}

func NewDonateActionStore(db *store2.DB) DonateActionStore {
	return &donateActionStore{&store{db: db}}
}

func (s *donateActionStore) AddDonateAction(ctx context.Context, action *DonateAction) error {
	return s.db.Create(action).Error
}

func (s *donateActionStore) QueryDonateActionsByIdentityNumber(ctx context.Context, ident string) ([]*DonateAction, error) {
	var actions []*DonateAction
	err := s.db.Where("identity_number = ?", ident).Find(&actions).Error
	return actions, err
}

func (s *donateActionStore) QueryDonateActionsByPID(ctx context.Context, pid string) ([]*DonateAction, error) {
	var actions []*DonateAction
	err := s.db.Where("pid = ?", pid).Find(&actions).Error
	return actions, err
}

type snapshotStore struct {
	*store
}

func NewSnapshotStore(db *store2.DB) SnapshotStore {
	return &snapshotStore{&store{db: db}}
}

func (s *snapshotStore) UpsertSnapshot(ctx context.Context, snapshot *Snapshot) error {
	return s.db.Tx(func(tx *store2.DB) error {
		return tx.Save(snapshot).Error
	})
}

func (s *snapshotStore) GetSnapshotCount(ctx context.Context) (int64, error) {
	var count int64
	if err := s.db.View().Model(&Snapshot{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *snapshotStore) InsertSnapshot(ctx context.Context, snapshot *Snapshot) error {
	return s.db.Tx(func(tx *store2.DB) error {
		return tx.Create(snapshot).Error
	})
}

func (s *snapshotStore) GetSnapshotById(ctx context.Context, snapshotId string) (*Snapshot, error) {
	var snapshot Snapshot
	if err := s.db.View().Where("snapshot_id = ?", snapshotId).First(&snapshot).Error; err != nil {
		return nil, err
	}
	return &snapshot, nil
}

func (s *snapshotStore) GetLastestSnapshot(ctx context.Context) (*Snapshot, error) {
	var snapshot Snapshot
	if err := s.db.View().Order("created_at DESC").First(&snapshot).Error; err != nil {
		return nil, err
	}
	return &snapshot, nil
}
