package service

import (
	"errors"

	"awesomeProject/dao"
	"awesomeProject/model"
	"awesomeProject/utils"
)

func Register(mem *model.MemberRequest) error {
	if mem.Username == nil || mem.Password == nil || mem.CaptchaID == nil || mem.CaptchaData == nil {
		return errors.New("register 参数缺失")
	}
	// 接收UserName和Password 然后持久化到数据库
	memPO := &model.MemberPO{
		Username: mem.Username,
		Password: mem.Password,
	}
	// 验证码检查部分
	isValid := utils.CaptchaStore.Verify(*mem.CaptchaID, *mem.CaptchaData, true)
	if !isValid {
		return errors.New("验证码验证失败")
	}
	err := dao.InsertMember(memPO)
	if err != nil {
		return err
	}
	return nil
}
