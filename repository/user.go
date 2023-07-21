package repository

import (
	"game_assistantor/model"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var UserRepos UserRepository

type UserRepository struct {
}

// 更新密码
func (u *UserRepository) SaveUserPassword(user *model.User) (err error) {
	hashStr, err := bcrypt.GenerateFromPassword([]byte(user.PassWord), bcrypt.DefaultCost) //加密处理
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	log.Info().Msgf("hash string is: %s", hashStr)
	err = engine.Model(&model.User{}).Where("user_id = ?", user.UserId).Update("password", string(hashStr)).Error
	return
}

// 获取单个用户
func (u *UserRepository) GetUser(userId string) (user *model.User, err error) {
	user = new(model.User)
	user.UserId = userId
	engine.First(user)
	return
}

// 更新数据
func (u *UserRepository) UpdateUser(user *model.User) (err error) {
	err = engine.Updates(user).Error
	return
}

// 获取所有用户
func (u *UserRepository) GetUserList() (userList []*model.User, err error) {
	err = engine.Model(model.User{}).Find(&userList).Error
	return
}
