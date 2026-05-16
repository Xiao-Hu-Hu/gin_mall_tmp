package dao

import (
	"context"
	"gin_mall_tmp/model"
	"gorm.io/gorm"
)

type SeckillDao struct {
	DB *gorm.DB
}

func NewSeckillDao(ctx context.Context) *SeckillDao {
	return &SeckillDao{DB: NewDBClient(ctx)}
}

// GetActivityByID 根据 ID 获取秒杀活动
func (dao *SeckillDao) GetActivityByID(id uint) (*model.SeckillActivity, error) {
	var activity model.SeckillActivity
	err := dao.DB.Where("id = ?", id).First(&activity).Error
	return &activity, err
}

// GetActiveActivityByID 获取进行中的秒杀活动
func (dao *SeckillDao) GetActiveActivityByID(id uint) (*model.SeckillActivity, error) {
	var activity model.SeckillActivity
	err := dao.DB.Where("id = ? AND status = 1", id).First(&activity).Error
	return &activity, err
}

// CreateSeckillOrder 创建秒杀订单
func (dao *SeckillDao) CreateSeckillOrder(order *model.SeckillOrder) error {
	return dao.DB.Create(order).Error
}

// GetSeckillOrderByUserAndActivity 根据用户和活动获取秒杀订单
func (dao *SeckillDao) GetSeckillOrderByUserAndActivity(userID, activityID uint) (*model.SeckillOrder, error) {
	var order model.SeckillOrder
	err := dao.DB.Where("user_id = ? AND activity_id = ?", userID, activityID).First(&order).Error
	return &order, err
}

// UpdateSeckillOrderStatus 更新秒杀订单状态
func (dao *SeckillDao) UpdateSeckillOrderStatus(orderID uint, status uint) error {
	return dao.DB.Model(&model.SeckillOrder{}).Where("id = ?", orderID).Update("status", status).Error
}

// ListSeckillOrdersByUser 获取用户的秒杀订单列表
func (dao *SeckillDao) ListSeckillOrdersByUser(userID uint, page, pageSize int) ([]*model.SeckillOrder, int64, error) {
	var orders []*model.SeckillOrder
	var total int64

	db := dao.DB.Where("user_id = ?", userID)
	db.Model(&model.SeckillOrder{}).Count(&total)

	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&orders).Error
	return orders, total, err
}
