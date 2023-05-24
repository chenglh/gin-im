package model

import "time"

// 好友和群都存在这个表里面
// 可根据具体业务做拆分
type Contact struct {
	Id        int64      `gorm:"primary_key;comment:自增ID" json:"id"`                                          // 自增ID
	Ownerid   int64      `gorm:"type:int;default:0;comment:用户ID" json:"ownerid"`                              // 记录是谁的
	Dstobj    int64      `gorm:"type:int;default:0;comment:对端ID" json:"dstobj"`                               // 对端用户ID/群ID
	Cate      int64      `gorm:"type:tinyint;default:0;comment:对端类型,0好友,1群聊" json:"cate"`                     // 什么类型
	Memo      string     `gorm:"type:varchar(30);comment:备注信息" json:"memo"`                                   // 备注信息
	CreatedAt *time.Time `gorm:"<-:create;type:datetime;index:idx_created_at;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt *time.Time `gorm:"<-;type:datetime;index:idx_updated_at;comment:更新时间" json:"updated_at"`        // 更新时间
	DeletedAt *time.Time `gorm:"type:datetime;index:idx_deleted_at;comment:删除时间" json:"deleted_at"`           // 删除时间
}
