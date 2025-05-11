package model

import (
	"time"
)

// MemberRequest 接收请求的结构体（小驼峰参数）
type MemberRequest struct {
	UID         int     `json:"uid"`
	Username    *string `json:"userName"` // 小驼峰
	Password    *string `json:"password"`
	Name        *string `json:"name,omitempty"`
	Tel         *string `json:"tel,omitempty"`
	Gender      *int    `json:"gender,omitempty"`
	ClassGrade  *string `json:"classGrade,omitempty"` // 小驼峰
	Team        *string `json:"team,omitempty"`
	Portrait    *string `json:"portrait,omitempty"`
	MienImg     *string `json:"mienImg,omitempty"` // 小驼峰
	Company     *string `json:"company,omitempty"`
	GraduateImg *string `json:"graduateImg,omitempty"` // 小驼峰
	IsGraduate  *int    `json:"isGraduate,omitempty"`  // 小驼峰
	Signature   *string `json:"signature,omitempty"`
	Year        *int    `json:"year,omitempty"`
	Status      *int    `json:"status,omitempty"`
	CaptchaID   *string `json:"captcha_id,omitempty"`
	CaptchaData *string `json:"captcha_data,omitempty"`
}

// MemberPO 数据库实体结构体
type MemberPO struct {
	UID         int       `gorm:"primaryKey;column:uid" json:"uid"`
	Username    *string   `gorm:"primaryKey;column:username" json:"userName"` // 小驼峰
	Password    *string   `gorm:"column:password" json:"password"`
	Name        *string   `gorm:"column:name" json:"name,omitempty"`
	Tel         *string   `gorm:"column:tel" json:"tel,omitempty"`
	Gender      *int      `gorm:"column:gender" json:"gender,omitempty"`
	ClassGrade  *string   `gorm:"column:class_grade" json:"classGrade,omitempty"` // 小驼峰
	Team        *string   `gorm:"column:team" json:"team,omitempty"`
	Portrait    *string   `gorm:"column:portrait" json:"portrait,omitempty"`
	MienImg     *string   `gorm:"column:mien_img" json:"mienImg,omitempty"` // 小驼峰
	Company     *string   `gorm:"column:company" json:"company,omitempty"`
	GraduateImg *string   `gorm:"column:graduate_img" json:"graduateImg,omitempty"` // 小驼峰
	IsGraduate  *int      `gorm:"column:is_graduate" json:"isGraduate,omitempty"`   // 小驼峰
	Signature   *string   `gorm:"column:signature" json:"signature,omitempty"`
	Year        *int      `gorm:"column:year" json:"year,omitempty"`
	Status      *int      `gorm:"column:status" json:"status,omitempty"`
	ModifyTime  time.Time `gorm:"column:modify_time;autoUpdateTime" json:"modifyTime,omitempty"`
}

// TableName 指定表名
func (*MemberPO) TableName() string {
	return "member"
}

// UpdateMemberRequest 定义更新用户信息的请求体
type UpdateMemberRequest struct {
	UID        int     `json:"uid"`                  // 主键，必填
	Portrait   *string `json:"portrait,omitempty"`   // 用户头像
	ClassGrade *string `json:"classGrade,omitempty"` // 班级
	Phone      *string `json:"phone,omitempty"`      // 电话
	Username   *string `json:"username,omitempty"`   // 用户名，可选
	Name       *string `json:"name,omitempty"`       // 姓名
	Team       *string `json:"team,omitempty"`       // 团队
	MienImg    *string `json:"mienImg,omitempty"`    // 风采图片
	Signature  *string `json:"signature,omitempty"`  // 个性签名
}
