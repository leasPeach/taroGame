package model

import "time"

// UserModel 微信用户
type UserModel struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OpenId    string    `gorm:"column:openId;uniqueIndex;size:64;not null" json:"openId"`
	UnionId   string    `gorm:"column:unionId;size:64" json:"unionId,omitempty"`
	NickName  string    `gorm:"column:nickName;size:64" json:"nickName"`
	AvatarUrl string    `gorm:"column:avatarUrl;size:512" json:"avatarUrl"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// TableName 表名
func (UserModel) TableName() string {
	return "Users"
}
