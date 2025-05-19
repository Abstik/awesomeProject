package model

import (
	"time"
)

// 接收请求的结构体
type MemberRequest struct {
	UID         int     `json:"uid"`
	Username    *string `json:"userName"`              // 用户名，中文拼音
	Password    *string `json:"password"`              // 密码
	Name        *string `json:"name,omitempty"`        // 中文用户名
	Tel         *string `json:"tel,omitempty"`         // 电话
	Gender      *int    `json:"gender,omitempty"`      // 性别
	ClassGrade  *string `json:"classGrade,omitempty"`  // 专业班级
	Team        *string `json:"team,omitempty"`        // 组别
	Portrait    *string `json:"portrait,omitempty"`    // 头像，有默认值
	MienImg     *string `json:"mienImg,omitempty"`     // 风采图片，有默认值
	Company     *string `json:"company,omitempty"`     // 入职公司
	GraduateImg *string `json:"graduateImg,omitempty"` // 毕业照，有默认值
	IsGraduate  *int    `json:"isGraduate,omitempty"`  // 是否毕业
	Signature   *string `json:"signature,omitempty"`   // 个性签名
	Year        *int    `json:"year,omitempty"`        // 入学年份
	Status      *int    `json:"status,omitempty"`      // 0管理员，1已毕业，2未毕业
	CaptchaID   *string `json:"captchaID,omitempty"`   // 验证码ID，必填
	CaptchaData *string `json:"captchaData,omitempty"` // 验证码数据，必填
}

// 更新用户信息的请求体
type UpdateMemberRequest struct {
	UID        int     `json:"uid"`                  // 主键，必填
	Portrait   *string `json:"portrait,omitempty"`   // 用户头像
	ClassGrade *string `json:"classGrade,omitempty"` // 专业班级
	Phone      *string `json:"phone,omitempty"`      // 电话
	Username   *string `json:"username,omitempty"`   // 用户名
	Name       *string `json:"name,omitempty"`       // 姓名
	Team       *string `json:"team,omitempty"`       // 组别
	MienImg    *string `json:"mienImg,omitempty"`    // 风采图片
	Signature  *string `json:"signature,omitempty"`  // 个性签名
}

// 数据库实体结构体
type MemberPO struct {
	UID         int       `gorm:"primaryKey;column:uid" json:"uid"` // 自增主键
	Username    *string   `gorm:"column:username" json:"username"`  // 非空，与UID构成联合唯一索引，用户名的中文拼音
	Password    *string   `gorm:"column:password" json:"-"`
	Name        *string   `gorm:"column:name" json:"name"`                // 中文用户名
	Tel         *string   `gorm:"column:tel" json:"tel"`                  // 手机号
	Gender      *int      `gorm:"column:gender" json:"gender"`            // 性别，0为男生，1为女生
	ClassGrade  *string   `gorm:"column:class_grade" json:"classGrade"`   // 专业班级
	Team        *string   `gorm:"column:team" json:"team"`                // 组别
	Portrait    *string   `gorm:"column:portrait" json:"portrait"`        // 肖像，有默认值
	MienImg     *string   `gorm:"column:mien_img" json:"mienImg"`         // 风采图片，有默认值
	Company     *string   `gorm:"column:company" json:"company"`          // 入职公司
	GraduateImg *string   `gorm:"column:graduate_img" json:"graduateImg"` // 毕业照，有默认值
	IsGraduate  *int      `gorm:"column:is_graduate" json:"isGraduate"`   // 是否毕业，0未毕业，1已毕业，默认为0
	Signature   *string   `gorm:"column:signature" json:"signature"`      // 个性签名
	Year        *int      `gorm:"column:year" json:"year"`                // 入学年份
	Status      *int      `gorm:"column:status" json:"status"`            // 0为管理员，1为用户，默认为1
	ModifyTime  time.Time `gorm:"column:modify_time;autoUpdateTime" json:"-"`
}

// TableName 指定表名
func (*MemberPO) TableName() string {
	return "member"
}
