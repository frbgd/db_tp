package main

import (
	api2 "db_tp/internal/api"
	services2 "db_tp/internal/services"
	"db_tp/internal/storage"
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

	services2.DatabaseSrv = services2.NewDatabaseService(storage.PgEngine)
	services2.ForumSrv = services2.NewForumService(storage.PgEngine)
	services2.PostSrv = services2.NewPostService(storage.PgEngine)
	services2.ThreadSrv = services2.NewThreadService(storage.PgEngine)
	services2.UserSrv = services2.NewUserService(storage.PgEngine)

	r := router.New()

	r.POST("/api/forum/create", api2.CreateForum)
	r.GET("/api/forum/{slug}/details", api2.GetForumDetails)
	r.POST("/api/forum/{slug}/create", api2.CreateThread)
	r.GET("/api/forum/{slug}/users", api2.GetForumUsers)
	r.GET("/api/forum/{slug}/threads", api2.GetForumThreads)

	r.GET("/api/post/{id}/details", api2.GetPostDetails)
	r.POST("/api/post/{id}/details", api2.EditPost)

	r.POST("/api/service/clear", api2.ClearDatabase)
	r.GET("/api/service/status", api2.GetDatabaseStatus)

	r.POST("/api/thread/{slug_or_id}/create", api2.CreatePost)
	r.GET("/api/thread/{slug_or_id}/details", api2.GetThreadDetails)
	r.POST("/api/thread/{slug_or_id}/details", api2.EditThread)
	r.GET("/api/thread/{slug_or_id}/posts", api2.GetThreadPosts)
	r.POST("/api/thread/{slug_or_id}/vote", api2.VoteForThread)

	r.POST("/api/user/{nickname}/create", api2.CreateUser)
	r.GET("/api/user/{nickname}/profile", api2.GetUserDetails)
	r.POST("/api/user/{nickname}/profile", api2.EditUser)

	log.Fatal(fasthttp.ListenAndServe(":5000", r.Handler))
}
