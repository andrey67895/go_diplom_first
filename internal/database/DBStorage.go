package database

import (
	"context"
	"database/sql"
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

func (db DBStorageModel) GetOrders(orderId int) (*model.OrdersModel, error) {
	var data model.OrdersModel
	rows := db.DB.QueryRow("SELECT * from orders where orders_id = $1", orderId)
	err := rows.Scan(&data.OrdersID, &data.LoginId)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (db DBStorageModel) CreateOrders(ordersModel model.OrdersModel) error {
	_, err := db.DB.ExecContext(db.ctx, `INSERT INTO orders(orders_id, login) values ($1,$2)`, ordersModel.OrdersID, ordersModel.LoginId)
	return err
}

func (db DBStorageModel) GetAuth(login string) (*model.UserModel, error) {
	var data model.UserModel
	rows := db.DB.QueryRow("SELECT * from auth where login = $1", login)
	err := rows.Scan(&data.Login, &data.Password)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (db DBStorageModel) InitTable(ctx context.Context) {
	_, err := db.DB.ExecContext(ctx, `DROP TABLE IF EXISTS auth; DROP TABLE IF EXISTS orders;`)
	if err != nil {
		helpers.TLog.Error(err.Error())
	}
	_, err = db.DB.ExecContext(ctx, `
		CREATE TABLE auth (
        	"login" text not null primary key,
        	"hash_pass" text not null);
		CREATE TABLE orders (
			"orders_id" bigint primary key,
			"login" text not null);
	`)
	if err != nil {
		helpers.TLog.Error(err.Error())
	}
}
