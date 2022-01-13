package services

import (
	"db_tp/db"
	"db_tp/models"
)

var UserSrv UserService

type UserService struct {
	db *db.PostgresDbEngine
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
