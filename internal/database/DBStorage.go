package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/andrey67895/go_diplom_first/internal/config"
	"github.com/andrey67895/go_diplom_first/internal/database/migrator"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DBStorage DBStorageModel

type DBStorageModel struct {
	DB  *sql.DB
	ctx context.Context
}

//go:embed migrations/*.sql
var MigrationsFS embed.FS

const migrationsDir = "migrations"

func openDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", config.DatabaseDsn)
	if err != nil {
		helpers.TLog.Error(err.Error())
	}
	return db, err
}

func InitDB(ctx context.Context) error {
	db, err := openDB()
	tCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := db.PingContext(tCtx); err != nil {
		helpers.TLog.Error(err.Error())
		return err
	}
	dbStorage := DBStorageModel{DB: db, ctx: ctx}

	tMigrator := migrator.NewMigrator(MigrationsFS, migrationsDir)
	helpers.TLog.Info("Запуск миграции DB")
	db, err = openDB()
	if err != nil {
		return err
	}
	err = tMigrator.ApplyMigrations(db)
	if err != nil {
		return err
	}
	helpers.TLog.Info("Завершение миграции DB")
	DBStorage = dbStorage
	return nil
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

func (db DBStorageModel) GetOrdersByNotFinalStatus() ([]*model.OrdersModel, error) {
	data := make([]*model.OrdersModel, 0)

	rows, err := db.DB.QueryContext(db.ctx, "SELECT * from orders where status in ('NEW', 'PROCESSING') order by uploaded_at")
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
		data = append(data, &v)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (db DBStorageModel) GetOrdersByLogin(login string) ([]*model.OrdersModel, error) {
	data := make([]*model.OrdersModel, 0, 10)

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
		data = append(data, &v)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return data, nil
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
		_, err = tx.ExecContext(db.ctx, `UPDATE orders SET accrual=$1,status=$2 WHERE orders_id=$3`, ordersAccrualModel.Accrual, ordersAccrualModel.Status, ordersAccrualModel.OrderID)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.ExecContext(db.ctx, `INSERT INTO current_balance as ca (login, current) values ($1,$2) on conflict (login) do update set current = (EXCLUDED.current  + ca."current")`, login, ordersAccrualModel.Accrual)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		return err
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

func (db DBStorageModel) GetWithdrawnBalanceByLogin(login string) ([]*model.WithdrawnBalanceModel, error) {
	data := make([]*model.WithdrawnBalanceModel, 0)

	rows, err := db.DB.QueryContext(db.ctx, "SELECT * from withdrawn_balance where login = $1", login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var v model.WithdrawnBalanceModel
		err = rows.Scan(&v.Login, &v.Order, &v.ProcessedAT, &v.Withdrawn)
		if err != nil {
			return nil, err
		}
		data = append(data, &v)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return data, nil
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

func (db DBStorageModel) WithdrawnBalanceByLogin(withdrawnBalanceModel model.WithdrawnBalanceModel) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	_, err = tx.ExecContext(db.ctx, `INSERT INTO withdrawn_balance as wb (login,"order",withdrawn) values ($1,$2,$3)`, withdrawnBalanceModel.Login, withdrawnBalanceModel.Order, withdrawnBalanceModel.Withdrawn)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(db.ctx, `UPDATE current_balance as cb SET current = (cb.current-$1) WHERE login=$2`, withdrawnBalanceModel.Withdrawn, withdrawnBalanceModel.Login)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (db DBStorageModel) CreateOrUpdateCurrentBalance(currentBalanceModel model.CurrentBalanceModel) error {
	_, err := db.DB.ExecContext(db.ctx, `INSERT INTO current_balance as ca (login, current) values ($1,$2) on conflict (login) do update set current = (EXCLUDED.current  + ca."current")`, *currentBalanceModel.Login, *currentBalanceModel.Balance)
	return err
}
