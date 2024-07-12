package services

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/andrey67895/go_diplom_first/internal/config"
	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

var client http.Client

func GerAndUpdateOrderStatusByAccrual(number string) {
	url := config.AccrualSystemAddress + "/api/orders/" + number
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	body, err := client.Do(r)
	if err != nil {
		helpers.TLog.Error(err.Error())
	}
	if body.StatusCode == http.StatusOK {
		var tModel model.OrdersAccrualModel
		err = json.NewDecoder(body.Body).Decode(&tModel)
		if err != nil {
			helpers.TLog.Error(err.Error())
			return
		}
		if *tModel.Status == "PROCESSING" {
			database.DBStorage.UpdateOrders(tModel)
		} else if *tModel.Status == "INVALID" {
			database.DBStorage.UpdateOrders(tModel)
		} else if *tModel.Status == "PROCESSED" {
			database.DBStorage.UpdateOrders(tModel)
		}
	}
}

func OrdersStatusJob() {
	for {
		helpers.TLog.Info("Запуск проверки статусов")
		orders, err := database.DBStorage.GetOrdersByNotFinalStatus()
		if err != nil {
			helpers.TLog.Error(err.Error())
			return
		}
		for _, order := range *orders {
			GerAndUpdateOrderStatusByAccrual(*order.OrdersID)
		} //status.
		helpers.TLog.Info("Окончание проверки статусов")
		time.Sleep(20 * time.Second)

	}

}
