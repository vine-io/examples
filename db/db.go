package main

import (
	"fmt"
	"log"
	"time"

	"github.com/vine-io/plugins/dao/mysql"
	"github.com/vine-io/vine/lib/dao"
	"github.com/vine-io/vine/lib/dao/logger"
)

type User struct {
	Id   int32  `dao:"column:id;autoIncrement;primaryKey"`
	Name string `dao:"column:name"`
	Age  int32  `dao:"column:age"`
}

func main() {
	dns := `root:123456@tcp(192.168.3.111:3306)/vine?charset=utf8&parseTime=True&loc=Local`
	dialect := mysql.NewDialect(dao.DSN(dns), dao.Logger(logger.New(logger.Options{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Info,
	})))
	if err := dialect.Init(); err != nil {
		log.Fatalln(err)
	}

	db := dialect.NewTx()

	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalln(err)
	}

	u := &User{Name: "Mimi", Age: 11}
	if err := db.Create(u).Error; err != nil {
		log.Fatalln(err)
	}

	u1 := &User{}
	if err := db.Find(&u1, "name = ?", "Mimi").Error; err != nil {
		log.Fatalln(err)
	}

	fmt.Println(u1)

	if err := db.Where("name = ?", "Mimi").First(&u1).Error; err != nil {
		log.Fatalln(err)
	}

	fmt.Println(u1)

	u1.Name = "Mimi_update"
	if err := db.Updates(u1).Error; err != nil {
		log.Fatalln(err)
	}


	if err := db.Delete(&User{}, "id = ?", 1).Error; err != nil {
		log.Fatalln(err)
	}
}
