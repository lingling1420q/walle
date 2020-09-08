package main

import (
	"fmt"
	"github.com/Gitforxuyang/walle/config"
	"github.com/Gitforxuyang/walle/server"
	"github.com/Gitforxuyang/walle/util/logger"
	"github.com/Gitforxuyang/walle/util/sentry"
	"github.com/Gitforxuyang/walle/util/trace"
)

func main() {
	config.Init()
	conf := config.GetConfig()
	logger.Init(conf.GetName())
	trace.Init(fmt.Sprintf("%s_%s", conf.GetName(), conf.GetENV()), conf.GetTraceConfig().Endpoint, conf.GetTraceConfig().Ratio)
	sentry.Init()
	server.InitServer()
}
