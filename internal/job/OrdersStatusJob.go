package job

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/andrey67895/go_diplom_first/internal/config"
	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

var client http.Client

func GetAndUpdateOrderStatusByAccrual(login string, number string) (*http.Response, error) {
	url := config.AccrualSystemAddress + "/api/orders/" + number
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	body, err := client.Do(r)
	if err != nil {
		helpers.TLog.Error(err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			helpers.TLog.Error(err.Error())
		}
	}(body.Body)
	if body.StatusCode == http.StatusOK {
		var tModel model.OrdersAccrualModel
		err = json.NewDecoder(body.Body).Decode(&tModel)
		if err != nil {
			helpers.TLog.Error(err.Error())
			return body, err
		}
		if *tModel.Status == "PROCESSING" || *tModel.Status == "INVALID" || *tModel.Status == "PROCESSED" {
			err := database.DBStorage.UpdateOrders(tModel, *tModel.Status == "PROCESSED", login)
			if err != nil {
				helpers.TLog.Error(err.Error())
				return body, err
			}
		}
	}
	return body, nil
}

func OrdersStatusJob(ctx context.Context, wg *sync.WaitGroup) {
	second := 1
	ticker := time.NewTicker(time.Duration(second) * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				helpers.TLog.Info("Job: Завершение работы")
				return
			case <-ticker.C:
				helpers.TLog.Info("Job: Запуск проверки статусов")
				orders, err := database.DBStorage.GetOrdersByNotFinalStatus()
				if err != nil {
					helpers.TLog.Error(err.Error())
					return
				}
				for _, order := range orders {
					body, _ := GetAndUpdateOrderStatusByAccrual(*order.Login, *order.OrdersID)
					if body.StatusCode == 429 {
						i, err := strconv.Atoi(body.Header.Get("Retry-After"))
						if err != nil {
							helpers.TLog.Error(err.Error())
						}
						ticker.Reset(time.Duration(i) * time.Second)
					} else {
						ticker.Reset(time.Duration(second) * time.Second)
					}
					err := body.Body.Close()
					if err != nil {
						helpers.TLog.Error(err.Error())
					}
				}
				helpers.TLog.Info("Job: Окончание проверки статусов")

			}
		}
	}()

	//for range ticker.C {
	//	helpers.TLog.Info("Job: Запуск проверки статусов")
	//	orders, err := database.DBStorage.GetOrdersByNotFinalStatus()
	//	if err != nil {
	//		helpers.TLog.Error(err.Error())
	//		return
	//	}
	//	for _, order := range orders {
	//		body, _ := GetAndUpdateOrderStatusByAccrual(*order.Login, *order.OrdersID)
	//		if body.StatusCode == 429 {
	//			i, err := strconv.Atoi(body.Header.Get("Retry-After"))
	//			if err != nil {
	//				helpers.TLog.Error(err.Error())
	//			}
	//			ticker.Reset(time.Duration(i) * time.Second)
	//		} else {
	//			ticker.Reset(time.Duration(second) * time.Second)
	//		}
	//		body.Body.Close()
	//	}
	//	helpers.TLog.Info("Job: Окончание проверки статусов")
	//}
}
