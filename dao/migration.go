package dao

import "gin_mall_tmp/model"

func Migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&model.User{},
			&model.Address{},
			&model.Admin{},
			&model.Carousel{},
			&model.Category{},
			&model.BasePage{},
			&model.Cart{},
			&model.Favorite{},
			&model.Notice{},
			&model.Product{},
			&model.Order{},
			&model.ProductImg{},
			&model.SeckillActivity{},
			&model.SeckillOrder{},
		)
	if err != nil {
		panic(err)
		return
	}

	return
}
