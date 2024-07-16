package model

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andrey67895/go_diplom_first/internal/helpers"
)

type UserModel struct {
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

func UserModelDecode(w http.ResponseWriter, r *http.Request) UserModel {
	var tModel UserModel
	err := json.NewDecoder(r.Body).Decode(&tModel)
	if err != nil {
		helpers.TLog.Error(err.Error())
		http.Error(w, "Ошибка десериализации!", http.StatusBadRequest)
	}
	return tModel
}

func (u UserModel) IsValid(w http.ResponseWriter) {
	if u.Login == nil || u.Password == nil || *u.Login == "" || *u.Password == "" {
		helpers.TLog.Error(fmt.Errorf("ошибка валидации! Обязательные поля: password и login, не могут быть пустыми или null: %+v", u).Error())
		http.Error(w, fmt.Errorf("ошибка валидации! Обязательные поля: password и login, не могут быть пустыми или null: %+v", u).Error(), http.StatusBadRequest)
	}
}
