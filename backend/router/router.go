package router

import (
	"context"
	"donate/clock"
	"donate/config"
	"donate/model"
	"donate/model/mixin_client_wrapper"
	"donate/router/api"
	publicMiddleware "donate/router/middleware"

	"github.com/fox-one/pkg/store2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Service struct {
	clock clock.Clock

	store       model.Store
	mixinClient *mixin_client_wrapper.MixinClientWrapper
	apiServer   *api.ApiServer
	router      *gin.Engine
}

func NewService(conf *config.Config, db *store2.DB) *Service {
	mixinClient, err := mixin_client_wrapper.NewMixinClientWrapper(conf.MixinConfig)
	if err != nil {
		panic(err)
	}
	store := model.NewStore(db)

	srv := &Service{
		clock:       clock.New(),
		store:       store,
		mixinClient: mixinClient,
		apiServer:   api.New(mixinClient, store),
	}
	srv.initRouter()

	return srv
}

func (s *Service) initRouter() {
	router := gin.New()
	logger := log.Logger.With().Logger()

	router.Use(
		publicMiddleware.Cors(),
		publicMiddleware.GinXid(&logger),
		publicMiddleware.GinLogger(&logger),
		publicMiddleware.GinRecovery(&logger, true),
	)
	router.GET("/project/:item", s.apiServer.GetProject)
	router.GET("/donate-users/:pid", s.apiServer.GetDonateUsersByPid)
	router.GET("/projects", s.apiServer.GetProjects)
	router.GET("/user/:ident", s.apiServer.GetUserByIdentityNumber)
	router.GET("/users/search", s.apiServer.SearchUser)
	router.GET("/users-donate/:ident", s.apiServer.GetProjectsByIdentityNumber)
	router.GET("/assets", s.apiServer.GetAssets) // 提供支持捐赠的资产 以及资产价格

	s.router = router
}

func (s *Service) Run(addr string) error {

	// g, ctx := errgroup.WithContext(ctx)
	// g.Go(func() error {
	// 	return s.router.Run(addr)
	// })
	go s.RunMixinLoop(context.Background())
	return s.router.Run(addr)

	// g.Go(func() error {
	// 	// 后台轮询转账
	// 	return nil
	// })

	// return g.Wait()
}
