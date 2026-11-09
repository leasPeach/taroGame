package dao

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"

	"gorm.io/gorm"
)

const userTableName = "Users"

// UserInterface 用户数据接口
type UserInterface interface {
	GetByOpenId(openId string) (*model.UserModel, error)
	UpsertByOpenId(user *model.UserModel) (*model.UserModel, error)
}

// UserInterfaceImp 用户数据实现
type UserInterfaceImp struct{}

// UserImp 用户实现实例
var UserImp UserInterface = &UserInterfaceImp{}

// GetByOpenId 按 openId 查询用户
func (imp *UserInterfaceImp) GetByOpenId(openId string) (*model.UserModel, error) {
	var user model.UserModel
	err := db.Get().Table(userTableName).Where("openId = ?", openId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpsertByOpenId 按 openId 写入或更新用户
func (imp *UserInterfaceImp) UpsertByOpenId(user *model.UserModel) (*model.UserModel, error) {
	cli := db.Get()
	existing, err := imp.GetByOpenId(user.OpenId)
	now := time.Now()

	if err == gorm.ErrRecordNotFound {
		user.CreatedAt = now
		user.UpdatedAt = now
		if err := cli.Table(userTableName).Create(user).Error; err != nil {
			return nil, err
		}
		return user, nil
	}
	if err != nil {
		return nil, err
	}

	if user.UnionId != "" {
		existing.UnionId = user.UnionId
	}
	if user.NickName != "" {
		existing.NickName = user.NickName
	}
	if user.AvatarUrl != "" {
		existing.AvatarUrl = user.AvatarUrl
	}
	existing.UpdatedAt = now
	if err := cli.Table(userTableName).Save(existing).Error; err != nil {
		return nil, err
	}
	return existing, nil
}

// IsNotFound 是否未找到记录
func IsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
