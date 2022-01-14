package services

import (
	"context"
	"database/sql"
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
	rows, err := postSrv.db.CP.Query(
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
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil
	} else {
		fullPost := new(models.FullPost)

		var postObj models.Post
		parent := sql.NullInt64{}
		err = rows.Scan(&postObj.Id,
			&postObj.Thread,
			&postObj.Author,
			&postObj.Forum,
			&postObj.IsEdited,
			&postObj.Message,
			&parent,
			&postObj.Created,
		)
		if parent.Valid {
			postObj.Parent = int(parent.Int64)
		}
		if err != nil {
			panic(err)
		}

		fullPost.Post = &postObj
		fullPost.Author = UserSrv.GetByNickname(postObj.Author)
		fullPost.Forum = ForumSrv.GetBySlug(postObj.Forum)
		fullPost.Thread = ThreadSrv.GetById(postObj.Thread)

		return fullPost
	}
}

func (postSrv *PostService) UpdateById(id int, item *models.PostUpdate) *models.Post {

}
