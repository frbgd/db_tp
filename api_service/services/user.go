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

func (userSrv *UserService) GetUsersByNicknameOrEmail(nickname string, email string) (*models.User[], error) {

}

func (userSrv *UserService) CreateUser(item *models.User) (*models.User, error) {

}

func (userSrv *UserService) GetByEmail(email string) (*models.User, error) {

}

func (userSrv *UserService) GetByNickname(nickname string) *models.User {

}

func (userSrv *UserService) UpdateByNickname(nickname string, item *models.UserUpdate) (*models.User, error) {

}
