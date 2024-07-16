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

func GetOrdersAndSortByLogin(login string, w http.ResponseWriter) *[]model.OrdersModel {
	orders, err := database.DBStorage.GetOrdersByLogin(login)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, "Ошибка сервера!", http.StatusInternalServerError)
	}
	tOrders := *orders
	sort.Slice(tOrders, func(i, j int) bool {
		return tOrders[i].UploadedAT.After(*tOrders[j].UploadedAT)
	})
	return orders
}

func GetOrderIDAndValid(w http.ResponseWriter, req *http.Request) string {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, "Неверный формат номера заказа!", http.StatusUnprocessableEntity)
	}
	orderID, err := strconv.Atoi(string(b))
	if !helpers.LuhnValid(orderID) || err != nil {
		http.Error(w, "Неверный формат номера заказа!", http.StatusUnprocessableEntity)
	}

	return string(b)
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
