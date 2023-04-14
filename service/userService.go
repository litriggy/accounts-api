package service

import (
	"accounts/api/pkg/config"
	db "accounts/api/platform/database/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func GetUser(userID int32) (*db.User, error) {
	query, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}

	userInfo, err := query.GetUser(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func CreateSoftWallet(walletAddr string, userId int32) error {
	query, err := config.ConnectDB()
	if err != nil {
		return err
	}

	err = query.CreateUserWallet(context.Background(), db.CreateUserWalletParams{
		UserID:       userId,
		IsIntegrated: 0,
		WalletAddr:   walletAddr,
		SecPk:        sql.NullString{Valid: false},
	})
	if err != nil {
		return err
	}
	return nil
}

func CreateHardWallet(walletAddr string, secPk string, userId int32) error {
	query, err := config.ConnectDB()
	if err != nil {
		return err
	}

	err = query.CreateUserWallet(context.Background(), db.CreateUserWalletParams{
		UserID:       userId,
		IsIntegrated: 1,
		WalletAddr:   walletAddr,
		SecPk:        sql.NullString{String: secPk, Valid: true},
		WalletType:   sql.NullString{String: "eth", Valid: true},
	})
	if err != nil {
		return err
	}
	return nil
}

func GetBalances(userID int32, selected string) ([]db.GetBalancesRow, error) {
	query, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	rows, err := query.GetOptBalances(context.Background(), userID, selected)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func getPK(walletAddr string) (string, error) {
	query, err := config.ConnectDB()
	if err != nil {
		return "", err
	}
	result, err := query.GetPKFromWallet(context.Background(), walletAddr)
	if err != nil {
		return "", err
	}
	if !result.Valid {
		return "", errors.New("no pk found")
	}
	return result.String, nil
}

func CreatedSecondPassword(userID int32, secPw string) error {
	query, err := config.ConnectDB()
	if err != nil {
		return err
	}
	if err := query.CreateSecondPassword(
		context.Background(),
		db.CreateSecondPasswordParams{
			UserID: userID,
			SecPw:  sql.NullString{String: secPw, Valid: true},
		},
	); err != nil {
		return err
	}
	return nil
}

func CheckPassword(userID int32, pass string) (bool, error) {
	query, err := config.ConnectDB()
	if err != nil {
		return false, err
	}
	count, err := query.CheckUserPass(
		context.Background(),
		db.CheckUserPassParams{
			UserID: userID,
			SecPw: sql.NullString{
				Valid:  true,
				String: pass,
			},
		})
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckPasswordExists(userID int32) (bool, error) {
	query, err := config.ConnectDB()
	if err != nil {
		return false, err
	}
	count, err := query.CheckPassExists(context.Background(), userID)

	if err != nil {
		return false, nil

	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func FindWallet(walletAddr string) (bool, error) {
	query, err := config.ConnectDB()
	if err != nil {
		return false, err
	}
	result, err := query.FindWallet(context.Background(), walletAddr)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return true, nil
		}
		fmt.Println(err.Error())
		return false, err
	}
	if result != (db.UserWallet{}) {
		return false, errors.New("exists")
	} else {
		return true, nil
	}
}
