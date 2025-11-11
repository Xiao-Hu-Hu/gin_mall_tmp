package main

import (
	"gin_mall_tmp/conf"
	"gin_mall_tmp/routes"
)

func main() {
	conf.Init()
	r := routes.NewRouter()

	_ = r.Run(conf.HttpPort)
}
