package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/model"
)

func GetAuth(tModel model.UserModel, create bool) (*model.UserModel, *model.ApiError) {
	auth, err := database.DBStorage.GetAuth(*tModel.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if create {
				err := database.DBStorage.CreateAuth(tModel)
				if err != nil {
					helpers.TLog.Error(err.Error())
					return nil, &model.ApiError{
						Status: http.StatusInternalServerError,
						Error:  err,
					}
				}
			} else {
				err := fmt.Errorf("неверная пара логин/пароль")
				helpers.TLog.Error(err.Error())
				return nil, &model.ApiError{
					Status: http.StatusUnauthorized,
					Error:  err,
				}
			}
		} else {
			return nil, &model.ApiError{
				Status: http.StatusInternalServerError,
				Error:  err,
			}
		}
	}
	if auth != nil && create {
		return nil, &model.ApiError{Status: http.StatusConflict,
			Error: fmt.Errorf("пользователь уже существует"),
		}
	} else if auth == nil && !create {
		return nil, &model.ApiError{
			Status: http.StatusUnauthorized,
			Error:  fmt.Errorf("неверная пара логин/пароль"),
		}
	}
	return auth, nil
}
