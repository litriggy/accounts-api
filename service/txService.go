package service

import (
	"accounts/api/app/model"
	"accounts/api/pkg/config"
	db "accounts/api/platform/database/models"
	"accounts/api/utils/chain"
	"context"
	"database/sql"
	"errors"
	"strconv"
)

type TransactionDetailModel struct {
	From      string
	To        string
	ServiceID int32
	Amount    int64
	Txhash    string
}

func TransferBalance(fromID int32, toID int32, serviceID int32, onChainLog []model.OnTransactionDetail, offChainLog model.OffTransactionDetail, secPw string) error {
	totalAmount := int64(0)
	//strFromID := strconv.Itoa(int(fromID))
	ctx := context.Background()
	query, database, err := config.DBConn()
	if err != nil {
		return err
	}

	tx, err := database.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	q := query.WithTx(tx)

	txData, err := q.CreateTx(ctx, db.CreateTxParams{
		FromID:      fromID,
		ToID:        toID,
		ServiceID:   serviceID,
		TotalAmount: 0,
	})
	if err != nil {
		return err
	}
	txID, err := txData.LastInsertId()
	if err != nil {
		return err
	}
	if (model.OffTransactionDetail{}) != offChainLog {
		result, err := q.GetBalance(ctx, db.GetBalanceParams{
			UserID:    fromID,
			ServiceID: serviceID,
		})
		if err != nil {
			return err
		}

		if result.Amount < offChainLog.Amount {
			return errors.New("not enough balance")
		}
		if err := q.UpdateBalance(ctx, db.UpdateBalanceParams{
			Amount:    -(offChainLog.Amount),
			UserID:    fromID,
			ServiceID: serviceID,
		}); err != nil {
			return err
		}

		if err := q.UpdateBalance(ctx, db.UpdateBalanceParams{
			Amount:    offChainLog.Amount,
			UserID:    toID,
			ServiceID: serviceID,
		}); err != nil {
			return err
		}
		totalAmount += offChainLog.Amount
		q.CreateTxDetails(ctx, db.CreateTxDetailsParams{
			TransactionID: int32(txID),
			From:          strconv.Itoa(int(fromID)),
			To:            offChainLog.To,
			Amount:        offChainLog.Amount,
			IsOnchain:     int32(1),
			Txhash: sql.NullString{
				String: "",
				Valid:  false,
			},
		})

	}
	if len(onChainLog) > 0 {
		serviceInfo, err := GetService(serviceID)
		isNative := true
		if serviceInfo.IsNative != 1 {
			isNative = false
		}
		if err != nil {
			return err
		}
		for _, log := range onChainLog {
			encPrivateKey, err := getPK(log.From)
			if err != nil {
				return err
			}
			privateKey, err := chain.Decrypt(encPrivateKey, secPw)
			if err != nil {
				return err
			}
			txHash, err := chain.Transfer(privateKey, log.To, log.Amount, serviceID, serviceInfo.NetType.String, serviceInfo.ContractAddr.String, isNative)
			if err != nil {
				return err
			}
			q.CreateTxDetails(ctx, db.CreateTxDetailsParams{
				TransactionID: int32(txID),
				From:          log.From,
				To:            log.To,
				Amount:        log.Amount,
				IsOnchain:     int32(1),
				Txhash: sql.NullString{
					String: txHash,
					Valid:  true,
				},
			})
			totalAmount += log.Amount
		}
	}
	// if err := q.UpdateTxHash(ctx, db.UpdateTxHashParams{
	// 	Txhash: sql.NullString{String:},
	// })
	if err := q.UpdateTxTotalAmount(ctx, db.UpdateTxTotalAmountParams{
		TotalAmount: int64(totalAmount),
		ID:          int32(txID),
	}); err != nil {
		return err
	}
	return tx.Commit()
}
