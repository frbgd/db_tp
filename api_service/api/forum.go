package api

import (
	"db_tp/models"
	"db_tp/services"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

func CreateForum(ctx *fasthttp.RequestCtx) {
	item := &models.Forum{}
	easyjson.Unmarshal(ctx.PostBody(), item)

	forum, notUnique := services.ForumSrv.CreateForum(item)
	if forum == nil {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Can't find user with nickname:  %s", item.User)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}
	if notUnique {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Not unique fields")}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(409)
		ctx.SetContentType("application/json")
		return
	}

	resp, _ := easyjson.Marshal(forum)
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
	ctx.Response.SetStatusCode(201)
}

func GetForumDetails(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	forum := services.ForumSrv.GetBySlug(slug)

	if forum == nil {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Can't find forum with slug:  %s", slug)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}

	resp, _ := easyjson.Marshal(forum)
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
	ctx.Response.SetStatusCode(200)
}

func CreateThread(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	item := &models.Thread{}
	easyjson.Unmarshal(ctx.PostBody(), item)
	if item.Forum == "" {
		item.Forum = slug
	}

	thread, notUnique := services.ForumSrv.CreateThread(item)
	if thread == nil {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Can't find")}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}
	if notUnique {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Not unique fields")}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(409)
		ctx.SetContentType("application/json")
		return
	}

	resp, _ := easyjson.Marshal(thread)
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
	ctx.Response.SetStatusCode(201)
}

func GetForumUsers(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	desc, _ := strconv.ParseBool(string(ctx.QueryArgs().Peek("desc")))
	limitStr := string(ctx.QueryArgs().Peek("limit"))
	var limit int
	if limitStr == "" {
		limit = 100
	} else {
		limit, _ = strconv.Atoi(limitStr)
	}
	since := string(ctx.QueryArgs().Peek("since"))

	users := services.ForumSrv.GetForumUsers(slug, desc, limit, since)

	if users == nil {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Can't find users for forum with slug:  %s", slug)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}

	resp, _ := easyjson.Marshal(users)
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
	ctx.Response.SetStatusCode(200)
}

func GetForumThreads(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	desc, _ := strconv.ParseBool(string(ctx.QueryArgs().Peek("desc")))
	limitStr := string(ctx.QueryArgs().Peek("limit"))
	var limit int
	if limitStr == "" {
		limit = 100
	} else {
		limit, _ = strconv.Atoi(limitStr)
	}

	var threads models.Threads

	since := string(ctx.QueryArgs().Peek("since"))
	if since != "" {
		sinceTime, _ := time.Parse(time.RFC3339, since)
		threads = services.ForumSrv.GetForumThreads(slug, desc, limit, &sinceTime)
	} else {
		threads = services.ForumSrv.GetForumThreads(slug, desc, limit, nil)
	}

	if threads == nil {
		errMsg := models.ErrMsg{Message: fmt.Sprintf("Can't find threads for forum with slug:  %s", slug)}
		response, _ := easyjson.Marshal(errMsg)
		ctx.SetBody(response)
		ctx.SetStatusCode(404)
		ctx.SetContentType("application/json")
		return
	}

	resp, _ := easyjson.Marshal(threads)
	ctx.Response.SetBody(resp)
	ctx.SetContentType("application/json")
	ctx.Response.SetStatusCode(200)
}
