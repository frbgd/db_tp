package main

import (
	"db_tp/api"
	"db_tp/services"
	"db_tp/storage"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

var err error

func main() {
	storage.PgEngine, err = storage.NewPostgresDbEngine(
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBNAME"),
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
	)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer storage.PgEngine.Close()

	services.DatabaseSrv = services.NewDatabaseService(storage.PgEngine)
	services.ForumSrv = services.NewForumService(storage.PgEngine)
	services.PostSrv = services.NewPostService(storage.PgEngine)
	services.ThreadSrv = services.NewThreadService(storage.PgEngine)
	services.UserSrv = services.NewUserService(storage.PgEngine)

	r := router.New()

	r.POST("/api/forum/create", api.CreateForum)
	r.GET("/api/forum/{slug}/details", api.GetForumDetails)
	r.POST("/api/forum/{slug}/create", api.CreateThread)
	r.GET("/api/forum/{slug}/users", api.GetForumUsers)
	r.GET("/api/forum/{slug}/threads", api.GetForumThreads)

	r.GET("/api/post/{id}/details", api.GetPostDetails)
	r.POST("/api/post/{id}/details", api.EditPost)

	r.POST("/api/service/clear", api.ClearDatabase)
	r.GET("/api/service/status", api.GetDatabaseStatus)

	r.POST("/api/thread/{slug_or_id}/create", api.CreatePost)
	r.GET("/api/thread/{slug_or_id}/details", api.GetThreadDetails)
	r.POST("/api/thread/{slug_or_id}/details", api.EditThread)
	r.GET("/api/thread/{slug_or_id}/posts", api.GetThreadPosts)
	r.POST("/api/thread/{slug_or_id}/vote", api.VoteForThread)

	r.POST("/api/user/{nickname}/create", api.CreateUser)
	r.GET("/api/user/{nickname}/profile", api.GetUserDetails)
	r.POST("/api/user/{nickname}/profile", api.EditUser)

	log.Fatal(fasthttp.ListenAndServe(":5000", r.Handler))
}
