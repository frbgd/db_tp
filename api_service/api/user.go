package api

import (
	"db_tp/models"
	"db_tp/services"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

//func CreateUser(ctx *fasthttp.RequestCtx) {
//
//}

func GetUserDetails(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	user := services.UserSrv.GetByNickname(nickname)

	if user == nil {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Can't find user with nickname:  %s", nickname)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}

	resp, _ := easyjson.Marshal(user)
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
	ctx.Response.SetStatusCode(200)
}

//func EditUser(ctx *fasthttp.RequestCtx) {
//
//}
