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

func (q *Queries) RawSQLExec(ctx context.Context, stmt string) error {
	_, err := q.db.ExecContext(ctx, stmt)
	return err
}
