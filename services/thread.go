package services

import (
	"context"
	"database/sql"
	"db_tp/models"
	"db_tp/storage"
	"fmt"
	"github.com/jackc/pgx/v4"
	"strconv"
	"strings"
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
	db *storage.PostgresDbEngine
}

func NewThreadService(db *storage.PostgresDbEngine) *ThreadService {
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
	rows, err := threadSrv.db.CP.Query(
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
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		thread := new(models.Thread)
		err = rows.Scan(
			&thread.Id,
			&thread.Slug,
			&thread.Forum,
			&thread.Author,
			&thread.Title,
			&thread.Message,
			&thread.Votes,
			&thread.Created,
		)
		if err != nil {
			panic(err)
		}
		return thread
	} else {
		return nil
	}
}

func (threadSrv *ThreadService) GetById(id int) *models.Thread {
	rows, err := threadSrv.db.CP.Query(
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
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		thread := new(models.Thread)
		slug := sql.NullString{}
		err = rows.Scan(
			&thread.Id,
			&slug,
			&thread.Forum,
			&thread.Author,
			&thread.Title,
			&thread.Message,
			&thread.Votes,
			&thread.Created,
		)
		if err != nil {
			panic(err)
		}
		if slug.Valid {
			thread.Slug = slug.String
		}
		return thread
	} else {
		return nil
	}
}

func (threadSrv *ThreadService) CreatePosts(posts *models.Posts) (*models.Posts, bool) {
	resultPosts := make(models.Posts, 0)
	for _, post := range *posts {
		var row pgx.Row
		if post.Parent == 0 {
			row = threadSrv.db.CP.QueryRow(
				context.Background(),
				`INSERT INTO posts (thread_id, user_nickname, forum_slug, message, parent, created)
                VALUES ($1, $2, $3, $4, NULL, $5)
                RETURNING id, thread_id, user_nickname, forum_slug, is_edited, message, parent, created`,
				post.Thread,
				post.Author,
				post.Forum,
				post.Message,
				post.Created,
			)
		} else {
			row = threadSrv.db.CP.QueryRow(
				context.Background(),
				`INSERT INTO posts (thread_id, user_nickname, forum_slug, message, parent, created)
                VALUES ($1, $2, $3, $4, $5, $6)
                RETURNING id, thread_id, user_nickname, forum_slug, is_edited, message, parent, created`,
				post.Thread,
				post.Author,
				post.Forum,
				post.Message,
				post.Parent,
				post.Created,
			)
		}

		newPost := &models.Post{}
		parent := sql.NullInt64{}
		err := row.Scan(
			&newPost.Id,
			&newPost.Thread,
			&newPost.Author,
			&newPost.Forum,
			&newPost.IsEdited,
			&newPost.Message,
			&parent,
			&newPost.Created,
		)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "not-null") || strings.Contains(fmt.Sprint(err), "foreign key") {
				return nil, false
			}
			if strings.Contains(fmt.Sprint(err), "not exists") {
				return nil, true
			}
		}
		if parent.Valid {
			newPost.Parent = int(parent.Int64)
		}
		resultPosts = append(resultPosts, *newPost)
	}
	return &resultPosts, false
}

func (threadSrv *ThreadService) UpdateBySlugOrId(slugOrId string, item *models.ThreadUpdate) *models.Thread {
	var thread *models.Thread = nil
	if id, err := strconv.Atoi(slugOrId); err == nil {
		thread = threadSrv.UpdateById(id, item)
	}
	if thread != nil {
		return thread
	}
	return threadSrv.UpdateBySlug(slugOrId, item)
}

func (threadSrv *ThreadService) UpdateBySlug(slug string, item *models.ThreadUpdate) *models.Thread {
	threadTitle := new(string)
	if item.Title != "" {
		threadTitle = &item.Title
	}
	threadMessage := new(string)
	if item.Message != "" {
		threadMessage = &item.Message
	}

	thread := &models.Thread{}
	row := threadSrv.db.CP.QueryRow(context.Background(),
		`UPDATE threads
			SET title   = COALESCE(NULLIF($2, ''), title),
				message = COALESCE(NULLIF($3, ''), message)
			WHERE slug = $1
			RETURNING id, slug, forum_slug, user_nickname, title, message, votes, created`,
		slug, threadTitle, threadMessage)
	slugFromDb := sql.NullString{}
	err := row.Scan(
		&thread.Id,
		&slugFromDb,
		&thread.Forum,
		&thread.Author,
		&thread.Title,
		&thread.Message,
		&thread.Votes,
		&thread.Created,
	)
	if err == pgx.ErrNoRows {
		return nil
	}
	if slugFromDb.Valid {
		thread.Slug = slugFromDb.String
	}
	return thread
}

func (threadSrv *ThreadService) UpdateById(id int, item *models.ThreadUpdate) *models.Thread {
	threadTitle := new(string)
	if item.Title != "" {
		threadTitle = &item.Title
	}
	threadMessage := new(string)
	if item.Message != "" {
		threadMessage = &item.Message
	}

	thread := &models.Thread{}
	row := threadSrv.db.CP.QueryRow(context.Background(),
		`UPDATE threads
			SET title   = COALESCE(NULLIF($2, ''), title),
				message = COALESCE(NULLIF($3, ''), message)
			WHERE id = $1
			RETURNING id, slug, forum_slug, user_nickname, title, message, votes, created`,
		id, threadTitle, threadMessage)
	slugFromDb := sql.NullString{}
	err = row.Scan(
		&thread.Id,
		&slugFromDb,
		&thread.Forum,
		&thread.Author,
		&thread.Title,
		&thread.Message,
		&thread.Votes,
		&thread.Created,
	)
	if err == pgx.ErrNoRows {
		return nil
	}
	if slugFromDb.Valid {
		thread.Slug = slugFromDb.String
	}
	return thread
}

func (threadSrv *ThreadService) VoteById(id int, item *models.Vote) error {
	_, err := threadSrv.db.CP.Exec(context.Background(),
		`INSERT INTO votes (thread_id, user_nickname, voice)
			VALUES ($1, $2, $3)
			ON CONFLICT (thread_id, user_nickname) DO UPDATE SET voice = $3`,
		id,
		item.Nickname,
		item.Voice,
	)
	return err
}

func (threadSrv *ThreadService) GetPosts(threadId int, desc bool, limit int, since int, sort string) models.Posts {
	var rows pgx.Rows
	if sort == "" {
		sort = "flat"
	}

	if since != 0 {
		rows, err = threadSrv.db.CP.Query(context.Background(),
			sqlGetSortedPostsSince[desc][sort],
			threadId,
			since,
			limit,
		)
		if err != nil {
			panic(err)
		}
	} else {
		rows, err = threadSrv.db.CP.Query(context.Background(),
			sqlGetSortedPosts[desc][sort],
			threadId,
			limit,
		)
		if err != nil {
			panic(err)
		}
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		var post models.Post
		parent := sql.NullInt64{}
		err = rows.Scan(
			&post.Id,
			&post.Thread,
			&post.Author,
			&post.Forum,
			&post.IsEdited,
			&post.Message,
			&parent,
			&post.Created,
		)
		if err != nil {
			panic(err)
		}
		if parent.Valid {
			post.Parent = int(parent.Int64)
		}
		posts = append(posts, post)
	}
	return posts
}
