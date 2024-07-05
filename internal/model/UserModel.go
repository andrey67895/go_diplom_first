package model

import "fmt"

type UserModel struct {
	ID       *int64  `json:"id,omitempty"`
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

func (u UserModel) IsValid() error {
	if u.Login == nil || u.Password == nil || *u.Login == "" || *u.Password == "" {
		return fmt.Errorf("ошибка валидации! Обязательные поля: password и login, не могут быть пустыми или null: %+v", u)
	}
	return nil
}
