package dao

import (
	"context"
	"gin_mall_tmp/model"
	"gorm.io/gorm"
)

type OrderDao struct {
	*gorm.DB
}

func NewOrderDao(ctx context.Context) *OrderDao {
	return &OrderDao{NewDBClient(ctx)}
}

func NewOrderDaoByDB(db *gorm.DB) *OrderDao {
	return &OrderDao{db}
}

func (dao *OrderDao) CreateOrder(in *model.Order) error {
	return dao.DB.Model(&model.Order{}).Create(&in).Error
}

func (dao *OrderDao) ShowOrderByOid(oId uint, uId uint) (order *model.Order, err error) {
	err = dao.DB.Model(&model.Order{}).Where("id=? ANd user_id=?", oId, uId).First(&order).Error
	return
}

func (dao *OrderDao) ListOrderByUserId(uId uint) (orders []*model.Order, err error) {
	err = dao.DB.Model(&model.Order{}).Where("user_id=?", uId).Find(&orders).Error
	return
}

func (dao *OrderDao) UpdateOrderByUserId(oId uint, uId uint, order *model.Order) error {
	return dao.DB.Model(&model.Order{}).Where("id=? AND user_id=?", oId, uId).Updates(&order).Error
}

func (dao *OrderDao) DeleteOrderByOrderId(oId uint, uId uint) error {
	return dao.DB.Model(&model.Order{}).Where("id=? AND user_id=?", oId, uId).Delete(&model.Order{}).Error
}

func (dao *OrderDao) UpdateOrderNumByUserId(oId uint, uId uint, num int) error {
	return dao.DB.Model(&model.Order{}).Where("id=? AND user_id=?", oId, uId).Update("num", num).Error
}

func (dao *OrderDao) ListOrderCondition(condition map[string]interface{}, page model.BasePage) (order []*model.Order, total int64, err error) {
	err = dao.DB.Model(&model.Order{}).Where(condition).Count(&total).Error
	err = dao.DB.Model(&model.Order{}).
		Where(condition).
		Offset((page.PageNum - 1) * (page.PageSize)).
		Limit(page.PageSize).Find(&order).Error
	return
}

func (dao *OrderDao) UpdateOrderByOrderId(oId uint, order *model.Order, uId uint) error {
	return dao.DB.Model(&model.Order{}).Where("id = ? AND user_id=?", oId, uId).Updates(&order).Error
}
