package handler

import (
	"github.com/hibiken/asynq"
	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
	"kw_tool/util/constants"
	"kw_tool/util/protocol/http2"
	"kw_tool/util/protocol/httpresp"
	"net/http"
	"time"
)

type handler struct {
	s   *http2.Server
	c   *cache.Cache
	cli *asynq.Client
}

func New(s *http2.Server, cli *asynq.Client) *handler {
	return &handler{
		s:   s,
		c:   cache.New(cache.NoExpiration, 10*time.Minute),
		cli: cli,
	}
}

func (h *handler) Init() {
	h.s.RegisterHandler(http.MethodGet, constants.HealthzURI, h.healthz)
	h.s.RegisterHandler(http.MethodPost, constants.SuggestURI, h.suggest)
}

func (h *handler) healthz(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	httpresp.Send(http.StatusOK, "health check successful", w)
}
