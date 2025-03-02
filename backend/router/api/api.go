package api

import (
	"context"
	"donate/logger"
	"donate/model"
	"donate/model/mixin_client_wrapper"
	"donate/pkg/cacheflight"
	"donate/router/middleware"
	"donate/utils"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ApiServer struct {
	mixinClient *mixin_client_wrapper.MixinClientWrapper
	store       model.Store
	assetCf     *cacheflight.Group
}

func New(mixinClient *mixin_client_wrapper.MixinClientWrapper, store model.Store) *ApiServer {
	return &ApiServer{
		mixinClient: mixinClient,
		store:       store,
		assetCf:     cacheflight.New(time.Minute, time.Minute<<1),
	}
}

type GetProjectResponse struct {
	model.Project
	User *model.User `json:"user"`
}

// 1. 根据 base64 编码获取项目信息 + 捐赠过的用户
func (a *ApiServer) GetProject(ctx *gin.Context) {
	logger := ctx.MustGet(middleware.DefaultLoggerKey).(*logger.CtxLogger)
	encodedItem := ctx.Param("item")

	// base54decode
	pidUuid, err := uuid.FromString(encodedItem)
	if err == nil && pidUuid != uuid.Nil {
		project, err := a.store.GetProject(ctx, pidUuid.String())
		if err != nil {
			logger.Error().Err(err).Msg("failed to get project")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "failed to get project",
			})
			return
		}
		user, _ := a.store.GetUserByIdentityNumber(ctx, project.IdentityNumber)
		response := GetProjectResponse{
			Project: *project,
			User:    user,
		}

		ctx.JSON(http.StatusOK, response)
		return
	}

	// base64 解码
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedItem)
	if err != nil {
		logger.Error().Err(err).Msg("failed to decode base64 string")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid base64 string",
		})
		return
	}

	var donateItem struct {
		Title          string `json:"title"`          // required
		Description    string `json:"description"`    // optional
		ImgUrl         string `json:"imgUrl"`         // optional
		Link           string `json:"link"`           // optional
		IdentityNumber string `json:"identityNumber"` // required
	}
	if err := json.Unmarshal(decodedBytes, &donateItem); err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal json")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return
	}

	if len(donateItem.Title) == 0 || len(donateItem.IdentityNumber) == 0 {
		logger.Error().Msg("title or mixin_uid is empty")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "title or mixin_uid is empty",
		})
		return
	}
	mixinUser, err := a.mixinClient.Client.ReadUser(ctx, donateItem.IdentityNumber)
	if err != nil {
		logger.Error().Err(err).Msg("failed to read user")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read user",
		})
		return
	}

	// 检查用户在本地库是否存在,存在更新,不存在创建
	duser, err := a.store.GetUserByIdentityNumber(ctx, mixinUser.IdentityNumber)
	switch err {
	case nil:
		// 更新
		duser = &model.User{
			IdentityNumber: mixinUser.IdentityNumber,
			FullName:       mixinUser.FullName,
			MixinUID:       mixinUser.UserID,
			AvatarUrl:      mixinUser.AvatarURL,
			Biography:      mixinUser.Biography,
			MixinCreatedAt: mixinUser.CreatedAt,
			UpdatedAt:      time.Now(),
		}
		if err := a.store.UpdateUserBymuid(ctx, mixinUser.UserID, duser); err != nil {
			logger.Error().Err(err).Msg("failed to update user")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update user",
			})
			return
		}
	case gorm.ErrRecordNotFound:
		// 创建
		duser = &model.User{
			IdentityNumber: mixinUser.IdentityNumber,
			FullName:       mixinUser.FullName,
			MixinUID:       mixinUser.UserID,
			AvatarUrl:      mixinUser.AvatarURL,
			Biography:      mixinUser.Biography,
			MixinCreatedAt: mixinUser.CreatedAt,
			CreatedAt:      time.Now(),
		}
		if err := a.store.AddUser(ctx, duser); err != nil {
			logger.Error().Err(err).Msg("failed to create user")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create user",
			})
			return
		}
	default:
		logger.Error().Err(err).Msg("failed to get user")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user",
		})
		return
	}

	// 生成项目 ID
	pid := utils.GenUuidFromStrings("donate", donateItem.Title, donateItem.Description, donateItem.ImgUrl, mixinUser.UserID)
	project, err := a.store.GetProject(ctx, pid)
	switch err {
	case gorm.ErrRecordNotFound:
		item := &model.Project{
			PID:            pid,
			Title:          donateItem.Title,
			Description:    donateItem.Description,
			IdentityNumber: donateItem.IdentityNumber,
			ImgUrl:         donateItem.ImgUrl,
			Link:           donateItem.Link,
			MixinUID:       mixinUser.UserID,
			CreatedAt:      time.Now(),
		}
		user, _ := a.store.GetUserByIdentityNumber(ctx, mixinUser.IdentityNumber)

		err = a.store.AddProject(ctx, item)
		if err != nil {
			logger.Error().Err(err).Msg("failed to add project")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to add project",
			})
			return
		}
		project, _ = a.store.GetProject(ctx, pid)
		response := GetProjectResponse{
			Project: *project,
			User:    user,
		}

		// 直接将项目信息返回出去
		ctx.JSON(http.StatusOK, response)
		return
	case nil:
	default:
		logger.Error().Err(err).Msg("failed to get project")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get project",
		})
		return
	}
	user, _ := a.store.GetUserByIdentityNumber(ctx, project.IdentityNumber)
	response := GetProjectResponse{
		Project: *project,
		User:    user,
	}

	ctx.JSON(http.StatusOK, response)
	return
}

// 查询捐赠记录
type UserAction struct {
	IdentityNumber string          `json:"identityNumber"`
	FullName       string          `json:"fullName"`
	AvatarUrl      string          `json:"avatarUrl"`
	Biography      string          `json:"biography"`
	AssetID        string          `json:"assetId"`
	Amount         decimal.Decimal `json:"amount"`
	Asset          model.Asset     `json:"asset"`
	Project        model.Project   `json:"project"`
	User           model.User      `json:"user"` // 被捐赠者
}

// 1. 根据 pid 查询捐赠过的用户
func (a *ApiServer) GetDonateUsersByPid(ctx *gin.Context) {
	logger := ctx.MustGet(middleware.DefaultLoggerKey).(*logger.CtxLogger)
	pid := ctx.Param("pid")
	logger.Info().Str("pid", pid).Msg("GetDonateUsersByPid")

	_, err := uuid.FromString(pid)
	if err == nil {
		project, err := a.store.GetProject(ctx, pid)
		switch err {
		case gorm.ErrRecordNotFound:
			logger.Error().Err(err).Msg("project not found")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "project not found",
			})
			return
		case nil:
		default:
			logger.Error().Err(err).Msg("failed to get project")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get project",
			})
			return
		}
		var recipientUser *model.User
		recipientUser, _ = a.store.GetUserByIdentityNumber(ctx, project.IdentityNumber)

		var response []UserAction

		donateActions, err := a.store.QueryDonateActionsByPID(ctx, pid)
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "project not found",
			})
			return
		case nil:
		default:
			logger.Error().Err(err).Msg("failed to get donate actions")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get donate actions",
			})
			return
		}

		assetMap := a.getAssetMap()
		for _, action := range donateActions {
			user, err := a.store.GetUserByIdentityNumber(ctx, action.IdentityNumber)
			if err != nil {
				logger.Error().Err(err).Msg("get user failed")
				continue
			}
			asset, ok := assetMap[action.AssetID]
			if !ok {
				continue
			}

			userAction := &UserAction{
				IdentityNumber: action.IdentityNumber,
				FullName:       user.FullName,
				AvatarUrl:      user.AvatarUrl,
				Biography:      user.Biography,
				AssetID:        action.AssetID,
				Amount:         action.Amount,
				Asset:          *asset,
				Project:        *project,
				User:           *recipientUser,
			}
			response = append(response, *userAction)
		}

		ctx.JSON(http.StatusOK, response)
	} else {
		encodedItem := pid
		// base64 解码
		decodedBytes, err := base64.StdEncoding.DecodeString(encodedItem)
		if err != nil {
			logger.Error().Err(err).Msg("failed to decode base64 string")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid base64 string",
			})
			return
		}
		var donateItem struct {
			Title          string `json:"title"`          // required
			Description    string `json:"description"`    // optional
			ImgUrl         string `json:"imgUrl"`         // optional
			Link           string `json:"link"`           // optional
			IdentityNumber string `json:"identityNumber"` // required
		}
		if err := json.Unmarshal([]byte(decodedBytes), &donateItem); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal json")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid json",
			})
			return
		}

		if len(donateItem.Title) == 0 || len(donateItem.IdentityNumber) == 0 {
			logger.Error().Msg("title or mixin_uid is empty")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "title or mixin_uid is empty",
			})
			return
		}
		mixinUser, err := a.mixinClient.Client.ReadUser(ctx, donateItem.IdentityNumber)
		if err != nil {
			logger.Error().Err(err).Msg("failed to read user")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "failed to read user",
			})
			return
		}

		pid := utils.GenUuidFromStrings("donate", donateItem.Title, donateItem.Description, donateItem.ImgUrl, mixinUser.UserID)
		project, err := a.store.GetProject(ctx, pid)
		switch err {
		case gorm.ErrRecordNotFound:
			logger.Error().Err(err).Msg("project not found")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "project not found",
			})
			return
		case nil:
		default:
			logger.Error().Err(err).Msg("failed to get project")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get project",
			})
			return
		}
		var recipientUser *model.User
		recipientUser, _ = a.store.GetUserByIdentityNumber(ctx, project.IdentityNumber)

		var response []UserAction

		donateActions, err := a.store.QueryDonateActionsByPID(ctx, pid)
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "project not found",
			})
			return
		case nil:
		default:
			logger.Error().Err(err).Msg("failed to get donate actions")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get donate actions",
			})
			return
		}

		assetMap := a.getAssetMap()
		for _, action := range donateActions {
			user, err := a.store.GetUserByIdentityNumber(ctx, action.IdentityNumber)
			if err != nil {
				logger.Error().Err(err).Msg("get user failed")
				continue
			}
			asset, ok := assetMap[action.AssetID]
			if !ok {
				continue
			}

			userAction := &UserAction{
				IdentityNumber: action.IdentityNumber,
				FullName:       user.FullName,
				AvatarUrl:      user.AvatarUrl,
				Biography:      user.Biography,
				AssetID:        action.AssetID,
				Amount:         action.Amount,
				Asset:          *asset,
				Project:        *project,
				User:           *recipientUser,
			}
			response = append(response, *userAction)
		}

		ctx.JSON(http.StatusOK, response)
	}

	return
}

// 2. 项目列表
func (a *ApiServer) GetProjects(ctx *gin.Context) {
	logger := ctx.MustGet(middleware.DefaultLoggerKey).(*logger.CtxLogger)

	query := ctx.Request.URL.Query()
	limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
	if err != nil || offset < 0 {
		offset = 0
	}

	var response struct {
		Items []struct {
			Project *model.Project `json:"project"`
			User    *model.User    `json:"user"`
		} `json:"items"`
	}

	var projects []*model.Project
	if ident := query.Get("identity_number"); ident != "" {
		projects, err = a.store.GetProjectsByIdentityNumber(ctx, ident, limit, offset)
	} else {
		projects, err = a.store.ListProjects(ctx, limit, offset)
	}

	switch err {
	case gorm.ErrRecordNotFound:
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "projects not found",
		})
		return
	case nil:
	default:
		logger.Error().Err(err).Msg("failed to get projects")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get projects",
		})
		return
	}

	for _, project := range projects {
		user, _ := a.store.GetUserByIdentityNumber(ctx, project.IdentityNumber)
		response.Items = append(response.Items, struct {
			Project *model.Project `json:"project"`
			User    *model.User    `json:"user"`
		}{
			Project: project,
			User:    user,
		})
	}

	ctx.JSON(http.StatusOK, response)
	return

}

// SearchUser searches for users by identity number or name prefix
func (a *ApiServer) SearchUser(ctx *gin.Context) {
	logger := ctx.MustGet(middleware.DefaultLoggerKey).(*logger.CtxLogger)

	query := ctx.Request.URL.Query()
	ident, prefix := query.Get("identity_number"), query.Get("prefix")
	searchTerm := ident + prefix

	if searchTerm == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "identity_number or prefix is required"})
		return
	}

	users, err := a.store.UserStore.ListUsers(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list users")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search users"})
		return
	}

	var matchedUsers []*model.User
	for _, u := range users {
		if strings.HasPrefix(u.IdentityNumber, ident) || strings.HasPrefix(u.FullName, prefix) {
			matchedUsers = append(matchedUsers, u)
		}
	}

	if len(matchedUsers) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no users found"})
		return
	}

	ctx.JSON(http.StatusOK, matchedUsers)
}

// 3. 用户列表
func (a *ApiServer) GetUsers(ctx *gin.Context) {
	logger := ctx.MustGet(middleware.DefaultLoggerKey).(*logger.CtxLogger)
	users, err := a.store.UserStore.ListUsers(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list users")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (a *ApiServer) GetUserByIdentityNumber(ctx *gin.Context) {
	logger := ctx.MustGet(middleware.DefaultLoggerKey).(*logger.CtxLogger)
	ident := ctx.Param("ident")
	user, err := a.store.UserStore.GetUserByIdentityNumber(ctx, ident)
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		mixinUser, err := a.mixinClient.Client.ReadUser(ctx, ident)
		if err != nil {
			logger.Error().Err(err).Msg("failed to read user")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read user"})
			return
		}
		user = &model.User{
			IdentityNumber: mixinUser.IdentityNumber,
			FullName:       mixinUser.FullName,
			MixinUID:       mixinUser.UserID,
			AvatarUrl:      mixinUser.AvatarURL,
			Biography:      mixinUser.Biography,
			MixinCreatedAt: mixinUser.CreatedAt,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		err = a.store.UserStore.AddUser(ctx, user)
		if err != nil {
			logger.Error().Err(err).Msg("failed to create user")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}
		user, err = a.store.UserStore.GetUserByIdentityNumber(ctx, ident)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get user")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
			return
		}
		ctx.JSON(http.StatusOK, user)
		return
	default:
	}
	ctx.JSON(http.StatusOK, user)
	return
}

// 4. 查询用户捐赠的项目列表
func (a *ApiServer) GetProjectsByIdentityNumber(ctx *gin.Context) {
	logger := ctx.MustGet(middleware.DefaultLoggerKey).(*logger.CtxLogger)

	ident := ctx.Param("ident")
	if ident == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "identity_number is required"})
		return
	}
	user, err := a.store.UserStore.GetUserByIdentityNumber(ctx, ident)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	actions, err := a.store.DonateActionStore.QueryDonateActionsByIdentityNumber(ctx, ident)
	if err != nil {
		logger.Error().Err(err).Msg("failed to query donate actions")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query donate actions"})
		return
	}

	assetMap := a.getAssetMap()
	var response []*UserAction

	for _, action := range actions {
		project, err := a.store.GetProject(ctx, action.PID)
		if err != nil {
			continue
		}
		var recipUser *model.User
		recipUser, _ = a.store.GetUserByIdentityNumber(ctx, project.IdentityNumber)

		asset, ok := assetMap[action.AssetID]
		if !ok {
			continue
		}
		res := &UserAction{
			IdentityNumber: user.IdentityNumber,
			FullName:       user.FullName,
			AvatarUrl:      user.AvatarUrl,
			Biography:      user.Biography,
			AssetID:        action.AssetID,
			Amount:         action.Amount,
			Asset:          *asset,
			Project:        *project,
			User:           *recipUser,
		}
		response = append(response, res)
	}

	ctx.JSON(http.StatusOK, response)
}

func (a *ApiServer) GetAssets(ctx *gin.Context) {
	assetA, err := a.assetCf.Do("all_asset", func() (val interface{}, err error) {
		return a.mixinClient.Client.SafeReadAssets(ctx)
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get assets"})
		return
	}

	assets := assetA.([]*mixin.SafeAsset)

	response := make([]*model.Asset, len(assets))

	response = make([]*model.Asset, 0, len(assets))
	// 初始化一个非空的response slice避免nil panic
	assetMap := make(map[string]*model.Asset)

	for _, a := range assets {
		asset := &model.Asset{
			AssetID:  a.AssetID,
			Symbol:   a.Symbol,
			Name:     a.Name,
			IconURL:  a.IconURL,
			ChainID:  a.ChainID,
			PriceUSD: a.PriceUSD,
		}
		response = append(response, asset)
		assetMap[a.AssetID] = asset
	}

	// 只处理非空的response元素
	for _, r := range response {
		if r != nil && r.ChainID != "" {
			if c, ok := assetMap[r.ChainID]; ok {
				r.ChainIconURL = c.IconURL
				r.ChainSymbol = c.Symbol
			}
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func (a *ApiServer) getAssetMap() (assetMap map[string]*model.Asset) {
	assetA, err := a.assetCf.Do("all_asset", func() (val interface{}, err error) {
		return a.mixinClient.Client.SafeReadAssets(context.TODO())
	})
	if err != nil {
		return
	}

	assets := assetA.([]*mixin.SafeAsset)
	// 初始化一个非空的response slice避免nil panic
	assetMap = make(map[string]*model.Asset)
	response := make([]*model.Asset, 0, len(assets))

	for _, a := range assets {
		asset := &model.Asset{
			AssetID:  a.AssetID,
			Symbol:   a.Symbol,
			Name:     a.Name,
			IconURL:  a.IconURL,
			ChainID:  a.ChainID,
			PriceUSD: a.PriceUSD,
		}
		response = append(response, asset)
		assetMap[a.AssetID] = asset
	}
	// 只处理非空的response元素
	for _, r := range response {
		if r != nil && r.ChainID != "" {
			if c, ok := assetMap[r.ChainID]; ok {
				r.ChainIconURL = c.IconURL
				r.ChainSymbol = c.Symbol
			}
		}
	}

	for _, rr := range response {
		assetMap[rr.AssetID] = rr
	}

	return
}
