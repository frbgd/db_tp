package api

import (
	"db_tp/internal/models"
	"db_tp/internal/services"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
	"time"
)

func CreatePost(ctx *fasthttp.RequestCtx) {
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

	items := &models.Posts{}
	easyjson.Unmarshal(ctx.PostBody(), items)

	if items != nil && len(*items) == 0 {
		resp, _ := easyjson.Marshal(items)
		ctx.Response.SetBody(resp)
		ctx.SetContentType("application/json")
		ctx.Response.SetStatusCode(201)
		return
	}

	now := time.Now()
	input := make(models.Posts, 0)
	for _, item := range *items {

		input = append(input, models.Post{
			Message: item.Message,
			Created: now,
			Author:  item.Author,
			Forum:   thread.Forum,
			Parent:  item.Parent,
			Thread:  thread.Id,
		})
	}

	posts, invalidParents := services.ThreadSrv.CreatePosts(&input)
	if invalidParents {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Invalid parents")}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(409)
		ctx.SetContentType("application/json")
		return
	}
	if posts == nil {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Can't find user")}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}

	resp, _ := easyjson.Marshal(posts)
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
	ctx.Response.SetStatusCode(201)
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
	slugOrId := ctx.UserValue("slug_or_id").(string)
	item := &models.ThreadUpdate{}
	easyjson.Unmarshal(ctx.PostBody(), item)

	thread := services.ThreadSrv.UpdateBySlugOrId(slugOrId, item)
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

func GetThreadPosts(ctx *fasthttp.RequestCtx) {
	slugOrId := ctx.UserValue("slug_or_id").(string)
	limitStr := string(ctx.QueryArgs().Peek("limit"))
	var limit int
	if limitStr == "" {
		limit = 100
	} else {
		limit, _ = strconv.Atoi(limitStr)
	}
	since, _ := strconv.Atoi(string(ctx.QueryArgs().Peek("since")))
	desc, _ := strconv.ParseBool(string(ctx.QueryArgs().Peek("desc")))
	sort := string(ctx.QueryArgs().Peek("sort"))

	var thread *models.Thread
	id, err := strconv.Atoi(slugOrId)
	if err == nil {
		thread = services.ThreadSrv.GetById(id)
	}
	if thread == nil {
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
	slugOrId := ctx.UserValue("slug_or_id").(string)
	item := &models.Vote{}
	easyjson.Unmarshal(ctx.PostBody(), item)

	if id, err := strconv.Atoi(slugOrId); err == nil {
		voteErr := services.ThreadSrv.VoteById(id, item)
		if voteErr == nil {
			thread := services.ThreadSrv.GetById(id)
			resp, _ := easyjson.Marshal(thread)
			ctx.Response.SetBody(resp)
			ctx.SetContentType("application/json")
			ctx.Response.SetStatusCode(200)
			return
		}
	}

	thread := services.ThreadSrv.GetBySlug(slugOrId)
	if thread == nil {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Can't find posts for thread with slug or id:  %s", slugOrId)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}
	voteErr := services.ThreadSrv.VoteById(thread.Id, item)
	if strings.Contains(fmt.Sprint(voteErr), "foreign key") {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Can't find posts for thread with slug or id:  %s", slugOrId)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}
	thread = services.ThreadSrv.GetById(thread.Id)
	if thread == nil {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Can't find posts for thread with slug or id:  %s", slugOrId)}
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
