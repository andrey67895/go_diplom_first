package services

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

func GetOrdersAndSortByLogin(login string) ([]*model.OrdersModel, error) {
	orders, err := database.DBStorage.GetOrdersByLogin(login)
	if err != nil {
		helpers.TLog.Error(err.Error())
		return nil, err
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAT.After(*orders[j].UploadedAT)
	})
	return orders, nil
}

func GetOrderIDAndValid(body io.ReadCloser) (*string, error) {
	b, err := io.ReadAll(body)
	if err != nil {
		err := fmt.Errorf("неверный формат номера заказа! %s", err.Error())
		helpers.TLog.Error(err.Error())
		return nil, err
	}
	orderID, err := strconv.Atoi(string(b))
	if !helpers.LuhnValid(orderID) || err != nil {
		err := fmt.Errorf("неверный формат номера заказа")
		return nil, err
	}
	tOrderID := string(b)
	return &tOrderID, nil
}

func CreateOrders(tModel model.OrdersModel) error {
	err := database.DBStorage.CreateOrders(tModel)
	if err != nil {
		helpers.TLog.Error(err.Error())
		return err
	}
	return nil
}

func GetOrderByOrderIDOrCreate(tModel model.OrdersModel) (*model.OrdersModel, error) {
	orders, err := database.DBStorage.GetOrdersByOrderID(*tModel.OrdersID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := CreateOrders(tModel); err != nil {
				return nil, err
			}
			return nil, nil
		} else {
			helpers.TLog.Error(err.Error())
			return nil, err
		}
	}
	return orders, nil
}
