package services

import (
	"db_tp/db"
	"db_tp/models"
)

var ForumSrv ForumService

type ForumService struct {
	db *db.PostgresDbEngine
}

func (forum *ForumService) GetBySlug(slug string) (*models.Forum, error) {

}

// TODO переместить в другой сервис
func (forum *ForumService) GetThreadBySlug(slug string) (*models.Thread, error) {

}

func (forum *ForumService) GetForumUsers(slug string, desc string, limit string, since string) (*models.User[], error) {

}

func (forum *ForumService) GetForumThreads(slug string, desc string, limit string, since string) (*models.Thread[], error) {

}

func (forum *ForumService) CreateForum(item *models.Forum) (*models.Forum, error) {

}

func (forum *ForumService) CreateThread(item *models.Thread) (*models.Thread, error) {

}
