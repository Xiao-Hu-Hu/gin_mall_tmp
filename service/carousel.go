package service

import (
	"context"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/pkg/e"
	"gin_mall_tmp/serializer"
)

type CarouselService struct {
}

func (service *CarouselService) List(context context.Context) serializer.Response {
	carouselDao := dao.NewCarouselDao(context)
	code := e.Success
	carousels, err := carouselDao.ListCarousel()
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildCarousels(carousels), uint(len(carousels)))
}
