package dao

import (
	"context"
	"gin_mall_tmp/model"
	"gorm.io/gorm"
)

type CategoryDao struct {
	*gorm.DB
}

func NewCategoryDao(ctx context.Context) *CategoryDao {
	return &CategoryDao{NewDBClient(ctx)}
}

func NewCategoryDaoByDB(db *gorm.DB) *CategoryDao {
	return &CategoryDao{db}
}

// GetCategoryById 根据id获取Category
func (dao *CategoryDao) ListCategory() (Category []*model.Category, err error) {
	err = dao.DB.Model(&model.Category{}).Find(&Category).Error
	return
}
