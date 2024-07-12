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

func GerAndUpdateOrderStatusByAccrual(login string, number string) {
	url := config.AccrualSystemAddress + "/api/orders/" + number
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	body, err := client.Do(r)
	if err != nil {
		helpers.TLog.Error(err.Error())
	}
	defer body.Body.Close()
	if body.StatusCode == http.StatusOK {
		var tModel model.OrdersAccrualModel
		err = json.NewDecoder(body.Body).Decode(&tModel)
		if err != nil {
			helpers.TLog.Error(err.Error())
			return
		}
		if *tModel.Status == "PROCESSING" || *tModel.Status == "INVALID" || *tModel.Status == "PROCESSED" {
			err := database.DBStorage.UpdateOrders(tModel)
			if err != nil {
				helpers.TLog.Error(err.Error())
				return
			}
			if *tModel.Status == "PROCESSED" {
				helpers.TLog.Info("LOGIN::: ", login)
				helpers.TLog.Info("Balance::: ", *tModel.Accrual)
				err := database.DBStorage.CreateOrUpdateCurrentBalance(model.CurrentBalanceModel{Login: &login, Balance: tModel.Accrual})
				if err != nil {
					helpers.TLog.Error(err.Error())
					return
				}
			}
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
			GerAndUpdateOrderStatusByAccrual(*order.Login, *order.OrdersID)
		}
		helpers.TLog.Info("Окончание проверки статусов")
		time.Sleep(3 * time.Second)

	}

}
