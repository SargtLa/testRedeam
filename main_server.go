package main

import (
	"flag"
	"net"
	"path"
	"time"

	"github.com/ruslanBik4/httpgo/apis"
	httpgo "github.com/ruslanBik4/httpgo/httpGo"
	"github.com/ruslanBik4/logs"

	"github.com/SargtLa/testRedeam/api"
	"github.com/SargtLa/testRedeam/db"
)

var (
	httpServer *httpgo.HttpGo

	flagfPort      = flag.String("port", ":8080", "host address to listen on")
	flagSystemPath = flag.String("path", "./", "path to system files")
	flagCfgPath    = flag.String("config_path", "cfg", "path to config files")
	flagServerCfg  = flag.String("server_config", "httpgo.yml", "server config yml file name")
)

func init() {
	flag.Parse()
	listener, err := net.Listen("tcp", *flagfPort)
	if err != nil {
		logs.Fatal(err)
	}

	cfg, err := httpgo.NewCfgHttp(path.Join(*flagSystemPath, *flagCfgPath, *flagServerCfg))
	if err != nil || cfg == nil {
		logs.Fatal(err, cfg)
	}

	ctxApis := apis.NewCtxApis(4)

	ctxApis.AddValue("migration", path.Join(*flagCfgPath, "DB"))
	DB := db.GetDB(ctxApis)
	if DB == nil {
		panic("cannot init DB")
	}

	ctxApis.AddValue("DB", DB)
	ctxApis.AddValue(api.CFG_PATH, *flagCfgPath)
	ctxApis.AddValue(api.SYSTEM_PATH, *flagSystemPath)

	var auth apis.FncAuth

	a := apis.NewApis(ctxApis, api.Routes, auth)

	httpServer = httpgo.NewHttpgo(cfg, listener, a)
}

func main() {
	logs.StatusLog("server starting %s on port %s", time.Now(), *flagfPort)

	defer func() {
		errRec := recover()
		if err, ok := errRec.(error); ok {
			logs.ErrorStack(err)
		}
	}()

	err := httpServer.Run(
		false,
		"",
		"")
	if err != nil {
		logs.ErrorStack(err)
	} else {
		logs.StatusLog("Server https correct shutdown at %v", time.Now())
	}

}
