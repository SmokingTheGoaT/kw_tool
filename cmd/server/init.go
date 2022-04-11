package server

import (
	"kw_tool/cmd/server/app"
)

func Start() (err error) {
	application := app.KWTool{}
	err = application.Init()
	return
}
