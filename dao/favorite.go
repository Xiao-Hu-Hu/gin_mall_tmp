package dao

import (
	"context"
	"gin_mall_tmp/model"
	"gorm.io/gorm"
)

type FavoriteDao struct {
	*gorm.DB
}

func NewFavoriteDao(ctx context.Context) *FavoriteDao {
	return &FavoriteDao{NewDBClient(ctx)}
}

func NewFavoriteDaoByDB(db *gorm.DB) *FavoriteDao {
	return &FavoriteDao{db}
}

func (dao *FavoriteDao) ShowFavorite(uId uint) (favorites []*model.Favorite, err error) {
	err = dao.DB.Model(&model.Favorite{}).Where("user_id=?", uId).Find(&favorites).Error
	return
}

func (dao *FavoriteDao) FavoriteExistOrNot(pId uint, uId uint) (exist bool, err error) {
	var count int64
	err = dao.DB.Model(&model.Favorite{}).Where("product_id=? AND user_id=?", pId, uId).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (dao *FavoriteDao) CreateFavorite(in *model.Favorite) error {
	return dao.DB.Model(&model.Favorite{}).Create(in).Error
}

func (dao *FavoriteDao) DeleteFavorite(uId, fId uint) error {
	return dao.DB.Model(&model.Favorite{}).Where("id = ? AND user_id = ?", fId, uId).Delete(&model.Favorite{}).Error
}
