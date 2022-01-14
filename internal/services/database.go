package services

import (
	"context"
	"db_tp/internal/models"
	"db_tp/internal/storage"
)

var DatabaseSrv *DatabaseService

type DatabaseService struct {
	db *storage.PostgresDbEngine
}

func NewDatabaseService(db *storage.PostgresDbEngine) *DatabaseService {
	srv := new(DatabaseService)
	srv.db = db
	return srv
}

func (databaseSrv *DatabaseService) GetStatus() *models.Status {
	status := new(models.Status)
	row := databaseSrv.db.CP.QueryRow(context.Background(), `SELECT * FROM
		(SELECT COUNT(1) FROM users) as user_count,
 		(SELECT COUNT(1) FROM forums) as forum_count,
		(SELECT COUNT(1) FROM threads) as thread_count,
		(SELECT COUNT(1) FROM posts) as post_count;`)
	row.Scan(&status.User, &status.Forum, &status.Thread, &status.Post)

	return status
}

func (databaseSrv *DatabaseService) ClearDatabase() {
	databaseSrv.db.CP.Query(context.Background(),
		`TRUNCATE users, forums, threads, votes, posts, forums_users_nicknames`)
}
