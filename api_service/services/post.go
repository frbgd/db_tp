package services

import (
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

func (forum *PostService) GetUserByNickname(nickname string) (*models.User, error) {

}

// TODO переместить в другой сервис
func (forum *PostService) GetThreadById(id int64) (*models.Thread, error) {

}

func (forum *PostService) GetForumBySlug(slug string) (*models.Forum, error) {

}

func (forum *PostService) GetById(id int64) (*models.FullPost, error) {

}

func (forum *PostService) UpdateById(id int64) (*models.PostUpdate, error) {

}
