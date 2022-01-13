package services

import (
	"context"
	"db_tp/db"
	"db_tp/models"
	"github.com/jackc/pgx/v4"
	"strconv"
)

var sqlGetSortedPostsSince = map[bool]map[string]string{
	true: {
		"flat": `SELECT id,
				   thread_id,
				   user_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND id < $2
			ORDER BY id DESC
			LIMIT $3`,
		"tree": `SELECT id,
				   thread_id,
				   user_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND path < (SELECT path FROM posts WHERE id = $2)
			ORDER BY path DESC
			LIMIT $3`,
		"parent_tree": `WITH roots AS (
				SELECT DISTINCT path[1]
				FROM posts
				WHERE thread_id = $1
				  AND parent IS NULL
				  AND path[1] < (SELECT path[1] FROM posts WHERE id = $2)
         		ORDER BY path[1] DESC
				LIMIT $3
			)
			SELECT id,
				   thread_id,
				   user_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND path[1] IN (SELECT * FROM roots)
			ORDER BY path[1] DESC, path[2:]`,
	},
	false: {
		"flat": `SELECT id,
				   thread_id,
				   user_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND id > $2
			ORDER BY id
			LIMIT $3`,
		"tree": `SELECT id,
				   thread_id,
				   user_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND path > (SELECT path FROM posts WHERE id = $2)
			ORDER BY path
			LIMIT $3`,
		"parent_tree": `WITH roots AS (
				SELECT DISTINCT path[1]
				FROM posts
				WHERE thread_id = $1
				  AND parent IS NULL
				  AND path[1] > (SELECT path[1] FROM posts WHERE id = $2)
         		ORDER BY path[1]
				LIMIT $3
			)
			SELECT id,
				   thread_id,
				   user_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND path[1] IN (SELECT * FROM roots)
			ORDER BY path`,
	},
}

var sqlGetSortedPosts = map[bool]map[string]string{
	true: {
		"flat": `SELECT id,
				   thread_id,
				   user_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			ORDER BY id DESC
			LIMIT $2`,
		"tree": `SELECT id,
				   thread_id,
				   user_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			ORDER BY path DESC
			LIMIT $2`,
		"parent_tree": `WITH roots AS (
				SELECT DISTINCT path[1]
				FROM posts
				WHERE thread_id = $1
				ORDER BY path[1] DESC
				LIMIT $2
			)
			SELECT id,
				   thread_id,
				   user_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND path[1] IN (SELECT * FROM roots)
			ORDER BY path[1] DESC, path[2:]`,
	},
	false: {
		"flat": `SELECT id,
				   thread_id,
				   user_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			ORDER BY id
			LIMIT $2`,
		"tree": `SELECT id,
					   thread_id,
					   user_nickname,
					   forum_slug,
					   is_edited,
					   message,
					   parent,
					   created
				FROM posts
				WHERE thread_id = $1
				ORDER BY path
				LIMIT $2`,
		"parent_tree": `WITH roots AS (
				SELECT DISTINCT path[1]
				FROM posts
				WHERE thread_id = $1
				ORDER BY path[1]
				LIMIT $2
			)
			SELECT id,
				   thread_id,
				   user_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND path[1] IN (SELECT * FROM roots)
			ORDER BY path`,
	},
}

var ThreadSrv *ThreadService

type ThreadService struct {
	db *db.PostgresDbEngine
}

func NewThreadService(db *db.PostgresDbEngine) *ThreadService {
	srv := new(ThreadService)
	srv.db = db
	return srv
}

func (threadSrv *ThreadService) GetBySlugOrId(slugOrId string) *models.Thread {
	var threadObj *models.Thread
	if id, err := strconv.Atoi(slugOrId); err == nil {
		threadObj = threadSrv.GetById(id)
	}
	if threadObj != nil {
		return threadObj
	}
	return threadSrv.GetBySlug(slugOrId)
}

func (threadSrv *ThreadService) GetBySlug(slug string) *models.Thread {
	rows, _ := threadSrv.db.CP.Query(
		context.Background(),
		`SELECT id,
				   slug,
				   forum_slug,
				   user_nickname,
				   title,
				   message,
				   votes,
				   created
			FROM threads
			WHERE slug = $1`,
		slug)
	defer rows.Close()

	if rows.Next() {
		thread := new(models.Thread)
		rows.Scan(
			&thread.Id,
			&thread.Slug,
			&thread.Forum,
			&thread.Author,
			&thread.Title,
			&thread.Message,
			&thread.Votes,
			&thread.Created,
		)
		return thread
	} else {
		return nil
	}
}

func (threadSrv *ThreadService) GetById(id int) *models.Thread {
	rows, _ := threadSrv.db.CP.Query(
		context.Background(),
		`SELECT id,
				   slug,
				   forum_slug,
				   user_nickname,
				   title,
				   message,
				   votes,
				   created
			FROM threads
			WHERE id = $1`,
		id)
	defer rows.Close()

	if rows.Next() {
		thread := new(models.Thread)
		rows.Scan(
			&thread.Id,
			&thread.Slug,
			&thread.Forum,
			&thread.Author,
			&thread.Title,
			&thread.Message,
			&thread.Votes,
			&thread.Created,
		)
		return thread
	} else {
		return nil
	}
}

func (threadSrv *ThreadService) CreatePosts(posts *models.Post[]) (*models.Post, error) {

}

func (threadSrv *ThreadService) UpdateBySlugOrId(slugOrId string, item *models.ThreadUpdate) (*models.Thread, error) {

}

func (threadSrv *ThreadService) UpdateBySlug(slug string, item *models.ThreadUpdate) (*models.Thread, error) {

}

func (threadSrv *ThreadService) UpdateById(id int64, item *models.ThreadUpdate) (*models.Thread, error) {

}

func (threadSrv *ThreadService) VoteById(id int64, item *models.Vote) error {

}

func (threadSrv *ThreadService) GetPosts(threadId int, desc bool, limit int, since *int, sort string) []models.Post {
	var rows pgx.Rows
	if since != nil {
		rows, _ = threadSrv.db.CP.Query(context.Background(),
			sqlGetSortedPostsSince[desc][sort],
			threadId,
			&since,
			limit,
		)
	} else {
		rows, _ = threadSrv.db.CP.Query(context.Background(),
			sqlGetSortedPosts[desc][sort],
			threadId,
			limit,
		)
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		var post models.Post
		rows.Scan(
			&post.Id,
			&post.Thread,
			&post.Author,
			&post.Forum,
			&post.IsEdited,
			&post.Message,
			&post.Parent,
			&post.Created,
		)
		posts = append(posts, post)
	}
	return posts
}
