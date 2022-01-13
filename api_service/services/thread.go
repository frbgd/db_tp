package services

import (
	"db_tp/db"
	"db_tp/models"
)

var ThreadSrv *ThreadService

type ThreadService struct {
	db *db.PostgresDbEngine
}

func NewThreadService(db *db.PostgresDbEngine) *ThreadService {
	srv := new(ThreadService)
	srv.db = db
	return srv
}

func (forum *ThreadService) GetBySlugOrId(slugOrId string) *models.Thread {

}

func (forum *ThreadService) GetBySlug(slug string) *models.Thread {

}

func (forum *ThreadService) GetById(id int) *models.Thread {

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

func (forum *ThreadService) GetPosts(threadId int64, desc string, limit string, since string, sort string) *models.Post[] {

}
