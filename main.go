package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)
type Todo struct{
	ID int `json:"id"`
	Completed bool `json:"completed"`
	Body string	`json:"body"`
}

func main(){

	fmt.Println("Hello World.")
	app := fiber.New()

	todos := []Todo{}
	app.Get("/health",func(c *fiber.Ctx)error {
		return c.Status(200).JSON(fiber.Map{"message":"Healthy Server"})
	})

	app.Post("/api/todo", func(c *fiber.Ctx) error {
		todo := &Todo{}
		err := c.BodyParser(todo)
		if err != nil {
			return err
		}
		if todo.Body == ""{
			return c.Status(400).JSON(fiber.Map{"message":"Invalid Input"})
		}
		todo.ID = len(todos) + 1
		todos = append(todos, *todo)
		return c.Status(200).JSON(fiber.Map{"message":"Todo created!!","todo":todo})
	})

	app.Get("/api/all",func(c *fiber.Ctx) error{
		return c.Status(200).JSON(todos)
	})

	app.Patch("/api/todo/:id",func (c *fiber.Ctx) error  {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true;
				return c.Status(200).JSON(fiber.Map{"message":"Todo updated."})
			}
		}
		return c.Status(400).JSON(fiber.Map{"message":"No todo found with the id"})
	})

	app.Delete("/api/todo/:id",func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i],todos[i+1:]... )
				return c.Status(200).JSON(todos)
			}
		}
		return c.Status(400).JSON(fiber.Map{"message":"No todo found with the id"})
	})
	err := app.Listen(":9000")
	if err != nil {
		log.Fatal(err)
	} 
}