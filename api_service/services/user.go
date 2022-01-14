package services

import (
	"context"
	"db_tp/db"
	"db_tp/models"
	"github.com/jackc/pgx/v4"
	"strings"
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

func (userSrv *UserService) GetByNicknameOrEmail(nickname string, email string) models.Users {
	users := make(models.Users, 0)
	userByNickname := userSrv.GetByNickname(nickname)
	if userByNickname != nil {
		users = append(users, *userByNickname)
	}
	userByEmail := userSrv.GetByEmail(email)
	if userByEmail != nil {
		if userByNickname == nil || userByNickname != nil && strings.Compare(userByNickname.Email, userByEmail.Email) != 0 {
			users = append(users, *userByEmail)
		}
	}
	return users
}

func (userSrv *UserService) CreateUser(item *models.User) (models.Users, bool) {
	_, err := userSrv.db.CP.Exec(context.Background(),
		`INSERT INTO users (nickname, email, fullname, about)
				VALUES ($1, $2, $3, $4)`,
		item.Nickname,
		item.Email,
		item.Fullname,
		item.About,
	)
	if err != nil {
		return userSrv.GetByNicknameOrEmail(item.Nickname, item.Email), true
	}

	return []models.User{*item}, false
}

func (userSrv *UserService) GetByEmail(email string) *models.User {
	rows, err := userSrv.db.CP.Query(
		context.Background(),
		`SELECT 	nickname,
                        email,
                        fullname,
                        about
                FROM users
                WHERE email = $1`,
		email)
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

func (userSrv *UserService) UpdateByNickname(nickname string, item *models.UserUpdate) (*models.User, bool) {
	userEmail := new(string)
	if item.Email != "" {
		userEmail = &item.Email
	}
	userFullname := new(string)
	if item.Fullname != "" {
		userFullname = &item.Fullname
	}
	userAbout := new(string)
	if item.About != "" {
		userAbout = &item.About
	}

	user := &models.User{}
	row := userSrv.db.CP.QueryRow(context.Background(),
		`UPDATE users
                    SET email=COALESCE(NULLIF($2, ''), email),
                        fullname=COALESCE(NULLIF($3, ''), fullname),
                        about=COALESCE(NULLIF($4, ''), about)
                    WHERE nickname = $1
                    RETURNING nickname, email, fullname, about`,
		nickname, userEmail, userFullname, userAbout)
	err = row.Scan(
		&user.Nickname,
		&user.Email,
		&user.Fullname,
		&user.About,
	)
	if err == pgx.ErrNoRows {
		return nil, false
	}
	if err != nil {
		return nil, true
	}
	return user, false
}
