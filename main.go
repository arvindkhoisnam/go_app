package main

import (
	"log"
	"net/http"
	"os"

	"github.com/arvindkhoisnam/go_app/models"
	"github.com/arvindkhoisnam/go_app/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)
type Todo struct{
	ID int `json:"id"`
	Completed bool `json:"completed"`
	Body string	`json:"body"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) AddTodo(c *fiber.Ctx) error {
	todo := &Todo{} 
	err := c.BodyParser(todo)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message":"Request failed."})
		return err
	}
	if todo.Body == ""{
		return c.Status(400).JSON(fiber.Map{"message":"Invalid input"})
		
	}  

	err = r.DB.Create(todo).Error 
	if err != nil {
		c.Status(400).JSON(fiber.Map{"message":"Error while creating todo"})
		return err
	}

	c.Status(200).JSON(fiber.Map{"message":"Todo created."})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App){
	api := app.Group("/api")
	api.Post("/todo",r.AddTodo)
	// api.Get("/",r.GetAllTodos)
	// api.Patch("/todo/:id",r.UpdateTodo)
	// api.Delete("/todo/:id",r.DeleteTodo)
}
func main(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error while loading .env")
	}
	PORT := os.Getenv("PORT")
	
	config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User: os.Getenv("DB_USER"),
		DBName: os.Getenv("DB_DB_NAME"),
		SSLMode: "",
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("could not load the database")
	}

	err = models.MigrateTodos(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB : db,
	}
	app := fiber.New()

	r.SetupRoutes(app)

	app.Get("/health",func(c *fiber.Ctx)error {
		return c.Status(200).JSON(fiber.Map{"message":"Healthy Server"})
	})

	// app.Post("/api/todo", func(c *fiber.Ctx) error {
	// 	todo := &Todo{}
	// 	err := c.BodyParser(todo)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if todo.Body == ""{
	// 		return c.Status(400).JSON(fiber.Map{"message":"Invalid Input"})
	// 	}

	// 	todo.ID = len(todos) + 1
	// 	todos = append(todos, *todo)
	// 	return c.Status(200).JSON(fiber.Map{"message":"Todo created!!","todo":todo})
	// })

	// app.Get("/api/all",func(c *fiber.Ctx) error{
	// 	return c.Status(200).JSON(todos)
	// })

	// app.Patch("/api/todo/:id",func (c *fiber.Ctx) error  {
	// 	id := c.Params("id")

	// 	for i, todo := range todos {
	// 		if fmt.Sprint(todo.ID) == id {
	// 			todos[i].Completed = true;
	// 			return c.Status(200).JSON(fiber.Map{"message":"Todo updated."})
	// 		}
	// 	}
	// 	return c.Status(400).JSON(fiber.Map{"message":"No todo found with the id"})
	// })

	// app.Delete("/api/todo/:id",func(c *fiber.Ctx) error {
	// 	id := c.Params("id")

	// 	for i, todo := range todos {
	// 		if fmt.Sprint(todo.ID) == id {
	// 			todos = append(todos[:i],todos[i+1:]... )
	// 			return c.Status(200).JSON(todos)
	// 		}
	// 	}
	// 	return c.Status(400).JSON(fiber.Map{"message":"No todo found with the id"})
	// })

	// err := app.Listen(PORT)
	// if err != nil {
	// 	log.Fatal(err)
	// } 
	log.Fatal(app.Listen(":"+PORT))
}