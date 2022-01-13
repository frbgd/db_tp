package api

import (
	"db_tp/models"
	"db_tp/services"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"strconv"
)

func CreatePost(ctx *fasthttp.RequestCtx) {

}

func GetThreadDetails(ctx *fasthttp.RequestCtx) {
	slugOrId := ctx.UserValue("slug_or_id").(string)

	thread := services.ThreadSrv.GetBySlugOrId(slugOrId)

	if thread == nil {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Can't find thread with slug or id:  %s", slugOrId)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}

	resp, _ := easyjson.Marshal(thread)
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
	ctx.Response.SetStatusCode(200)
}

func EditThread(ctx *fasthttp.RequestCtx) {

}

func GetThreadPosts(ctx *fasthttp.RequestCtx) {
	slugOrId := ctx.UserValue("slug_or_id").(string)
	limit := string(ctx.QueryArgs().Peek("limit"))
	since := string(ctx.QueryArgs().Peek("since"))
	desc := string(ctx.QueryArgs().Peek("desc"))
	sort := string(ctx.QueryArgs().Peek("sort"))

	var thread *models.Thread
	id, err := strconv.Atoi(slugOrId)
	if err == nil {
		thread = services.ThreadSrv.GetById(id)
	} else {
		thread = services.ThreadSrv.GetBySlug(slugOrId)
	}

	if thread == nil {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Can't find posts for thread with slug or id:  %s", slugOrId)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}

	posts := services.ThreadSrv.GetPosts(thread.Id, desc, limit, since, sort)
	resp, _ := easyjson.Marshal(posts)
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
	ctx.Response.SetStatusCode(200)
}

func VoteForThread(ctx *fasthttp.RequestCtx) {

}