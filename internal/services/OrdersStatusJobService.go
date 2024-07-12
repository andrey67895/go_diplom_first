package services

import (
	"io"
	"net/http"
	"time"

	"github.com/andrey67895/go_diplom_first/internal/config"
	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
)

var client http.Client

func GerOrderStatusByAccrual(number string) {
	url := config.AccrualSystemAddress + "/api/orders/" + number
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	body, err := client.Do(r)
	if err != nil {
		helpers.TLog.Error(err.Error())
	}
	all, err := io.ReadAll(body.Body)
	if err != nil {
		helpers.TLog.Error(err.Error())
		return
	}
	helpers.TLog.Info("NEN :::: ", string(all))
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
			GerOrderStatusByAccrual(*order.OrdersID)
		} //status.
		helpers.TLog.Info("Окончание проверки статусов")
		time.Sleep(20 * time.Second)

	}

}
