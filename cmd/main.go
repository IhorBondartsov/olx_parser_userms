package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/IhorBondartsov/olx_parser_userms/cfg"
	"github.com/IhorBondartsov/olx_parser_userms/storage/userSQL"
	"github.com/IhorBondartsov/olx_parser_userms/webrpc"
	"time"
)

var log = logrus.New()

func main() {
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%s?timeout=5s",
		cfg.Storage.Login,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.DBName)
	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("[MAIN] Cant create db connection %v", err)
	}
	userStor := userSQL.NewUserMyClientMySQL(db)
	tokenStor := userSQL.NewTokenClientMySQL(db)

	apiCfg := webrpc.CfgAPI{
		AccessPublicKey:  []byte(cfg.PublicKey),
		AccessPrivateKey: []byte(cfg.PrivateKey),
		UserStor:         userStor,
		RefreshStor:      tokenStor,
		TTLAccessToken:   time.Duration(cfg.TTLAcessToken) * time.Second,
	}
	webrpc.Start(apiCfg)

	log.Info("Listening on ", (cfg.Route + ":" + cfg.Port))
	log.Panic(http.ListenAndServe((cfg.Route + ":" + cfg.Port), nil))
}
