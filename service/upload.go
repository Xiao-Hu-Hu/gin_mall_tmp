package service

import (
	"gin_mall_tmp/conf"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strconv"
)

func UploadAvatarToLocalStatic(file multipart.File, userId uint, userName string) (filepath string, err error) {
	bId := strconv.Itoa(int(userId))
	// 构架存储路径
	basePath := "." + conf.AvatarPath + "user" + bId + "/"
	if !DirExistOrNot(basePath) {
		CreatDir(basePath)
	}
	// 完整路径
	avatarPath := basePath + userName + ".jpg"
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(avatarPath, content, 0666)
	if err != nil {
		return "", err
	}
	return "user" + bId + "/" + userName + ".jpg", nil
}

func UploadProductToLocalStatic(file multipart.File, userId uint, productName string) (filepath string, err error) {
	bId := strconv.Itoa(int(userId))
	// 构架存储路径 ./static/imgs/product/boss{userId}/
	basePath := "." + conf.ProductPath + "boss" + bId + "/"

	// 检查目录是否存在，不存在则创建
	if !DirExistOrNot(basePath) {
		CreatDir(basePath)
	}
	// 完整路径
	productPath := basePath + productName + ".jpg"

	// 读取文件内容
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	// 写入文件到磁盘
	err = ioutil.WriteFile(productPath, content, 0666)
	if err != nil {
		return "", err
	}
	// 返回路径
	return "boss" + bId + "/" + productName + ".jpg", nil
}

// 判断文件夹路径是否存在
func DirExistOrNot(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// CreatDir 创建文件夹
func CreatDir(dirName string) bool {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return false
	}
	return true
}
