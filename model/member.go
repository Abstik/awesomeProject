package model

import (
	"time"
)

// RegisterReq 注册请求
type RegisterReq struct {
	Username *string `json:"username" binding:"required"` // 用户名，中文拼音
	Name     *string `json:"name" binding:"required"`     // 中文姓名
	Year     *int    `json:"year" binding:"required"`     // 入学年份
	Team     *string `json:"team" binding:"required"`     // 组别
}

// LoginReq 登录请求
type LoginReq struct {
	Username    *string `json:"username" binding:"required"`    // 用户名
	Password    *string `json:"password" binding:"required"`    // 密码
	CaptchaID   *string `json:"captchaID" binding:"required"`   // 验证码ID
	CaptchaData *string `json:"captchaData" binding:"required"` // 验证码数据
}

// UpdateMemberRequest 修改用户信息请求
type UpdateMemberRequest struct {
	Username    *string `json:"username" binding:"required"` // 用户名
	Name        *string `json:"name"`                        // 姓名
	Tel         *string `json:"tel"`                         // 电话
	Gender      *int    `json:"gender"`                      // 性别，0为男生，1为女生
	ClassGrade  *string `json:"classGrade"`                  // 专业班级
	Team        *string `json:"team"`                        // 组别
	Portrait    *string `json:"portrait"`                    // 用户头像
	MienImg     *string `json:"mienImg"`                     // 风采图片
	Company     *string `json:"company"`                     // 入职公司
	GraduateImg *string `json:"graduateImg"`                 // 毕业照，有默认值
	IsGraduate  *int    `json:"isGraduate"`                  // 是否毕业，0未毕业，1已毕业，默认为0
	Signature   *string `json:"signature"`                   // 个性签名
	Year        *int    `json:"year"`                        // 入学年份
}

// MemberPO 数据库实体结构体
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
