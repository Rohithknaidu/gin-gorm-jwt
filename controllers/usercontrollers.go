package controllers

import (
	"golang/jwt/initializers"
	"golang/jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var data struct {
		Email    string
		Password string
	}
	if c.Bind(&data) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed during the hashing the password",
		})
	}
	user := models.User{Email: data.Email, Password: string(hash)}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user data created",
	})
}
func Login(c *gin.Context) {
	var data struct {
		Email    string
		Password string
	}
	if c.Bind(&data) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
	}
	var user models.User
	initializers.DB.First(&user, "email= ?", data.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid eamil id or password",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid  password",
		})
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return 
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Autharizatiion",tokenString,3600*24*30,"","",true,true)
	c.JSON(http.StatusOK,gin.H{})

}

func Validate(c *gin.Context){
	user,_:=c.Get("user")
	c.JSON(http.StatusOK,gin.H{
		"message":user,
	})
}
