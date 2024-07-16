package services

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

func GetCurrentBalanceByLogin(login string, w http.ResponseWriter) *model.CurrentBalanceModel {
	currentBalanceModel, err := database.DBStorage.GetCurrentBalanceByLogin(login)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tFloat := 0.0
			currentBalanceModel = &model.CurrentBalanceModel{Balance: &tFloat, Login: &login}
		} else {
			helpers.TLog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	return currentBalanceModel
}
