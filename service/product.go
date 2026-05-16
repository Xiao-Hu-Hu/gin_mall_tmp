package service

import (
	"context"
	"gin_mall_tmp/cache"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
	"gin_mall_tmp/pkg/e"
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/serializer"
	"mime/multipart"
	"strconv"
	"sync"
)

type ProductService struct {
	Id            uint   `json:"id" form:"id"`
	Name          string `json:"name" form:"name"`
	CategoryId    uint   `json:"category_id" form:"category_id"`
	Title         string `json:"title" form:"title"`
	Info          string `json:"info" form:"info"`
	ImgPath       string `json:"img_path" form:"img_path"`
	Price         string `json:"price" form:"price"`
	DiscountPrice string `json:"discount_price" form:"discount_price"`
	OnSale        bool   `json:"on_sale" form:"on_sale"`
	Num           int    `json:"num" form:"num"`
	model.BasePage
}

func (service *ProductService) Create(ctx context.Context, uId uint, files []*multipart.FileHeader) serializer.Response {
	var boss *model.User
	var err error
	code := e.Success

	// 获取商家信息
	userDao := dao.NewUserDao(ctx)
	boss, _ = userDao.GetUserById(uId)

	// 上传商品封面图以第一张图片最为封面图
	temp, _ := files[0].Open()
	path, err := UploadProductToLocalStatic(temp, uId, service.Name)
	if err != nil {
		code = e.ErrorProductImgUpload
		util.LogrusObj.Infoln("UploadProductToLocalStatic err:", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	product := &model.Product{
		Name:          service.Name,
		CategoryId:    service.CategoryId,
		Title:         service.Title,
		Info:          service.Info,
		ImgPath:       path,
		Price:         service.Price,
		DiscountPrice: service.DiscountPrice,
		OnSale:        true,
		Num:           service.Num,
		BossID:        uId,
		BossName:      boss.UserName,
		BossAvatar:    boss.Avatar,
	}

	// 创建商品主记录
	productDao := dao.NewProductDao(ctx)
	err = productDao.CreateProduct(product)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("ProductCreat err", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	// 并发上传商品详情图片
	wg := new(sync.WaitGroup)
	wg.Add(len(files))
	for index, file := range files {
		num := strconv.Itoa(index)
		productImgDao := dao.NewProductImgDaoByDB(productDao.DB)
		// 打开并上传每张图片
		tmp, _ := file.Open()
		path, err = UploadProductToLocalStatic(tmp, uId, service.Name+num)
		if err != nil {
			code = e.ErrorProductImgUpload
			util.LogrusObj.Infoln("UploadProductToLocalStatic err", err)
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		// 创建商品图片记录
		productImg := model.ProductImg{
			ProductId: product.ID, // 图片ID
			ImgPath:   path,       // 图片路径
		}
		err = productImgDao.CreatProductImg(&productImg)
		if err != nil {
			code = e.Error
			util.LogrusObj.Infoln("CreatProductImg err", err)
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		wg.Done()
	}
	wg.Wait() // 等待所有商品图片上传完成

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildProduct(product),
	}
}

func (service *ProductService) List(ctx context.Context) serializer.Response {
	var products []*model.Product
	var err error
	code := e.Success
	if service.PageSize == 0 {
		service.PageSize = 15
	}
	condition := make(map[string]interface{})
	if service.CategoryId != 0 {
		condition["category_id"] = service.CategoryId
	}
	productDao := dao.NewProductDao(ctx)
	total, err := productDao.CountProductByCondition(condition)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("ListProduct err", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		productDao = dao.NewProductDaoByDB(productDao.DB)
		products, _ = productDao.ListProductByCondition(condition, service.BasePage)
		wg.Done()
	}()
	wg.Wait()

	return serializer.BuildListResponse(serializer.BuildProducts(products), uint(total))
}

func (service *ProductService) Search(ctx context.Context) serializer.Response {
	code := e.Success
	if service.PageSize == 0 {
		service.PageSize = 15
	}
	productDao := dao.NewProductDao(ctx)
	products, count, err := productDao.SearchProduct(service.Info, service.BasePage)
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("SearchProduct err", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildProducts(products), uint(count))
}

func (service *ProductService) Show(ctx context.Context, id string) serializer.Response {
	code := e.Success
	pId, _ := strconv.Atoi(id)
	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(uint(pId))
	if err != nil {
		code = e.Error
		util.LogrusObj.Infoln("GetProductById err", err)
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// 更新商品热度
	cache.UpdateProductHeat(uint(pId), "view")

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildProduct(product),
	}
}
