package services

import (
	"context"
	"database/sql"
	"db_tp/models"
	"db_tp/storage"
	"fmt"
	"github.com/jackc/pgx/v4"
	"strings"
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
	db *storage.PostgresDbEngine
}

func NewForumService(db *storage.PostgresDbEngine) *ForumService {
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

func (forumSrv *ForumService) CreateForum(item *models.Forum) (*models.Forum, bool) {
	row := forumSrv.db.CP.QueryRow(context.Background(),
		`INSERT INTO forums (slug, title, threads, posts, owner_nickname)
                    VALUES ($1, $2, $3, $4, (
                                SELECT nickname FROM users WHERE nickname = $5
                    )) RETURNING owner_nickname`,
		item.Slug,
		item.Title,
		item.Threads,
		item.Posts,
		item.User,
	)
	err := row.Scan(&item.User)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "unique") {
			return forumSrv.GetBySlug(item.Slug), true
		}
		if strings.Contains(fmt.Sprint(err), "not-null") {
			return nil, false
		}
	}
	return item, false
}

func (forumSrv *ForumService) CreateThread(item *models.Thread) (*models.Thread, bool) {
	var row pgx.Row
	if item.Slug == "" {
		row = forumSrv.db.CP.QueryRow(context.Background(),
			`INSERT INTO threads (user_nickname, title, message, created, forum_slug)
			VALUES ((SELECT nickname FROM users WHERE nickname = $1), $2, $3, $4,
					(SELECT slug FROM forums WHERE slug = $5))
			RETURNING id, user_nickname, forum_slug, created`,
			item.Author,
			item.Title,
			item.Message,
			item.Created,
			item.Forum,
		)
	} else {
		row = forumSrv.db.CP.QueryRow(context.Background(),
			`INSERT INTO threads (slug, user_nickname, title, message, created, forum_slug)
			VALUES ($1, (SELECT nickname FROM users WHERE nickname = $2), $3, $4, $5,
					(SELECT slug FROM forums WHERE slug = $6))
			RETURNING id, user_nickname, forum_slug, created`,
			item.Slug,
			item.Author,
			item.Title,
			item.Message,
			item.Created,
			item.Forum,
		)
	}

	err := row.Scan(&item.Id, &item.Author, &item.Forum, &item.Created)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "unique") {
			return ThreadSrv.GetBySlug(item.Slug), true
		}
		if strings.Contains(fmt.Sprint(err), "not-null") {
			return nil, false
		}
	}
	return item, false
}
