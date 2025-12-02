package service

import (
	"context"
	"fmt"
	"gin_mall_tmp/cache"
	"gin_mall_tmp/conf"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
	"gin_mall_tmp/pkg/e"
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/serializer"
	"math/rand"
	"strconv"
	"time"

	"gopkg.in/mail.v2"
)

type EmailVerifyService struct {
	Email string `json:"email" form:"email"`
}

type EmailRegisterService struct {
	UserName string `json:"user_name" form:"user_name"`
	NickName string `json:"nick_name" form:"nick_name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Key      string `json:"key" form:"key"`
	Code     string `json:"code" form:"code"`
}

// SendVerifyCode 发送邮箱验证码
func (service *EmailVerifyService) SendVerifyCode(ctx context.Context) serializer.Response {
	code := e.Success

	// 验证邮箱格式
	if service.Email == "" {
		code = e.InvalidParams
		return serializer.Response{
			Status: code,
			Msg:    "邮箱不能为空",
		}
	}

	// 检查邮箱是否已被注册
	userDao := dao.NewUserDao(ctx)
	if _, exist, err := userDao.ExistOrNotByEmail(service.Email); err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	} else if exist {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "该邮箱已被注册",
		}
	}

	// 生成6位随机验证码
	verifyCode := generateVerifyCode()

	// 检查Redis是否可用
	if cache.RedisClient == nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "验证码服务暂不可用，请稍后重试",
		}
	}

	// 将验证码存入Redis，有效期3分钟
	err := cache.RedisClient.Set(ctx, "email_verify:"+service.Email, verifyCode, 3*time.Minute).Err()
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "验证码发送失败",
			Error:  err.Error(),
		}
	}

	// 发送验证码邮件
	err = sendVerifyEmail(service.Email, verifyCode)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "邮件发送失败",
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    "验证码已发送到您的邮箱，请查收",
	}
}

// RegisterWithEmail 邮箱验证码注册
func (service *EmailRegisterService) RegisterWithEmail(ctx context.Context) serializer.Response {
	var user model.User
	code := e.Success

	// 验证密钥
	if service.Key == "" || len(service.Key) != 16 {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "密钥长度必须为16位",
		}
	}

	// 验证验证码
	if service.Code == "" {
		code = e.InvalidParams
		return serializer.Response{
			Status: code,
			Msg:    "验证码不能为空",
		}
	}

	// 检查Redis是否可用
	if cache.RedisClient == nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "验证码服务暂不可用，请稍后重试",
		}
	}

	// 从Redis获取验证码
	storedCode, err := cache.RedisClient.Get(ctx, "email_verify:"+service.Email).Result()
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "验证码已过期或无效",
		}
	}

	// 验证验证码
	if storedCode != service.Code {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "验证码错误",
		}
	}

	// 验证邮箱格式
	if service.Email == "" {
		code = e.InvalidParams
		return serializer.Response{
			Status: code,
			Msg:    "邮箱不能为空",
		}
	}

	// 初始金额10000 ---> 密文存储 对称加密
	util.Encrypt.SetKey(service.Key)

	userDao := dao.NewUserDao(ctx)
	// 根据名字判断是否存在该用户
	_, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	if exist {
		code = e.ErrorExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	user = model.User{
		UserName: service.UserName,
		NickName: service.NickName,
		Email:    service.Email,
		Status:   model.Avtive,
		Avatar:   "avatar.JPG",
		Money:    util.Encrypt.AesEncoding("10000"), // 初始金额的加密
	}

	// 密码加密
	if err = user.SetPassword(service.Password); err != nil {
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	// 创建用户
	err = userDao.CreateUser(&user)
	if err != nil {
		code = e.Error
	}

	// 注册成功后删除验证码
	if cache.RedisClient != nil {
		cache.RedisClient.Del(ctx, "email_verify:"+service.Email)
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// generateVerifyCode 生成6位随机验证码
func generateVerifyCode() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(900000) + 100000)
}

// sendVerifyEmail 发送验证码邮件
func sendVerifyEmail(toEmail, code string) error {
	m := mail.NewMessage()
	m.SetHeader("From", conf.SmptEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "邮箱验证码")

	content := fmt.Sprintf(`
		<html>
		<body>
			<h2>欢迎注册我们的商城</h2>
			<p>您的验证码是：<strong style="font-size: 24px; color: #1890ff;">%s</strong></p>
			<p>验证码有效期为3分钟，请尽快使用。</p>
			<p>如果这不是您本人的操作，请忽略此邮件。</p>
		</body>
		</html>
	`, code)

	m.SetBody("text/html", content)

	// 使用SMTP服务器配置创建拨号器
	// 参数：SMTP主机、端口(465)、邮箱账号、邮箱密码或授权码
	d := mail.NewDialer(conf.SmtpHost, 465, conf.SmptEmail, conf.SmptPass)

	// 强制使用SSL/TLS
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
