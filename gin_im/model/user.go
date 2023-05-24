package model

import "time"

type User struct {
	UserId    int64      `gorm:"primary_key;comment:用户ID"`
	NickName  string     `gorm:"type:varchar(30);comment:用户昵称"`
	Status    int32      `gorm:"type:tinyint;default:0;comment:用户状态,0正常,1禁用,2风控"`
	Password  string     `gorm:"type:varchar(150);default:null;comment:密码串"`
	Token     string     `gorm:"type:varchar(255);default:nll;comment:密钥串"`
	Mobile    string     `gorm:"type:varchar(11);index:idx_mobile;unique;not null;comment:手机号码"`
	Gender    string     `gorm:"type:tinyint(4);default:0;comment:用户性别,0未知,1男,2女"`
	LoginIp   string     `gorm:"type:varchar(60);comment:登录IP"`
	LoginTime *time.Time `gorm:"type:datetime;comment:登录时间"`
	HeadUrl   string     `gorm:"type:varchar(255);comment:用户头像"`
	CreatedAt time.Time  `gorm:"<-:create;type:datetime;index:idx_created_at;comment:创建时间"`
	UpdatedAt time.Time  `gorm:"<-;type:datetime;index:idx_updated_at;comment:更新时间"`
	DeletedAt *time.Time `gorm:"type:datetime;index:idx_deleted_at;comment:删除时间"`
}

// TableName 指定表名
func (User) TableName() string {
	return "im_user"
}
