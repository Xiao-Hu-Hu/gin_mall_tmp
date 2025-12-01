package service

import (
	"context"
	"fmt"
	"gin_mall_tmp/conf"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
	"gin_mall_tmp/pkg/e"
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/serializer"
	"mime/multipart"
	"strings"
	"time"

	"gopkg.in/mail.v2"
)

// 用户服务层
type UserService struct {
	NickName string `json:"nick_name" form:"nick_name"`
	UserName string `json:"user_name" form:"user_name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Key      string `json:"key" form:"key"` //密钥 前端验证
}

// 邮箱服务层
type SendEmailService struct {
	Email         string `json:"email" form:"email"`
	Password      string `json:"password" form:"password"`
	OperationType uint   `json:"operation_type" form:"operation_type"`
	// 1.绑定邮箱
	// 2.解绑邮箱
	// 3.修改密码
}

type ValidEmailService struct {
}

type ShowMoneyService struct {
	Key string `json:"key" form:"key"`
}

// Register用户注册
func (service *UserService) Register(ctx context.Context) serializer.Response {
	var user model.User
	code := e.Success
	if service.Key == "" || len(service.Key) != 16 {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "密钥长度不足",
		}
	}

	// 基础邮箱校验：必填
	if service.Email == "" {
		code = e.InvalidParams
		return serializer.Response{
			Status: code,
			Msg:    "邮箱不能为空",
			Error:  "email is empty",
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
	// 根据邮箱判断是否已被注册
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
			Msg:    "邮箱已被注册",
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
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// Login用户登录
func (service *UserService) Login(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	// 判断用户是否存在
	user, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if err != nil || !exist {
		code = e.ErrorExistUserNotFound
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户不存在，请先进行注册",
		}
	}
	// 校验密码
	if user.CheckPassword(service.Password) == false {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密码错误，请重新登录",
		}
	}

	// token 签发
	token, err := util.GenerateToken(user.ID, service.UserName, 0)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "token认证失败，请重新登录",
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.TokenData{User: serializer.BuildUser(user), Token: token},
	}

}

// Update用户修改信息
func (service *UserService) Update(ctx context.Context, uId uint) serializer.Response {
	var user *model.User
	var err error
	code := e.Success
	// 找到这个用户
	userDao := dao.NewUserDao(ctx)
	user, err = userDao.GetUserById(uId)
	// 修改昵称nickname
	if service.NickName != "" {
		user.NickName = service.NickName
	}
	err = userDao.UpdateUserById(uId, user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

// Post上传头像
func (service *UserService) Post(ctx context.Context, uId uint, file multipart.File, fileSize int64) serializer.Response {
	code := e.Success
	var user *model.User
	var err error
	userDao := dao.NewUserDao(ctx)
	user, err = userDao.GetUserById(uId)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	// 保存图片到本地
	path, err := UploadAvatarToLocalStatic(file, uId, user.UserName)
	if err != nil {
		code = e.ErrorUploadFail
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	user.Avatar = path
	err = userDao.UpdateUserById(uId, user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

// Send发送邮件
func (service *SendEmailService) Send(ctx context.Context, uId uint) serializer.Response {
	code := e.Success
	var address string
	var notice *model.Notice // 绑定邮箱，修改密码  模板通知

	// 1. 生成邮箱验证token
	// 根据用户ID、操作类型、邮箱和密码生成一个加密的token
	// 这个token会包含在邮件链接中，用于后续验证
	token, err := util.GenerateEmailToken(uId, service.OperationType, service.Email, service.Password)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 2. 获取邮件模板
	// 根据操作类型（如绑定邮箱、修改密码）从数据库获取对应的邮件模板
	noticeDao := dao.NewNoticeDao(ctx)
	notice, err = noticeDao.GetNoticeById(service.OperationType)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 3. 构建完整的验证链接
	// 将生成的token拼接到验证邮箱的基础URL上
	address = conf.ValidEmail + token // 发送方

	// 4. 处理邮件模板内容
	// 从数据库获取的邮件模板文本
	mailStr := notice.Text

	// 将模板中的"Email"占位符替换为实际的验证链接
	mailText := strings.Replace(mailStr, "Email", address, -1)
	// 5. 创建邮件消息
	m := mail.NewMessage()
	// 设置发件人
	m.SetHeader("From", conf.SmptEmail)
	// 设置收件人（用户输入的邮箱）
	m.SetHeader("To", service.Email)
	// 设置邮箱主题
	m.SetHeader("Subject", "XiaoHu")
	// 设置邮箱正文（html格式）
	m.SetBody("text/html", mailText)

	// 6. 配置邮件发送器
	// 使用SMTP服务器配置创建拨号器
	// 参数：SMTP主机、端口(465)、邮箱账号、邮箱密码或授权码
	d := mail.NewDialer(conf.SmtpHost, 465, conf.SmptEmail, conf.SmptPass)

	// 设置强制使用TLS加密连接，确保通信安全
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// 7. 发送邮件
	// 建立连接并发送邮件
	if err := d.DialAndSend(m); err != nil {
		code = e.ErrorSendEmail
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 8. 返回成功响应
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// Valid验证邮箱
func (service *ValidEmailService) Valid(ctx context.Context, emailToken string) serializer.Response {
	code := e.Success
	var userId uint
	var email string
	var password string
	var operationType uint
	// 验证token
	if emailToken == "" {
		fmt.Println("emailToken:", emailToken)
		code = e.InvalidParams
	} else {
		emailClaims, err := util.ParseEmailToken(emailToken)
		if err != nil {
			code = e.ErrorAuthToken
		} else if time.Now().Unix() > emailClaims.ExpiresAt.Unix() {
			code = e.ErrorAuthCheckTokenTimeOut
		} else {
			userId = emailClaims.UserID
			email = emailClaims.Email
			password = emailClaims.Password
			operationType = emailClaims.OpeartionType
		}
	}
	if code != e.Success {
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	// 获取用户的信息
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(userId)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if operationType == 1 {
		// 绑定邮箱
		user.Email = email
	} else if operationType == 2 {
		// 解绑邮箱
		user.Email = ""
	} else if operationType == 3 {
		err = user.SetPassword(password)
		if err != nil {
			code = e.Error
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
	}
	err = userDao.UpdateUserById(userId, user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

// Show显示金额
func (service *ShowMoneyService) Show(ctx context.Context, uId uint) serializer.Response {
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(uId)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildMoney(user, service.Key),
		Msg:    e.GetMsg(code),
	}
}
