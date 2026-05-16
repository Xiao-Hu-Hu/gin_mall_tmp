package dao

import (
	"context"
	"fmt"
	"gin_mall_tmp/model"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

func NewUserDaoByDB(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

// ExisDrNotByUserName 根据username判断是否存储在该名字
func (dao *UserDao) ExistOrNotByUserName(userName string) (user *model.User, exist bool, err error) {
	err = dao.DB.Model(&model.User{}).Where("user_name=?", userName).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("用户不存在", userName)
			return nil, false, nil
		}
		fmt.Println("数据库查询有误", err)
		return nil, false, err
	}
	return user, true, nil
}

// 创建用户
func (dao *UserDao) CreateUser(user *model.User) error {
	return dao.DB.Model(&model.User{}).Create(user).Error
}

// GetUserById 根据ID获取user
func (dao *UserDao) GetUserById(id uint) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).Where("id = ?", id).First(&user).Error
	return
}

// UpdateUserById 通过id更新user信息
func (dao *UserDao) UpdateUserById(id uint, user *model.User) (err error) {
	err = dao.DB.Model(&model.User{}).Where("id = ?", id).Updates(&user).Error
	return err
}

func (dao *UserDao) GetBossByProductId(pid uint) (boss *model.User, err error) {
	var bossId uint
	err = dao.DB.Model(&model.Product{}).Select("boss_id").Where("id = ?", pid).First(&bossId).Error
	if err != nil {
		return nil, err
	}
	err = dao.DB.Model(&model.User{}).Where("id = ?", bossId).First(&boss).Error
	return boss, err
}
