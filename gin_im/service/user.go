package service

import (
	"IM/gin_im/global"
	"IM/gin_im/model"
	"errors"
)

type UserService struct{}

// 通过手机号查询用户
func (s *UserService) GetUserByMobile(mobile string) (*model.User, error) {
	var userInfo model.User

	result := global.DB.Model(model.User{}).Where("mobile=?", mobile).First(&userInfo)
	if result.RowsAffected == 0 {
		return &userInfo, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &userInfo, nil
}

// 通过ID查询用户
func (s *UserService) GetUserByUserId(userId int64) (*model.User, error) {
	var userInfo model.User

	result := global.DB.Model(model.User{}).First(&userInfo, userId)
	if result.RowsAffected == 0 {
		return &userInfo, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &userInfo, nil
}

// 创建用户
func (s *UserService) CreateUser(userInfo *model.User) (userId int64, err error) {
	//判断用户是否存在
	var user model.User
	result := global.DB.Where(model.User{Mobile: userInfo.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return 0, errors.New("用户已存在")
	}
	if result := global.DB.Create(&userInfo); result.Error != nil {
		return 0, result.Error
	}

	return userInfo.UserId, nil
}

// 更新用户
func (s *UserService) UpdateUser(userInfo *model.User) error {
	if err := global.DB.Save(userInfo).Error; err != nil {
		return err
	}

	return nil
}
