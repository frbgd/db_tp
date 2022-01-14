package services

import (
	"context"
	"database/sql"
	"db_tp/db"
	"db_tp/models"
	"github.com/jackc/pgx/v4"
	"time"
)

var err error

var sqlGetForumUserWithSince = map[bool]string{
	true: `SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		  AND user_nickname < $2
		ORDER BY user_nickname DESC
		LIMIT $3`,

	false: `SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		  AND user_nickname > $2
		ORDER BY user_nickname
		LIMIT $3`,
}

var sqlGetForumUser = map[bool]string{
	true: `SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		ORDER BY user_nickname DESC
		LIMIT $2`,

	false: `SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		ORDER BY user_nickname
		LIMIT $2`,
}

var sqlGetThreadsByForumSlugSince = map[bool]string{
	true: `SELECT id,
			   slug,
			   forum_slug,
			   user_nickname,
			   title,
			   message,
			   votes,
			   created
		FROM threads
		WHERE forum_slug = $1 AND created <= $2
		ORDER BY created DESC
		LIMIT $3`,

	false: `SELECT id,
			   slug,
			   forum_slug,
			   user_nickname,
			   title,
			   message,
			   votes,
			   created
		FROM threads
		WHERE forum_slug = $1 AND created >= $2
		ORDER BY created
		LIMIT $3`,
}

var sqlGetThreadsByForumSlug = map[bool]string{
	true: `SELECT id,
			   slug,
			   forum_slug,
			   user_nickname,
			   title,
			   message,
			   votes,
			   created
		FROM threads
		WHERE forum_slug = $1
		ORDER BY created DESC
		LIMIT $2`,

	false: `SELECT id,
			   slug,
			   forum_slug,
			   user_nickname,
			   title,
			   message,
			   votes,
			   created
		FROM threads
		WHERE forum_slug = $1
		ORDER BY created
		LIMIT $2`,
}

var ForumSrv *ForumService

type ForumService struct {
	db *db.PostgresDbEngine
}

func NewForumService(db *db.PostgresDbEngine) *ForumService {
	srv := new(ForumService)
	srv.db = db
	return srv
}

func (forumSrv *ForumService) GetBySlug(slug string) *models.Forum {
	rows, err := forumSrv.db.CP.Query(
		context.Background(),
		`SELECT slug, title, threads, posts, owner_nickname 
FROM forums 
WHERE slug = $1`,
		slug)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		forum := new(models.Forum)
		err = rows.Scan(&forum.Slug, &forum.Title, &forum.Threads, &forum.Posts, &forum.User)
		if err != nil {
			panic(err)
		}
		return forum
	} else {
		return nil
	}
}

//func (forumSrv *ForumService) GetThreadBySlug(slug string) (*models.Thread, error) {
//
//}

func (forumSrv *ForumService) GetForumUsers(slug string, desc bool, limit int, since string) models.Users {
	var rows pgx.Rows
	if since != "" {
		rows, _ = forumSrv.db.CP.Query(context.Background(),
			sqlGetForumUserWithSince[desc],
			slug,
			since,
			limit,
		)
	} else {
		rows, _ = forumSrv.db.CP.Query(context.Background(),
			sqlGetForumUser[desc],
			slug,
			limit,
		)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	foundUsers := false
	for rows.Next() {
		foundUsers = true
		var user models.User
		rows.Scan(&user.Nickname, &user.Email, &user.Fullname, &user.About)
		users = append(users, user)
	}

	if !foundUsers {
		forumObj := forumSrv.GetBySlug(slug)
		if forumObj == nil {
			return nil
		}
	}

	return users
}

func (forumSrv *ForumService) GetForumThreads(slug string, desc bool, limit int, since *time.Time) models.Threads {
	var rows pgx.Rows

	if since != nil {
		sql := sqlGetThreadsByForumSlugSince[desc]
		rows, err = forumSrv.db.CP.Query(context.Background(),
			sql,
			slug,
			since,
			limit,
		)
		if err != nil {
			panic(err)
		}
	} else {
		sql := sqlGetThreadsByForumSlug[desc]
		rows, err = forumSrv.db.CP.Query(context.Background(),
			sql,
			slug,
			limit,
		)
		if err != nil {
			panic(err)
		}
	}
	defer rows.Close()

	threads := make([]models.Thread, 0)
	foundThreads := false
	for rows.Next() {
		foundThreads = true
		var thread models.Thread
		slug := sql.NullString{}
		err = rows.Scan(&thread.Id,
			&slug,
			&thread.Forum,
			&thread.Author,
			&thread.Title,
			&thread.Message,
			&thread.Votes,
			&thread.Created)
		if err != nil {
			panic(err)
		}
		if slug.Valid {
			thread.Slug = slug.String
		}
		threads = append(threads, thread)
	}

	if !foundThreads {
		forumObj := forumSrv.GetBySlug(slug)
		if forumObj == nil {
			return nil
		}
	}

	return threads
}

//func (forumSrv *ForumService) CreateForum(item *models.Forum) (*models.Forum, error) {
//
//}
//
//func (forumSrv *ForumService) CreateThread(item *models.Thread) (*models.Thread, error) {
//
//}
