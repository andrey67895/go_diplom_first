package model

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/andrey67895/go_diplom_first/internal/helpers"
)

type UserModel struct {
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

func UserModelDecode(body io.ReadCloser) (*UserModel, error) {
	var tModel UserModel
	err := json.NewDecoder(body).Decode(&tModel)
	if err != nil {
		helpers.TLog.Error(err.Error())
		return nil, fmt.Errorf("ошибка десериализации")
	}
	return &tModel, nil
}

func (u UserModel) IsValid() error {
	if u.Login == nil || u.Password == nil || *u.Login == "" || *u.Password == "" {
		err := fmt.Errorf("ошибка валидации! Обязательные поля: password и login, не могут быть пустыми или null: %+v", u)
		helpers.TLog.Error(err.Error())
		return err
	}
	return nil
}
