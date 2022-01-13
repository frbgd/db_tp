package services

import (
	"context"
	"db_tp/db"
	"db_tp/models"
)

var PostSrv *PostService

type PostService struct {
	db *db.PostgresDbEngine
}

func NewPostService(db *db.PostgresDbEngine) *PostService {
	srv := new(PostService)
	srv.db = db
	return srv
}

func (postSrv *PostService) GetById(id int) *models.FullPost {
	rows, _ := postSrv.db.CP.Query(
		context.Background(),
		`SELECT id,
				   thread_id,
				   user_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE id = $1`,
		id)
	defer rows.Close()

	if !rows.Next() {
		return nil
	} else {
		fullPost := new(models.FullPost)

		var postObj models.Post
		rows.Scan(&postObj.Thread,
			&postObj.Author,
			&postObj.Forum,
			&postObj.IsEdited,
			&postObj.Message,
			&postObj.Parent,
			&postObj.Created,
		)

		fullPost.Post = postObj
		fullPost.Author = *UserSrv.GetByNickname(postObj.Author)
		fullPost.Forum = *ForumSrv.GetBySlug(postObj.Forum)
		fullPost.Thread = *ThreadSrv.GetById(postObj.Thread)

		return fullPost
	}
}

//func (postSrv *PostService) UpdateById(id int64) (*models.PostUpdate, error) {
//
//}
