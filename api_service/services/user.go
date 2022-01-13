package services

import (
	"db_tp/db"
	"db_tp/models"
)

var UserSrv *UserService

type UserService struct {
	db *db.PostgresDbEngine
}

func NewUserService(db *db.PostgresDbEngine) *UserService {
	srv := new(UserService)
	srv.db = db
	return srv
}

func (forum *UserService) GetUsersByNicknameOrEmail(nickname string, email string) (*models.User[], error) {

}

func (forum *UserService) CreateUser(item *models.User) (*models.User, error) {

}

func (forum *UserService) GetByEmail(email string) (*models.User, error) {

}

func (forum *UserService) GetByNickname(nickname string) (*models.User, error) {

}

func (forum *UserService) UpdateByNickname(nickname string, item *models.UserUpdate) (*models.User, error) {

}
