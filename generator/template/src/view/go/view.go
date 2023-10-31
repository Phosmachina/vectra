package _go

import (
	. "Vectra/src/model/service"
	. "Vectra/src/model/storage"
	. "github.com/phosmachina/FluentKV/reldb"
)

type GlobalCtx struct {
	IsDev    bool
	Domain   string
	TabTitle string
	Lang     string

	User UserCtx
}

func NewGlobalCtx(tabSuffix string, userId string) GlobalCtx {

	config := GetApiV1().GetStore().Config

	ctx := GlobalCtx{
		IsDev:    config.IsDev,
		TabTitle: config.TabPrefix + tabSuffix,
		Lang:     config.DefaultLang,
		User:     newUserCtx(userId),
	}

	if ctx.IsDev {
		ctx.Domain = config.DevDomain
	} else {
		ctx.Domain = config.Domain
	}

	return ctx
}

type UserCtx struct {
	ID          string
	Role        Role
	IsActivated bool
	Firstname   string
	Lastname    string
	Email       string
}

func newUserCtx(userId string) UserCtx {

	if userId == "" {
		return UserCtx{}
	}

	db := *GetApiV1().GetStore().DB
	userWrp := Get[User](db, userId)
	user := userWrp.Value
	role := AllFromLink[User, Role](db, userId)[0].Value

	return UserCtx{
		ID:          userId,
		Role:        role,
		IsActivated: user.IsActivated,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Email:       user.Email,
	}
}
