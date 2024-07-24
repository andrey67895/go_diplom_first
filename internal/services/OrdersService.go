package services

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

func GetOrdersAndSortByLogin(login string, w http.ResponseWriter) []*model.OrdersModel {
	orders, err := database.DBStorage.GetOrdersByLogin(login)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, "Ошибка сервера!", http.StatusInternalServerError)
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAT.After(*orders[j].UploadedAT)
	})
	return orders
}

func GetOrderIDAndValid(w http.ResponseWriter, req *http.Request) *string {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, "Неверный формат номера заказа!", http.StatusUnprocessableEntity)
		return nil
	}
	orderID, err := strconv.Atoi(string(b))
	if !helpers.LuhnValid(orderID) || err != nil {
		http.Error(w, "Неверный формат номера заказа!", http.StatusUnprocessableEntity)
		return nil
	}
	tOrderID := string(b)
	return &tOrderID
}

func CreateOrders(tModel model.OrdersModel, w http.ResponseWriter) {
	err := database.DBStorage.CreateOrders(tModel)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, "Ошибка сервера!", http.StatusInternalServerError)
	}
}

func GetOrderByOrderIDOrCreate(tModel model.OrdersModel, w http.ResponseWriter) *model.OrdersModel {
	orders, err := database.DBStorage.GetOrdersByOrderID(*tModel.OrdersID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			CreateOrders(tModel, w)
			w.WriteHeader(http.StatusAccepted)
			return nil
		} else {
			helpers.TLog.Error(err.Error())
			http.Error(w, "Ошибка сервера!", http.StatusInternalServerError)
			return nil
		}
	}
	return orders
}
