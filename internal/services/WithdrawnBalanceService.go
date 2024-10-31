package services

import (
	"fmt"
	"sort"

	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

func GetWithdrawnBalanceAndSortByLogin(login string) ([]*model.WithdrawnBalanceModel, error) {
	withdrawnHistory, err := database.DBStorage.GetWithdrawnBalanceByLogin(login)
	if err != nil {
		helpers.TLog.Error(err.Error())
		return nil, fmt.Errorf("ошибка сервера: %s", err.Error())
	}
	sort.Slice(withdrawnHistory, func(i, j int) bool {
		return withdrawnHistory[i].ProcessedAT.After(*withdrawnHistory[j].ProcessedAT)
	})
	return withdrawnHistory, nil
}

func WithdrawnBalanceByLogin(tModel model.WithdrawnBalanceModel) error {
	err := database.DBStorage.WithdrawnBalanceByLogin(tModel)
	if err != nil {
		helpers.TLog.Error(err.Error())
		return err
	}
	return nil
}

func GetWithdrawnBalanceSum(login string) (*float64, error) {
	withdrawnBalanceSum, err := database.DBStorage.GetWithdrawnBalanceSumByLogin(login)
	if err != nil {
		helpers.TLog.Error(err.Error())
		return nil, err
	}
	return withdrawnBalanceSum, nil
}
