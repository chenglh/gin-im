package service

import (
	"IM/gin_im/global"
	"IM/gin_im/model"
	"context"
	"errors"
	"fmt"
	"github.com/golang-module/carbon"
	"gorm.io/gorm"
)

// 好友/群服务
type ContactService struct{}

// 添加用户
func (s *ContactService) CreateFriend(userId, dstobj int) error {
	// 判断用户是否存在
	var userIds []int64
	result := global.DB.Model(&model.User{}).Where("user_id in(?,?)", userId, dstobj).Pluck("user_id", &userIds)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	} else if result.RowsAffected != 2 {
		return errors.New("用户信息错误")
	}

	// 判断是否已经添加
	var contact, contact2 model.Contact
	result = global.DB.Model(&model.Contact{}).Where("ownerid=? and dstobj=? and cate=0", userId, dstobj).First(&contact)
	if result.RowsAffected != 0 {
		return errors.New("当前账号已经是您的好友")
	}

	// 写入记录表
	nowTime := carbon.Now().Carbon2Time()
	contact.CreatedAt = &nowTime
	contact.UpdatedAt = &nowTime
	contact.Ownerid = int64(userId)
	contact.Dstobj = int64(dstobj)
	contact.Cate = 0

	// 对端好友
	contact2.CreatedAt = &nowTime
	contact2.UpdatedAt = &nowTime
	contact2.Ownerid = int64(dstobj)
	contact2.Dstobj = int64(userId)
	contact2.Cate = 0

	// 保存记录(事务)
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&contact).Error; err != nil {
			return err
		}
		if err := tx.Create(&contact2).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

// 好友列表
func (s *ContactService) FriendList(userId int64) ([]*model.User, error) {
	var friendUids []int64
	var userList []*model.User

	// 查询好友列表
	result := global.DB.Model(&model.Contact{}).Where("ownerid=? and cate=0", userId).Pluck("dstobj", &friendUids)
	if result.Error != nil || result.RowsAffected == 0 {
		return userList, result.Error
	}

	// 查询好友昵称
	result = global.DB.Model(&model.User{}).Where("user_id in(?)", friendUids).Find(&userList)

	return userList, result.Error
}

// 添加群
func (s *ContactService) AddCommunity(userId, dstobj int) error {
	// 判断用户是否存在
	var userInfo model.User
	result := global.DB.Model(&model.User{}).Where("user_id=?", userId).First(&userInfo)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	// 判断群是否存在
	var community model.Community
	result = global.DB.Model(&model.Community{}).Where("id=?", dstobj).First(&community)
	if community.Id == 0 {
		return errors.New("没有找到群ID")
	}

	// 判断是否已经添加
	var contact model.Contact
	result = global.DB.Model(&model.Contact{}).Where("ownerid=? and dstobj=? and cate=1", userId, dstobj).First(&contact)
	if result.RowsAffected != 0 {
		return errors.New("你已经加入了该群")
	}

	// 写入记录表
	nowTime := carbon.Now().Carbon2Time()
	contact.CreatedAt = &nowTime
	contact.UpdatedAt = &nowTime
	contact.Ownerid = int64(userId)
	contact.Dstobj = int64(dstobj)
	contact.Cate = 1

	// 保存记录(事务)
	if err := global.DB.Create(&contact).Error; err != nil {
		return err
	}

	return nil
}

// 创建群
func (s *ContactService) CreateCommunity(community model.Community, contact model.Contact) (int64, error) {
	// 判断群主ID是否存在
	var userInfo model.User
	result := global.DB.Model(&model.User{}).Where("user_id=?", community.Ownerid).First(&userInfo)
	if result.RowsAffected == 0 {
		return 0, errors.New("群主ID不存在")
	}

	// 需要使用事务，闭包事务
	// 1、创建群
	// 2、群与群主绑定
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作
		if err := tx.Create(&community).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}

		contact.Dstobj = community.Id // 填充群ID
		if err := tx.Create(&contact).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})

	return community.Id, err
}

// 用户群列表
func (s *ContactService) CommunityList(userId int) ([]*model.Community, error) {
	var contactIds []int64
	var communityList []*model.Community

	// 查询所属群
	result := global.DB.Model(&model.Contact{}).Where("ownerid=? and cate=1", userId).Pluck("dstobj", &contactIds)
	if result.Error != nil || result.RowsAffected == 0 {
		return communityList, result.Error
	}

	// 查询群信息
	result = global.DB.Model(&model.Community{}).Where("id in(?)", contactIds).Find(&communityList)

	return communityList, result.Error
}

// 聊天记录
func (s *ContactService) GetMessageList(userIdA, userIdB, start, end, msgType int64, isRev bool) ([]string, int64) {
	var key string
	if msgType == 1 {
		if userIdA > userIdB {
			key = fmt.Sprintf("single_msg_%v_%v", userIdA, userIdB)
		} else {
			key = fmt.Sprintf("single_msg_%v_%v", userIdB, userIdA)
		}
	} else {
		key = fmt.Sprintf("qunliao_msg_%v", userIdB)
	}

	fmt.Println(key)

	var rels []string
	ctx := context.Background()
	if isRev {
		rels, _ = global.Rdb.ZRange(ctx, key, start, end).Result()
	} else {
		rels, _ = global.Rdb.ZRevRange(ctx, key, start, end).Result()
	}

	total, _ := global.Rdb.ZCard(ctx, key).Result()

	return rels, total
}
