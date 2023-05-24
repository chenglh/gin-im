package model

import "time"

type Community struct {
	Id        int64      `gorm:"primary_key;comment:群ID" json:"id"`                                           // 自增ID/群ID
	Name      string     `gorm:"type:varchar(30);comment:群名称" json:"name"`                                    // 名称
	Ownerid   int64      `gorm:"type:int;default:0;comment:群主ID" json:"ownerid"`                              // 群主ID
	Icon      string     `gorm:"type:varchar(250);comment:群Logo" json:"icon"`                                 // 群logo
	Cate      int        `gorm:"type:tinyint(4);default:0;comment:群的类型" json:"cate"`                          // 群的类型
	Memo      string     `gorm:"type:varchar(120);comment:群描述" json:"memo"`                                   // 群描述
	CreatedAt *time.Time `gorm:"<-:create;type:datetime;index:idx_created_at;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt *time.Time `gorm:"<-;type:datetime;index:idx_updated_at;comment:更新时间" json:"updated_at"`        // 更新时间
	DeletedAt *time.Time `gorm:"type:datetime;index:idx_deleted_at;comment:删除时间" json:"deleted_at"`           // 删除时间
}
