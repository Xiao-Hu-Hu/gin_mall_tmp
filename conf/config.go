package conf

import (
	"gin_mall_tmp/dao"
	"gopkg.in/ini.v1"
	"strings"
)

var (
	AppModel string
	HttpPort string

	Db     string
	DbHost string
	DbPort string
	DbUser string

	DbPassword    string
	DbName        string
	RedisDb       string
	RedisAddr     string
	RedisPassword string
	RedisDbName   string

	ValidEmail string
	SmtpHost   string
	SmptEmail  string
	SmptPass   string

	Host        string
	ProductPath string
	AvatarPath  string
)

func Init() {
	// 本地读取环境变量
	//file, err := ini.Load("./conf/config.docker.ini")
	file, err := ini.Load("conf/config.ini")
	if err != nil {
		panic(err)
	}
	LoadServer(file)
	LoadMysql(file)
	LoadRedis(file)
	LoadEmail(file)
	LoadPhotoPath(file)

	//mysql 读(8) 主
	pathRead := strings.Join([]string{DbUser, ":", DbPassword, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=True&loc=Local"}, "")
	//mysql 写(2) 从  主从复制
	pathWrite := strings.Join([]string{DbUser, ":", DbPassword, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=True&loc=Local"}, "")

	dao.Database(pathRead, pathWrite)
}
func LoadServer(file *ini.File) {
	AppModel = file.Section("service").Key("AppModel").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func LoadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassword = file.Section("mysql").Key("DbPassword").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func LoadRedis(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPassword = file.Section("redis").Key("RedisPassword").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}

func LoadEmail(file *ini.File) {
	ValidEmail = file.Section("email").Key("ValidEmail").String()
	SmtpHost = file.Section("email").Key("SmtpHost").String()
	SmptEmail = file.Section("email").Key("SmptEmail").String()
	SmptPass = file.Section("email").Key("SmptPass").String()
}

func LoadPhotoPath(file *ini.File) {
	ProductPath = file.Section("path").Key("ProductPath").String()
	AvatarPath = file.Section("path").Key("AvatarPath").String()
	Host = file.Section("path").Key("Host").String()
}
