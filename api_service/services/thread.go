package services

import (
	"db_tp/db"
	"db_tp/models"
)

var ThreadeSrv ThreadService

type ThreadService struct {
	db *db.PostgresDbEngine
}

func (forum *ThreadService) GetBySlugOrId(slugOrId string) (*models.Thread, error) {

}

func (forum *ThreadService) GetBySlug(slug string) (*models.Thread, error) {

}

func (forum *ThreadService) GetById(id int64) (*models.Thread, error) {

}

func (forum *ThreadService) CreatePosts(posts *models.Post[]) (*models.Post, error) {

}

func (forum *ThreadService) UpdateBySlugOrId(slugOrId string, item *models.ThreadUpdate) (*models.Thread, error) {

}

func (forum *ThreadService) UpdateBySlug(slug string, item *models.ThreadUpdate) (*models.Thread, error) {

}

func (forum *ThreadService) UpdateById(id int64, item *models.ThreadUpdate) (*models.Thread, error) {

}

func (forum *ThreadService) VoteById(id int64, item *models.Vote) error {

}

func (forum *ThreadService) GetPosts(threadId int64, desc string, limit string, since string, sort string) error {

}
