package persist

import (
	"github.com/go-redis/redis"
	"github.com/mneumi/reading-crawler/site/xcar/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var RDB *redis.Client
var DB *gorm.DB

func init() {
	var err error

	dsn := "root:123456@tcp(127.0.0.1:3306)/xcar?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(
		&model.CarDetail{},
		&model.CarModel{},
	)

	///

	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err = RDB.Ping().Result()

	if err != nil {
		panic(err)
	}
}
