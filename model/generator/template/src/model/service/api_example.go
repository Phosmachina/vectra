package service

import (
	. "Vectra/src/model/storage"
	. "github.com/phosmachina/FluentKV/reldb"
	"golang.org/x/crypto/bcrypt"
)

func (api *ApiV1) CreateUser(info UserExch) error {

	db := *api.store.DB
	if FindFirst(db, func(id string, value *User) bool {
		return value.Email == info.Email
	}) != nil {
		return ErrorUserExist
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	userWrp := Insert(db, User{ // BUG: use constructor to init IObject
		Password:  password,
		Firstname: info.Firstname,
		Lastname:  info.Lastname,
		Email:     info.Email,
	})
	Link(userWrp, true, api.accessManager.DefaultRoles["registered"])

	return nil
}
