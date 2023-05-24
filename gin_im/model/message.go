package model

import (
	"time"
)

type Message struct {
	ID        int64      `gorm:"primary_key;comment:自增ID" json:"id"`                                          // 自增ID
	Ownerid   int64      `gorm:"type:int;default:0;comment:用户ID" json:"ownerid"`                              // 发送者
	Dstobj    int64      `gorm:"type:int;default:0;comment:对端ID" json:"dstobj"`                               // 接收者:用户ID/群ID
	MsgType   int32      `gorm:"type:tinyint;default:0;comment:聊天类型" json:"msg_type"`                         // 聊天类型:1私聊;2群聊;3心跳
	Media     int32      `gorm:"type:tinyint;default:0;comment:消息类型" json:"media"`                            // 消息类型:1文本;2表情包;3语音;4图片
	Content   string     `gorm:"type:tinyint;default:0;comment:消息类型" json:"content"`                          // 消息内容
	Picture   string     `gorm:"type:varchar(255);comment:预览图片" json:"picture"`                               // 预览图片
	Url       string     `gorm:"type:varchar(255);comment:备注信息" json:"url"`                                   // URL连接
	Desc      string     `gorm:"type:varchar(500);comment:备注信息" json:"desc"`                                  // 简单描述
	Amount    int32      `gorm:"type:int;default:0;comment:金额数据" json:"amount"`                               // 红包金额,单位分;语音时长
	TpCount   int32      `gorm:"type:int;default:0;comment:其他统计" json:"tp_count"`                             // 其他数据统计,如红包个数
	CreatedAt *time.Time `gorm:"<-:create;type:datetime;index:idx_created_at;comment:创建时间" json:"created_at"` // 创建时间
	ReadAt    *time.Time `gorm:"<-;type:datetime;index:idx_updated_at;comment:读取时间" json:"read_at"`           // 信息读取时间
	UpdatedAt *time.Time `gorm:"<-;type:datetime;index:idx_updated_at;comment:更新时间" json:"updated_at"`        // 更新时间
	DeletedAt *time.Time `gorm:"type:datetime;index:idx_deleted_at;comment:删除时间" json:"deleted_at"`           // 删除时间
}
