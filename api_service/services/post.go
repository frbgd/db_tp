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

func (forum *PostService) GetById(id int) *models.FullPost {
	rows, _ := forum.db.CP.Query(
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

		var post models.Post
		rows.Scan(&post.Thread,
			&post.Author,
			&post.Forum,
			&post.IsEdited,
			&post.Message,
			&post.Parent,
			&post.Created,
		)

		fullPost.Post = post
		fullPost.Author = *UserSrv.GetByNickname(post.Author)
		fullPost.Forum = *ForumSrv.GetBySlug(post.Forum)
		fullPost.Thread = *ThreadSrv.GetById(post.Thread)

		return fullPost
	}
}

func (forum *PostService) UpdateById(id int64) (*models.PostUpdate, error) {

}
