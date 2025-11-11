package service

import (
	"context"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
	"gin_mall_tmp/pkg/e"
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/serializer"
	"strconv"
)

type FavoriteService struct {
	ProductId  uint `json:"product_id" form:"product_id"`
	BossId     uint `json:"boss_id" form:"boss_id"`
	FavoriteId uint `json:"favorite_id" form:"favorite_id"`
	model.BasePage
}

func (service *FavoriteService) Show(ctx context.Context, uId uint) serializer.Response {
	favoriteDao := dao.NewFavoriteDao(ctx)
	code := e.Success
	favorites, err := favoriteDao.ShowFavorite(uId)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildFavorites(ctx, favorites), uint(len(favorites)))
}

func (service *FavoriteService) Create(ctx context.Context, uId uint) serializer.Response {
	code := e.Success
	favoriteDao := dao.NewFavoriteDao(ctx)
	exist, _ := favoriteDao.FavoriteExistOrNot(service.ProductId, uId)
	if exist {
		code = e.ErrorFavoriteExist
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  e.GetMsg(code),
		}
	}
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(uId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("CreatFavorite err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	bossDao := dao.NewUserDao(ctx)
	boss, err := bossDao.GetUserById(service.BossId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("CreatFavorite err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(service.ProductId)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("CreatFavorite err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	favorite := &model.Favorite{
		User:      *user,
		UserId:    uId,
		Product:   *product,
		ProductId: service.ProductId,
		Boss:      *boss,
		BossId:    service.BossId,
	}
	err = favoriteDao.CreateFavorite(favorite)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("CreatFavorite err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildFavorite(favorite, product, boss),
	}
}

func (service *FavoriteService) Delete(ctx context.Context, uId uint, fId string) serializer.Response {
	favoriteDao := dao.NewFavoriteDao(ctx)
	favoriteId, err := strconv.Atoi(fId)
	code := e.Success
	err = favoriteDao.DeleteFavorite(uId, uint(favoriteId))
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("DeleteFavorite err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
