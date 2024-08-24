package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
	"ts/internal/models"
)

//	type DBConnection interface {
//		transport.Manager
//		UsersControllerInterface
//	}
type UsersControllerInterface interface {
	CreateUser(user interface{}) error
	FindFirst(user interface{}, scond string, cond interface{}) error
}

type UHandler struct {
	Repo   UsersControllerInterface
	Logger *zerolog.Logger
}

func NewHandler(repo UsersControllerInterface) *UHandler {
	return &UHandler{Repo: repo}
}

func SignUp(db *UHandler) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Получение логина и пароля из тела запроса
		var body struct {
			Login    string
			Password string
		}
		err := c.Bind(&body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read body",
			})
			return
		}
		// Хеширование пароля
		hashPassword, errH := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

		if errH != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to hash password",
			})
			return
		}

		// Создание пользователя
		user := models.User{Login: body.Login, Password: string(hashPassword)}

		result := db.Repo.CreateUser(&user)

		if result != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create user",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"INFO": "User created successfully"})
	}
}

func Login(db *UHandler) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Получение логина и пароля из тела запроса
		var body struct {
			Login    string
			Password string
		}
		err := c.Bind(&body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read body",
			})
			return
		}
		// Поиск запрошенного пользователя
		var user models.User

		result := db.Repo.FindFirst(&user, "login = ?", body.Login)

		if result != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create user",
			})
			return
		}
		if user.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid login or password",
			})
			return
		}
		// Сравенение пароля с базой данных
		errP := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
		if errP != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid login or password",
			})
			return
		}
		// Генерация jwt токена
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		tokenString, errT := token.SignedString([]byte(os.Getenv("SECRET")))

		if errT != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create token",
			})
			return
		}

		// Отправить токен
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"INFO": user.Login + " logged in",
		})
	}

}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"INFO": "In validation",
	})
}

// интерфейст для работы с юзерами,
//
