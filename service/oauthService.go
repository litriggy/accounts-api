package service

import (
	"context"
	"database/sql"
	"strconv"

	"accounts/api/pkg/config"
	db "accounts/api/platform/database/models"
	"accounts/api/platform/logger"
)

var logr = logger.GetLogger()

func FindUser(oauthID string, oauthType string) (db.User, error) {
	ctx := context.Background()
	query, err := config.ConnectDB()
	if err != nil {
		return db.User{}, err
	}
	result, err := query.GetUserByOauth(ctx, db.GetUserByOauthParams{OauthType: oauthType, OauthID: oauthID})
	if err != nil {
		return db.User{}, err
	}
	return result, nil
}

func CreateUser(oauthID string, oauthType string, version int32, nickname string, email string, userType string, picture string) (string, error) {
	ctx := context.Background()
	query, database, err := config.DBConn()
	if err != nil {
		return "", err
	}

	tx, err := database.Begin()

	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	q := query.WithTx(tx)
	user, err := q.CreateUser(ctx, db.CreateUserParams{
		Nickname: sql.NullString{
			String: nickname,
			Valid:  true,
		},
		Email: email,
		Type:  userType,
		Picture: sql.NullString{
			String: picture,
			Valid:  true,
		},
	})
	if err != nil {
		return "", err
	}
	userid, err := user.LastInsertId()
	if err != nil {
		return "", err
	}

	if _, err := q.CreateOauth(ctx, db.CreateOauthParams{UserID: int32(userid), OauthID: oauthID, OauthType: oauthType, Version: version}); err != nil {
		return "", err
	}

	return strconv.Itoa(int(userid)), tx.Commit()
}
