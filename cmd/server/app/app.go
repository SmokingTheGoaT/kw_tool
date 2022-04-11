package app

import (
	"github.com/hibiken/asynq"
	"kw_tool/handler"
	local "kw_tool/util/asynq"
	"kw_tool/util/constants"
	"kw_tool/util/protocol/http2"
)

type KWTool struct {
	server       *http2.Server
	client       *local.Client
	schedulerSrv *local.Server
}

func (k *KWTool) Init() (err error) {
	redisOpts := asynq.RedisClientOpt{Addr: constants.RedisAddr}
	asynqCfs := asynq.Config{Concurrency: 10}
	client := &local.Client{}
	client.Init(redisOpts)
	k.client = client
	defer func(client *asynq.Client) {
		err = client.Close()
		if err != nil {

		}
	}(k.client.Cli)
	k.schedulerSrv = local.New(redisOpts, asynqCfs)
	if err = k.schedulerSrv.Init(); err != nil {
		return
	}
	k.server = http2.NewServer(http2.Config{ListenAddr: constants.ServerAddr})
	hdlr := handler.New(k.server, k.client)
	hdlr.Init()
	err = k.server.Listen()
	return
}
