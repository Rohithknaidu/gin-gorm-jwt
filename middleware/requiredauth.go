package middleware

import (
	"fmt"
	"golang/jwt/initializers"
	"golang/jwt/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequiredAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		log.Println("error during collecting cookie",err)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		log.Println("checking token err", err)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil || token == nil {
		log.Println("error during parsing token ",err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)

		}
		var user models.User
		initializers.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			log.Println("error when there is no user left")
		}
		c.Set("user", user)

		c.Next()

	} else {
		log.Println("error during final stage")

	}

}
