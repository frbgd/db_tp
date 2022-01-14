package api

import (
	"db_tp/models"
	"db_tp/services"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func CreateUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	item := &models.User{}
	easyjson.Unmarshal(ctx.PostBody(), item)
	item.Nickname = nickname

	users, notUnique := services.UserSrv.CreateUser(item)
	if notUnique {
		resp, _ := easyjson.Marshal(users)
		ctx.Response.SetBody(resp)
		ctx.SetContentType("application/json")
		ctx.Response.SetStatusCode(409)
	}

	resp, _ := easyjson.Marshal(users[0])
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
	ctx.Response.SetStatusCode(201)
}

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

func EditUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	item := &models.UserUpdate{}
	easyjson.Unmarshal(ctx.PostBody(), item)

	user, notUnique := services.UserSrv.UpdateByNickname(nickname, item)
	if notUnique {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Not unique fields")}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(409)
		ctx.SetContentType("application/json")
		return
	}
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
