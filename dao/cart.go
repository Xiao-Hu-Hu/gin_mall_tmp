package dao

import (
	"context"
	"gin_mall_tmp/model"
	"gorm.io/gorm"
)

type CartDao struct {
	*gorm.DB
}

func NewCartDao(ctx context.Context) *CartDao {
	return &CartDao{NewDBClient(ctx)}
}

func NewCartDaoByDB(db *gorm.DB) *CartDao {
	return &CartDao{db}
}

func (dao *CartDao) CreateCart(in *model.Cart) error {
	return dao.DB.Model(&model.Cart{}).Create(&in).Error
}

func (dao *CartDao) ShowCartByCid(cId uint, uId uint) (cart *model.Cart, err error) {
	err = dao.DB.Model(&model.Cart{}).Where("id=? ANd user_id=?", cId, uId).First(&cart).Error
	return
}

func (dao *CartDao) ListCartByUserId(uId uint) (carts []*model.Cart, err error) {
	err = dao.DB.Model(&model.Cart{}).Where("user_id=?", uId).Find(&carts).Error
	return
}

func (dao *CartDao) UpdateCartByUserId(cId uint, uId uint, cart *model.Cart) error {
	return dao.DB.Model(&model.Cart{}).Where("id=? AND user_id=?", cId, uId).Updates(&cart).Error
}

func (dao *CartDao) DeleteCartByCartId(cId uint, uId uint) error {
	return dao.DB.Model(&model.Cart{}).Where("id=? AND user_id=?", cId, uId).Delete(&model.Cart{}).Error
}

func (dao *CartDao) UpdateCartNumByUserId(cId uint, uId uint, num int) error {
	return dao.DB.Model(&model.Cart{}).Where("id=? AND user_id=?", cId, uId).Update("num", num).Error
}
