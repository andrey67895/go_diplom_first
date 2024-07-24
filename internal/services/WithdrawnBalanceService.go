package services

import (
	"net/http"
	"sort"

	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

func GetWithdrawnBalanceAndSortByLogin(login string, w http.ResponseWriter) []*model.WithdrawnBalanceModel {
	withdrawnHistory, err := database.DBStorage.GetWithdrawnBalanceByLogin(login)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, "Ошибка сервера!", http.StatusInternalServerError)

	}
	sort.Slice(withdrawnHistory, func(i, j int) bool {
		return withdrawnHistory[i].ProcessedAT.After(*withdrawnHistory[j].ProcessedAT)
	})
	return withdrawnHistory
}

func WithdrawnBalanceByLogin(tModel model.WithdrawnBalanceModel, w http.ResponseWriter) {
	err := database.DBStorage.WithdrawnBalanceByLogin(tModel)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetWithdrawnBalanceSum(login string, w http.ResponseWriter) *float64 {
	withdrawnBalanceSum, err := database.DBStorage.GetWithdrawnBalanceSumByLogin(login)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return withdrawnBalanceSum
}
