package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/FuturaInsTech/GoExcel/initializers"
	"github.com/FuturaInsTech/GoExcel/models"
	"github.com/FuturaInsTech/GoExcel/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequiredAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		//fmt.Println(claims["sub"], claims["aud"], claims["iat"], claims["exp"])

		//check expiry

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		var user models.User
		initializers.DB.First(&user, "phone = ?", claims["sub"])
		if user.Id == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		var message types.User
		message.Name = user.Name
		message.Email = user.Email
		message.Phone = user.Phone
		message.Id = user.Id
		message.UserGroupId = user.UserGroupID
		message.LanguageId = user.LanguageID
		message.CompanyId = user.CompanyID

		c.Set("user", message)

		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

func CheckRefreshToken(c *gin.Context) {

	tokenString, err := c.Cookie("RefreshToken")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET_REFRESH_TOKEN")), nil
	})

	if err != nil {

		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		var user models.User
		initializers.DB.First(&user, "id = ?", claims["sub"])
		if user.Id == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("user", user.Phone)

		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
