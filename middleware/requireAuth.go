package middleware

import (
	"JWTAUTH/initializers"
	"JWTAUTH/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Requireuth(c *gin.Context) {
	//Get the cookie from request

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed 1 ",
		})
		//c.AbortWithStatus(http.StatusUnauthorized)
	}
	//decode and validate it
	// Parse takes the token string and a function for looking up the key. The latter is especially
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println(ok, token, token.Valid, claims)
	if ok && token.Valid {
		//check the expiry of cookie
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		//Find the user withtoken sub
		var user models.UserOpen
		initializers.DB.First(&user, claims["sub"]) //checking user id from cookie as we had
		// passed user.id in cookie While login
		if user.ID == 0 {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed 2 ",
			})
			//c.AbortWithStatus(http.StatusUnauthorized)
		}
		//attach to request
		c.Set("user", user)

		//continue
		c.Next()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed 3",
		})
		//c.AbortWithStatus(http.StatusUnauthorized)

	}

}
