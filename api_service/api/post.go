package api

import (
	"db_tp/models"
	"db_tp/services"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

func GetPostDetails(ctx *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	related := string(ctx.QueryArgs().Peek("related"))

	post := services.PostSrv.GetById(id)

	if post == nil {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Can't find post with id:  %s", id)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}

	var response models.FullPost
	response.Post = post.Post
	if strings.Contains(related, "user") {
		response.Author = post.Author
	}
	if strings.Contains(related, "forum") {
		response.Forum = post.Forum
	}
	if strings.Contains(related, "thread") {
		response.Thread = post.Thread
	}

	resp, _ := easyjson.Marshal(response)
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
	ctx.Response.SetStatusCode(200)
}

//func EditPost(ctx *fasthttp.RequestCtx) {
//
//}
