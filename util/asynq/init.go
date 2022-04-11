package asynq

import (
	"github.com/hibiken/asynq"
	"kw_tool/tasks"
)

type Server struct {
	srv *asynq.Server
	mux *asynq.ServeMux
}

func New(opt asynq.RedisClientOpt, cfs asynq.Config) *Server {
	return &Server{
		srv: asynq.NewServer(opt, cfs),
		mux: asynq.NewServeMux(),
	}
}

func (s *Server) Init() (err error) {
	s.mux.HandleFunc(tasks.TypeRecursiveCrawlRequest, tasks.HandleRecursiveCrawlTask)
	err = s.srv.Run(s.mux)
	return
}
