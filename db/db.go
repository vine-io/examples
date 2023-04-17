package main

import (
	"context"
	"examples/db/proto"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	dns := `root:@tcp(127.0.0.1:8999)/vine?charset=utf8&parseTime=True&loc=Local`
	gormDB, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	store := userv1.NewUserStorage(gormDB, &userv1.User{})
	if err != nil {
		log.Fatalln(err)
	}

	// 注册 Schema, 会在数据库中创建对应的表
	if err := store.AutoMigrate(); err != nil {
		log.Fatal(err)
	}

	ctx := context.TODO()

	user := &userv1.User{
		Id:   1,
		Name: "Mimi",
		Age:  18,
	}

	fmt.Println("Create ==============>")
	s := userv1.NewUserStorage(gormDB, user)
	out, err := s.XXCreate(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)

	fmt.Println("Updates ==============>")
	s = userv1.NewUserStorage(gormDB, &userv1.User{Id: 1, Name: "Mimi_rename"})
	out, err = s.XXUpdates(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)

	fmt.Println("FindAll ==============>")
	s = userv1.NewUserStorage(gormDB, &userv1.User{Age: 18})
	outs, err := s.XXFindAll(ctx)
	fmt.Println(outs)

	fmt.Println("FindOne ==============>")
	s = userv1.NewUserStorage(gormDB, &userv1.User{Name: "Mimi_rename"})
	out, err = s.XXFindOne(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)

	fmt.Println("SoftDelete ==============>")
	s = userv1.NewUserStorage(gormDB, &userv1.User{Id: 1})
	err = s.XXDelete(ctx, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Delete ==============>")
	s = userv1.NewUserStorage(gormDB, &userv1.User{Id: 1})
	err = s.XXDelete(ctx, false)
	if err != nil {
		log.Fatal(err)
	}
}
