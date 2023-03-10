package repository

import (
	"game_assistantor/model"
)

var GameRoleRepos GameRoleRepository

type GameRoleRepository struct {
}

func (u *GameRoleRepository) GetAccountInfo(accountId string) (account *model.GameAccount, err error) {
	account = new(model.GameAccount)
	account.AccountId = accountId
	engine.First(account)
	return
}

func (u *GameRoleRepository) UpdateAccountInfo(account model.GameAccount) (err error) {
	err = engine.Save(&account).Error
	return
}

func (u *GameRoleRepository) GetAccountRoleList(accountId string) (roleList []model.GameRole, err error) {
	err = engine.Model(model.GameRole{}).Where("game_role.account = ?", accountId).Find(roleList).Error
	return
}

func (u *GameRoleRepository) SaveAccount(account model.GameAccount) (err error) {
	err = engine.Save(&account).Error
	return
}