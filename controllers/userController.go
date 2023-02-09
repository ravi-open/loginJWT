package controllers

import (
	"JWTAUTH/initializers"
	"JWTAUTH/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	//Get the email/password
	var body struct {
		Email    string `json:"email"`
		Password string `json:"pass"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//Hast the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to hash Password",
		})
		return
	}
	//Create the user
	user := models.UserOpen{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)
	//will store user in DB with hashed password
	if result.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to Create user",
		})
		return
	}

	//Respond
	c.JSON(http.StatusOK, gin.H{

		"message": "Successfully Signed up",
	})
}
func Login(c *gin.Context) {
	//Get email aand password of body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"pass"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadGateway, gin.H{

			"error": "Failed to read body",
		})
		return
	}
	//Look up requsted user
	var user models.UserOpen
	initializers.DB.First(&user, "email=?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or Password",
		})
		return
	}

	//compare sent in passwordith saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or Password",
		})
		return
	}
	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		//exp is expirt time we are adding current time+#0 days, if we don't
		//add extra time, it will expire immediately
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET ")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to Create Token ",
		})
		return
	}
	//end it back
	//SENDING TOKEN AS WELL BUT THEN WE WILL STORE IN COOKIE COZ ITS GOD
	// c.JSON(http.StatusOK, gin.H{
	// 	"token": tokenString,
	// })
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{

		"message": "Successfully Login",
	})

}
func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"user is ": user,
	})
}
