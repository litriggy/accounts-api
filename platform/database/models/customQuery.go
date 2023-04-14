package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

type ServiceP struct {
	ID           int32          `json:"serviceId"`
	Name         string         `json:"name"`
	Symbol       string         `json:"symbol"`
	Decimals     int32          `json:"decimals"`
	Image        sql.NullString `json:"img"`
	IsNative     int32          `json:"isNative"`
	ContractAddr sql.NullString `json:"contractAddr"`
	NetType      sql.NullString `json:"netType"`
}

type GetUserServicesRowP struct {
	ServiceID    int32          `json:"serviceId"`
	Amount       int64          `json:"amount"`
	Name         string         `json:"name"`
	Symbol       string         `json:"symbol"`
	Decimals     int32          `json:"decimals"`
	Image        sql.NullString `json:"img"`
	IsNative     int32          `json:"isNative"`
	ContractAddr sql.NullString `json:"contractAddr"`
	NetType      sql.NullString `json:"netType"`
}

func (q *Queries) GetOptBalances(ctx context.Context, userID int32, selected string) ([]GetBalancesRow, error) {
	rows, err := q.db.QueryContext(ctx, fmt.Sprintln("SELECT bal.amount, s.name, s.symbol, s.contract_addr, s.is_native, s.image, s.decimals FROM user_balances AS bal LEFT JOIN services AS s ON bal.service_id = s.id WHERE bal.user_id = '"+strconv.Itoa(int(userID))+"' AND bal.service_id IN "+selected))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBalancesRow
	for rows.Next() {
		var i GetBalancesRow
		if err := rows.Scan(
			&i.Amount,
			&i.Name,
			&i.Symbol,
			&i.ContractAddr,
			&i.IsNative,
			&i.Image,
			&i.Decimals,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type GetTransactionsRow struct {
	ID          int32
	FromID      int32
	ToID        int32
	Memo        sql.NullString
	TotalAmount int64
	Name        sql.NullString
	Symbol      sql.NullString
	Decimals    sql.NullInt32
	Image       sql.NullString
	IsNative    sql.NullInt32
	NetType     sql.NullString
	WalletType  sql.NullString
}

func (q *Queries) GetTransactions(ctx context.Context, userID string, lim string, off string) ([]GetTransactionsRow, error) {
	stmt := fmt.Sprintln(
		"SELECT t.id, t.from_id, t.to_id, t.memo, t.total_amount, s.name, s.symbol, s.decimals, s.image, s.is_native, s.net_type, s.wallet_type " +
			"FROM transactions AS t " +
			"LEFT JOIN services AS s " +
			"ON s.id = t.service_id " +
			"WHERE from_id = '" +
			userID +
			"' OR to_id IN ( " +
			"SELECT wallet_addr FROM user_wallets WHERE user_id = '" + userID + "' ) " +
			"OR to_id = '" +
			userID + "' " +
			"LIMIT " + lim +
			" OFFSET " + off,
	)
	rows, err := q.db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTransactionsRow
	for rows.Next() {
		var i GetTransactionsRow
		if err := rows.Scan(
			&i.ID,
			&i.FromID,
			&i.ToID,
			&i.Memo,
			&i.TotalAmount,
			&i.Name,
			&i.Symbol,
			&i.Decimals,
			&i.Image,
			&i.IsNative,
			&i.NetType,
			&i.WalletType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *Queries) RawSQLExec(ctx context.Context, stmt string) error {
	_, err := q.db.ExecContext(ctx, stmt)
	return err
}
