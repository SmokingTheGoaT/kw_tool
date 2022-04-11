package httpresp

import (
	"encoding/json"
	"kw_tool/util/constants"
	"net/http"
)

type response struct {
	Status      int
	ContentType string
	Data        []byte
	Header      http.Header
}

func (r response) write(w http.ResponseWriter) {
	if r.ContentType == constants.EmptyString {
		w.Header().Set("Content-Type", "application/json")
	} else {
		w.Header().Set("Content-Type", r.ContentType)
	}
	w.WriteHeader(r.Status)
	if _, err := w.Write(r.Data); err != nil {
		return
	}
}

func Send(status int, resp interface{}, w http.ResponseWriter) {
	res := response{}
	res.Status = status
	data, err := json.Marshal(resp)
	if err != nil {
		return
	}
	res.Data = data
	res.write(w)
}
