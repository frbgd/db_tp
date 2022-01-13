package services

import (
	"db_tp/db"
	"db_tp/models"
)

var ForumSrv *ForumService

type ForumService struct {
	db *db.PostgresDbEngine
}

func NewForumService(db *db.PostgresDbEngine) *ForumService {
	srv := new(ForumService)
	srv.db = db
	return srv
}

func (forum *ForumService) GetBySlug(slug string) *models.Forum {

}

// TODO переместить в другой сервис
func (forum *ForumService) GetThreadBySlug(slug string) (*models.Thread, error) {

}

func (forum *ForumService) GetForumUsers(slug string, desc string, limit string, since string) *models.User[] {

}

func (forum *ForumService) GetForumThreads(slug string, desc string, limit string, since string) *models.Thread[] {

}

func (forum *ForumService) CreateForum(item *models.Forum) (*models.Forum, error) {

}

func (forum *ForumService) CreateThread(item *models.Thread) (*models.Thread, error) {

}
