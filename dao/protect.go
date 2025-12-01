package dao

import (
	"context"
	"gin_mall_tmp/model"

	"gorm.io/gorm"
)

type ProductDao struct {
	*gorm.DB
}

func NewProductDao(ctx context.Context) *ProductDao {
	return &ProductDao{NewDBClient(ctx)}
}

func NewProductDaoByDB(db *gorm.DB) *ProductDao {
	return &ProductDao{db}
}

func (dao *ProductDao) CreateProduct(product *model.Product) (err error) {
	return dao.DB.Model(&model.Product{}).Create(&product).Error
}

func (dao *ProductDao) CountProductByCondition(condition map[string]interface{}) (total int64, err error) {
	err = dao.DB.Model(&model.Product{}).Where(condition).Count(&total).Error
	return total, err
}

func (dao *ProductDao) ListProductByCondition(condition map[string]interface{}, page model.BasePage) (products []*model.Product, err error) {
	err = dao.DB.Where(condition).Offset((page.PageNum - 1) * (page.PageSize)).Limit(page.PageSize).Find(&products).Error
	return products, err
}

func (dao *ProductDao) SearchProduct(info string, page model.BasePage) (products []*model.Product, count int64, err error) {
	// 支持按名称、标题、描述进行关键词模糊匹配
	like := "%" + info + "%"
	err = dao.DB.Model(&model.Product{}).
		Where("name LIKE ? OR title LIKE ? OR info LIKE ?", like, like, like).
		Count(&count).Error
	if err != nil {
		return
	}

	err = dao.DB.Model(&model.Product{}).
		Where("name LIKE ? OR title LIKE ? OR info LIKE ?", like, like, like).
		Offset((page.PageNum - 1) * (page.PageSize)).
		Limit(page.PageSize).Find(&products).Error
	return
}

func (dao *ProductDao) GetProductById(id uint) (product *model.Product, err error) {
	err = dao.DB.Model(&model.Product{}).Where("id = ?", id).First(&product).Error
	return
}

func (dao *ProductDao) UpdateProductById(pId uint, product *model.Product) (err error) {
	return dao.DB.Model(&model.Product{}).Where("id = ?", pId).Updates(&product).Error
}
