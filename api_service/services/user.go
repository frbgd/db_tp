package services

import (
	"context"
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

//func (userSrv *UserService) GetUsersByNicknameOrEmail(nickname string, email string) (*models.User[], error) {
//
//}
//
//func (userSrv *UserService) CreateUser(item *models.User) (*models.User, error) {
//
//}
//
//func (userSrv *UserService) GetByEmail(email string) (*models.User, error) {
//
//}

func (userSrv *UserService) GetByNickname(nickname string) *models.User {
	rows, err := userSrv.db.CP.Query(
		context.Background(),
		`SELECT 	nickname,
                        email,
                        fullname,
                        about
                FROM users
                WHERE nickname = $1`,
		nickname)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		user := new(models.User)
		err = rows.Scan(
			&user.Nickname,
			&user.Email,
			&user.Fullname,
			&user.About,
		)
		if err != nil {
			panic(err)
		}
		return user
	} else {
		return nil
	}
}

//func (userSrv *UserService) UpdateByNickname(nickname string, item *models.UserUpdate) (*models.User, error) {
//
//}
