package main

import (
	"context"
	"donate/config"
	"donate/model"
	"donate/router"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fox-one/pkg/db"
	"github.com/fox-one/pkg/store2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var _ model.DonateAction

var (
	configFile = flag.String("f", "~/.config/dome_loop_config_debug.json", "the config file")
	signalChan = make(chan os.Signal, 1)
)

func main() {
	flag.Parse()

	conf := new(config.Config)

	if err := config.Init(*configFile, conf); err != nil {
		log.Error().Err(err).Msgf("[configs.Init] err: %+v", err)
	}
	log.Debug().Any("conf", conf).Msg("init config success")

	db, err := provideDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// go func() {
	// 	sigCh := make(chan os.Signal, 1)
	// 	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	// 	<-sigCh
	// 	log.Info().Msg("Received shutdown signal, exiting...")
	// 	cancel()
	// }()

	router := router.NewService(conf, db)

	err = router.Run(conf.Port)
	if err != nil {
		log.Error().Err(err).Msg("run router failed")
	}
}

func connectDatabase(cfg db.Config, timeout time.Duration) (*store2.DB, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	dur := time.Millisecond
	// TODO gorm config
	gormCfg := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		CreateBatchSize:                          5000,
	}

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("connect db: %w", ctx.Err())
		case <-time.After(dur):
			if conn, err := store2.Open(cfg, gormCfg); err == nil {
				return conn, nil
			}

			dur = time.Second
		}
	}
}
func provideDatabase() (*store2.DB, error) {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("获取用户主目录失败: %w", err)
	}

	// 构建sqlite3数据库文件路径
	dbPath := filepath.Join(homeDir, "donate.sqlite3")

	// 确保数据库文件所在目录存在
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("创建数据库目录失败: %w", err)
	}

	cfg := db.Config{
		Dialect: "sqlite3",
		Host:    dbPath,
		Debug:   true,
	}

	conn, err := connectDatabase(cfg, 8*time.Second)
	if err != nil {
		return nil, err
	}

	if err := store2.Migrate(conn); err != nil {
		return nil, err
	}

	return conn, nil
}
