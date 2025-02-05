package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arvindkhoisnam/go_app/middleware"
	"github.com/arvindkhoisnam/go_app/models"
	"github.com/arvindkhoisnam/go_app/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)
type Todo struct{
	ID int `json:"id"`
	Completed bool `json:"completed"`
	Body string	`json:"body"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Repository struct {
	DB *gorm.DB
}

const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

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

func (r *Repository) GetAllTodos(c *fiber.Ctx) error {
	todos := &[]models.Todos{}

	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	fmt.Println(userID)
	fmt.Println(role)

	err := r.DB.Find(todos).Error
	if err != nil {
		c.Status(400).JSON(fiber.Map{"message":"Unable to fetch todos"})
	}
	c.Status(200).JSON(fiber.Map{"data":todos})
	return nil
}

func (r *Repository) UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	todo := &models.Todos{}

	result := r.DB.Model(&models.Todos{}).Where("id = ?",id).Update("completed", "true")
	if result.Error != nil {
		c.Status(400).JSON(fiber.Map{"message":"Unable to update todo"})
	}
	 // Fetch the updated todo from the database
	err := r.DB.Where("id = ?", id).First(todo).Error
	if err != nil {
		// Return an error response if fetching the updated todo fails
		return c.Status(400).JSON(fiber.Map{"message": "Unable to fetch updated todo", "error": err.Error()})
		}
	c.Status(200).JSON(fiber.Map{"data":todo})
	return nil
}

func (r *Repository)DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {	
		fmt.Println(id)
		 return c.Status(400).JSON(fiber.Map{"message":"Invalid id"})
	}

    // Delete the todo from the database
    result := r.DB.Where("id = ?", id).Delete(&models.Todos{})
    if result.Error != nil {
        // Return an error response if the deletion fails
        return c.Status(400).JSON(fiber.Map{"message": "Unable to delete todo", "error": result.Error.Error()})
    }

    // Check if any rows were affected
    if result.RowsAffected == 0 {
        return c.Status(404).JSON(fiber.Map{"message": "Todo not found"})
    }
	c.Status(200).JSON(fiber.Map{"message":"Todo deleted successfully"})
	return nil
}

func (r *Repository)Signup(c *fiber.Ctx) error {
	user := &models.User{}
	err := c.BodyParser(user)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message":"Request failed."})
		return err
	}
	if user.Username == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{"message":"Invalid credentials."})
		
	}  
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password.",
		})
	}
	// Store hashed password
	user.Password = string(hashedPassword)
	err = r.DB.Create(user).Error
	if err != nil {
		c.Status(400).JSON(fiber.Map{"message":"Error while signinup."})
		return err
	}
	c.Status(200).JSON(fiber.Map{"message":"User successfully signed up."})
	return nil
}

func (r *Repository)Signin(c *fiber.Ctx) error {
	user := &User{}
	userModel := &models.User{}

	// Parse request body
	if err := c.BodyParser(user); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Request failed.",
		})
	}

	// Check if user exists in the database
	if err := r.DB.Where("username = ?", user.Username).First(userModel).Error; err != nil {
		return c.Status(403).JSON(fiber.Map{
			"message": "Invalid username or password.",
		})
	}

	// Compare provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(user.Password)); err != nil {
		return c.Status(403).JSON(fiber.Map{
			"message": "Invalid username or password.",
		})
	}


    // Generate JWT Token
	jwtSec := os.Getenv("JWT_SECRET")
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userModel.ID,
		"role": "admin",
    })
    tokenString, err := token.SignedString([]byte(jwtSec))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Could not generate token.",
			"error":err.Error(),
        })
    }

    // Set cookie in response
    c.Cookie(&fiber.Cookie{
        Name:     "auth_token",
        Value:    tokenString,
        HTTPOnly: true, // Prevents JavaScript access to cookie
        Secure:   true, // Use secure flag for HTTPS
        SameSite: "lax",
    })

	// If login is successful, return success response (you might want to generate a JWT token here)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Login successful.",
		"token":tokenString,
	})
}
func (r *Repository) SetupRoutes(app *fiber.App){
	api := app.Group("/api")
	api.Post("/todo",r.AddTodo)
	api.Get("/all",middleware.AuthMiddleware,r.GetAllTodos)
	api.Patch("/todo/:id",r.UpdateTodo)
	api.Delete("/todo/:id",r.DeleteTodo)
	api.Post("/signup",r.Signup)
	api.Post("/signin",r.Signin)
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

	err = models.MigrateDB(db)
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
	log.Fatal(app.Listen(":"+PORT))
}