package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andrey67895/go_diplom_first/internal/config"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DBStorage DBStorageModel

type DBStorageModel struct {
	DB  *sql.DB
	ctx context.Context
}

func InitDB(ctx context.Context) {
	db, err := sql.Open("pgx", config.DatabaseDsn)
	if err != nil {
		helpers.TLog.Error(err.Error())
	}
	tCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err = db.PingContext(tCtx); err != nil {
		helpers.TLog.Error(err.Error())
	}
	dbStorage := DBStorageModel{DB: db, ctx: ctx}
	dbStorage.InitTable(tCtx)
	DBStorage = dbStorage
}

func (db DBStorageModel) CreateAuth(authModel model.UserModel) error {
	_, err := db.DB.ExecContext(db.ctx, `INSERT INTO auth(login, hash_pass) values ($1,$2)`, authModel.Login, helpers.EncodeHash(*authModel.Password))
	return err
}

func (db DBStorageModel) GetOrdersByOrderID(orderID string) (*model.OrdersModel, error) {
	var data model.OrdersModel
	rows := db.DB.QueryRowContext(db.ctx, "SELECT * from orders where orders_id = $1", orderID)
	err := rows.Scan(&data.OrdersID, &data.Login, &data.Accrual, &data.Status, &data.UploadedAT)
	if errors.Join(rows.Err(), err) != nil {
		return nil, err
	}
	return &data, nil
}

func (db DBStorageModel) GetOrdersByNotFinalStatus() (*[]model.OrdersModel, error) {
	data := make([]model.OrdersModel, 0)

	rows, err := db.DB.QueryContext(db.ctx, "SELECT * from orders where status in ('NEW', 'PROCESSING')")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var v model.OrdersModel
		err = rows.Scan(&v.OrdersID, &v.Login, &v.Accrual, &v.Status, &v.UploadedAT)
		if err != nil {
			return nil, err
		}
		data = append(data, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (db DBStorageModel) GetOrdersByLogin(login string) (*[]model.OrdersModel, error) {
	data := make([]model.OrdersModel, 0)

	rows, err := db.DB.QueryContext(db.ctx, "SELECT * from orders where login = $1", login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var v model.OrdersModel
		err = rows.Scan(&v.OrdersID, &v.Login, &v.Accrual, &v.Status, &v.UploadedAT)
		if err != nil {
			return nil, err
		}
		data = append(data, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (db DBStorageModel) CreateOrders(ordersModel model.OrdersModel) error {
	_, err := db.DB.ExecContext(db.ctx, `INSERT INTO orders(orders_id, login) values ($1,$2)`, ordersModel.OrdersID, ordersModel.Login)
	return err
}

func (db DBStorageModel) UpdateOrders(ordersAccrualModel model.OrdersAccrualModel, itProcessed bool, login string) error {
	var err error
	if itProcessed {
		tx, err := db.DB.Begin()
		if err != nil {
			return fmt.Errorf("begin transaction: %w", err)
		}
		tx.ExecContext(db.ctx, `UPDATE orders SET accrual=$1,status=$2 WHERE orders_id=$3`, ordersAccrualModel.Accrual, ordersAccrualModel.Status, ordersAccrualModel.OrderID)
		tx.ExecContext(db.ctx, `INSERT INTO current_balance as ca (login, current) values ($1,$2) on conflict (login) do update set current = (EXCLUDED.current  + ca."current")`, login, ordersAccrualModel.Accrual)
		err = tx.Commit()
		if err != nil {
			return err
		}
	} else {
		_, err = db.DB.ExecContext(db.ctx, `UPDATE orders SET accrual=$1,status=$2 WHERE orders_id=$3`, ordersAccrualModel.Accrual, ordersAccrualModel.Status, ordersAccrualModel.OrderID)
	}
	return err
}

func (db DBStorageModel) GetAuth(login string) (*model.UserModel, error) {
	var data model.UserModel
	rows := db.DB.QueryRowContext(db.ctx, "SELECT * from auth where login = $1", login)
	err := rows.Scan(&data.Login, &data.Password)
	if errors.Join(rows.Err(), err) != nil {
		return nil, err
	}
	return &data, nil
}

func (db DBStorageModel) GetCurrentBalanceByLogin(login string) (*model.CurrentBalanceModel, error) {
	var data model.CurrentBalanceModel
	rows := db.DB.QueryRowContext(db.ctx, "SELECT * from current_balance where login = $1", login)
	err := rows.Scan(&data.Login, &data.Balance)
	if errors.Join(rows.Err(), err) != nil {
		return nil, err
	}
	return &data, nil
}

func (db DBStorageModel) GetWithdrawnBalanceByLogin(login string) (*model.WithdrawnBalanceModel, error) {
	var data model.WithdrawnBalanceModel
	rows := db.DB.QueryRowContext(db.ctx, "SELECT * from withdrawn_balance where login = $1", login)
	err := rows.Scan(&data.Login, &data.Withdrawn)
	if errors.Join(rows.Err(), err) != nil {
		return nil, err
	}
	return &data, nil
}

func (db DBStorageModel) GetWithdrawnBalanceSumByLogin(login string) (*float64, error) {
	var data float64
	rows := db.DB.QueryRowContext(db.ctx, "select coalesce(sum(wb.withdrawn), 0) from withdrawn_balance wb where wb.login = $1", login)
	err := rows.Scan(&data)
	if errors.Join(rows.Err(), err) != nil {
		return nil, err
	}
	return &data, nil
}

func (db DBStorageModel) CreateOrUpdateCurrentBalance(currentBalanceModel model.CurrentBalanceModel) error {
	_, err := db.DB.ExecContext(db.ctx, `INSERT INTO current_balance as ca (login, current) values ($1,$2) on conflict (login) do update set current = (EXCLUDED.current  + ca."current")`, *currentBalanceModel.Login, *currentBalanceModel.Balance)
	return err
}

func (db DBStorageModel) InitTable(ctx context.Context) {
	_, err := db.DB.ExecContext(ctx, `DROP TABLE IF EXISTS auth; DROP TABLE IF EXISTS orders; DROP TABLE IF EXISTS withdrawn_balance; DROP TABLE IF EXISTS current_balance;`)
	if err != nil {
		helpers.TLog.Error(err.Error())
	}
	_, err = db.DB.ExecContext(ctx, `
		CREATE TABLE auth (
        	"login" text not null primary key,
        	"hash_pass" text not null);
		CREATE TABLE orders (
			"orders_id" varchar primary key,
			"login" text not null,
			"accrual"  double precision,
			"status" text not null default 'NEW',
			"uploaded_at" timestamp not null default now());
		CREATE TABLE withdrawn_balance (
			"login" varchar,
			"withdrawn" double precision not null);
		CREATE TABLE current_balance (
			"login" varchar primary key,
			"current" double precision not null);
		CREATE INDEX withdrawn_login_idx ON withdrawn_balance (login);
	`)
	if err != nil {
		helpers.TLog.Error(err.Error())
	}
}
