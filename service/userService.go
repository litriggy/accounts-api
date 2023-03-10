package service

import (
	"accounts/api/pkg/config"
	db "accounts/api/platform/database/models"
	"context"
	"database/sql"
	"errors"
)

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
