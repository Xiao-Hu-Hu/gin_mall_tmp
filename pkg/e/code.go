package e

const (
	Success       = 200
	Error         = 500
	InvalidParams = 400

	// user模块错误
	ErrorExistUser             = 3001
	ErrorFailEncryption        = 3002
	ErrorExistUserNotFound     = 3003
	ErrorNotCompare            = 3004
	ErrorAuthToken             = 3005
	ErrorAuthCheckTokenTimeOut = 3006
	ErrorUploadFail            = 3007
	ErrorSendEmail             = 3008

	// product模块
	ErrorProductImgUpload = 40001

	// 收藏夹错误
	ErrorFavoriteExist = 50001
)
