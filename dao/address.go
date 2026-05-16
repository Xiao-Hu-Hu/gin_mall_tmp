package dao

import (
	"context"
	"gin_mall_tmp/model"
	"gorm.io/gorm"
)

type AddressDao struct {
	*gorm.DB
}

func NewAddressDao(ctx context.Context) *AddressDao {
	return &AddressDao{NewDBClient(ctx)}
}

func NewAddressDaoByDB(db *gorm.DB) *AddressDao {
	return &AddressDao{db}
}

func (dao *AddressDao) CreateAddress(in *model.Address) error {
	return dao.DB.Model(&model.Address{}).Create(&in).Error
}

func (dao *AddressDao) GetAddressByAid(aId uint, uId uint) (address *model.Address, err error) {
	err = dao.DB.Model(&model.Address{}).Where("id=? ANd user_id=?", aId, uId).First(&address).Error
	return
}

func (dao *AddressDao) ListAddressByUserId(uId uint) (addresses []*model.Address, err error) {
	err = dao.DB.Model(&model.Address{}).Where("user_id=?", uId).Find(&addresses).Error
	return
}

func (dao *AddressDao) UpdateAddressByUserId(aId uint, uId uint, address *model.Address) error {
	return dao.DB.Model(&model.Address{}).Where("id=? AND user_id=?", aId, uId).Updates(&address).Error
}

func (dao *AddressDao) DeleteAddressByAddressId(aId uint, uId uint) error {
	return dao.DB.Model(&model.Address{}).Where("id=? AND user_id=?", aId, uId).Delete(&model.Address{}).Error
}
