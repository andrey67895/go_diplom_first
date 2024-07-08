package model

import "fmt"

type UserModel struct {
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

func (u UserModel) IsValid() error {
	if u.Login == nil || u.Password == nil || *u.Login == "" || *u.Password == "" {
		return fmt.Errorf("ошибка валидации! Обязательные поля: password и login, не могут быть пустыми или null: %+v", u)
	}
	return nil
}
