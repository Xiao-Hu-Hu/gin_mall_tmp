package service

import (
	"context"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/pkg/e"
	"gin_mall_tmp/serializer"
)

type CategoryService struct {
}

func (service *CategoryService) List(ctx context.Context) serializer.Response {
	categoryDao := dao.NewCategoryDao(ctx)
	code := e.Success
	categories, err := categoryDao.ListCategory()
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildCategories(categories), uint(len(categories)))
}
