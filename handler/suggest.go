package handler

import (
	"github.com/julienschmidt/httprouter"
	"kw_tool/types/requests"
	"kw_tool/util/protocol/httpresp"
	"net/http"
)

func (h *handler) suggest(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	req := &requests.Suggest{}
	req.Init(h.c, h.cli)
	if err := req.ValidateRequest(r); err != nil {
		httpresp.Send(http.StatusBadRequest, err, w)
	}
	if err := req.Run(); err != nil {
		httpresp.Send(http.StatusBadRequest, err, w)
	}
	httpresp.Send(http.StatusOK, "successful call", w)
}
