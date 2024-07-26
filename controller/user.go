package controller

import (
	"database1/database"
	"database1/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

type UserSerializer struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	NAME     string `json:"name"`
	EMAIL    string `json:"email"`
	PASSWORD string `json:"password"`
}

func CreateResponseUser(userModel model.User) UserSerializer {
	return UserSerializer{ID: userModel.ID, NAME: userModel.Name, EMAIL: userModel.Email, PASSWORD: userModel.Password}
}

type loginUser struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

func createToken(user model.User) (string, error) {
	key := os.Getenv("SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Token{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags user
// @Accept json
// @Produce json
// @Param user body UserSerializer true "User Data"
// @Router /user/register [post]
func Register(c *fiber.Ctx) error {
	var user model.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	validate := validator.New()

	err := validate.Struct(user)
	if err != nil {
		// Validation failed, handle the error
		errors := err.(validator.ValidationErrors)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": http.StatusBadRequest, "message": "validation Failed", "error": errors.Error()})
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	user.Password = string(hashedPassword)
	database.Database.DB.Create(&user)
	responseUser := CreateResponseUser(user)

	return c.Status(http.StatusOK).JSON(fiber.Map{"status": http.StatusOK, "message": "User Registered successfully", "data": responseUser})
}

// Login godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags user
// @Accept json
// @Produce json
// @Param user body UserSerializer true "User Data"
// @Router /user/login [post]
func Login(c *fiber.Ctx) error {
	var loginUser loginUser
	var user model.User
	if err := c.BodyParser(&loginUser); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	validate := validator.New()

	e := validate.Struct(user)
	if e != nil {
		// Validation failed, handle the error
		errors := e.(validator.ValidationErrors)
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"status": http.StatusBadGateway, "message": "validation Failed", "error": errors.Error()})
	}

	userData := database.Database.DB.Where("Email = ?", loginUser.Email).Find(&user)
	if userData == nil {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"status": http.StatusBadGateway, "message": "Email is not match"})
	}
	er := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if er != nil {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"status": http.StatusBadGateway, "message": "Password is not match"})
	}

	var token, err = createToken(user)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"status": http.StatusOK, "message": "User Login successfully", "token": token})
}
