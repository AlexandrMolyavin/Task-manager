package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
	"ts/internal/controllers"
	"ts/internal/models"
)

func RequireAuth(db *controllers.UHandler) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Получение cookie из ответа
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			db.Logger.Error().Msg("Failed to get cookie")
			fmt.Fprintf(c.Writer, "Failed to get cookie :'%v'", err)
			return
		}

		//
		token, errP := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("SECRET")), nil
		})
		if errP != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			//проверка времени действия токена
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			//поиск пользователя
			var user models.User

			result := db.Repo.FindFirst(&user, "", claims["sub"])
			if result != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Failed to find user",
				})
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			if user.ID == 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			//прикрепление к запросу
			c.Set("user", user)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

	}
}
