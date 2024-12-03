package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/teelek/go-test/database"
	m "github.com/teelek/go-test/models"
	"github.com/teelek/go-test/routes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDatabase() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		"root",
		"",
		"127.0.0.1",
		"3306",
		"golang_test",
	)
	var err error
	database.DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connected!")
	database.DBConn.AutoMigrate(&m.Dog{})
	database.DBConn.AutoMigrate(&m.Companies{})
	database.DBConn.AutoMigrate(&m.Profiles{})
}

func main() {
	app := fiber.New()
	routes.InetRoutes(app)
	initDatabase()
	app.Listen(":3000")
}
