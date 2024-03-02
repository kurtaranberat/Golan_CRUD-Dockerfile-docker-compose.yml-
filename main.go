package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5434
	user     = "postgres"
	password = "Password!" // Değişiklik burada
	dbname   = "User"      // Değişiklik burada
)

// User modelini tanımla
type User struct {
	gorm.Model
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(postgres.Open(getDBConnectionString()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{})
	app := fiber.New()

	app.Get("/users", getUsers)
	app.Get("/users/:id", getUser)
	app.Post("/users", createUser)
	app.Put("/users/:id", updateUser)
	app.Delete("/users/:id", deleteUser)

	log.Fatal(app.Listen(":3000"))
}

func getDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func getUsers(c *fiber.Ctx) error {
	var users []User
	db.Find(&users)
	return c.JSON(users)
}

func getUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	var user User
	db.First(&user, id)
	return c.JSON(user)
}

func createUser(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	db.Create(&user)
	return c.Status(fiber.StatusCreated).JSON(user)
}

func updateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	var user User
	db.First(&user, id)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	db.Save(&user)
	return c.Status(fiber.StatusOK).JSON(user)
}

func deleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	db.Delete(&User{}, id)
	return c.SendString("User deleted successfully")
}
