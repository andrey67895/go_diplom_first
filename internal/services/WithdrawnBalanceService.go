package services

import (
	"net/http"
	"sort"

	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

func GetWithdrawnBalanceAndSortByLogin(login string, w http.ResponseWriter) *[]model.WithdrawnBalanceModel {
	withdrawnHistory, err := database.DBStorage.GetWithdrawnBalanceByLogin(login)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, "Ошибка сервера!", http.StatusInternalServerError)

	}
	tOrders := *withdrawnHistory
	sort.Slice(tOrders, func(i, j int) bool {
		return tOrders[i].ProcessedAT.After(*tOrders[j].ProcessedAT)
	})

	return withdrawnHistory
}