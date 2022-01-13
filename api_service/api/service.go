package api

import (
	"db_tp/services"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func ClearDatabase(ctx *fasthttp.RequestCtx) {

}

func GetDatabaseStatus(ctx *fasthttp.RequestCtx) {
	status := services.DatabaseSrv.GetStatus()

	resp, _ := easyjson.Marshal(status)
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
	ctx.Response.SetStatusCode(200)
}
