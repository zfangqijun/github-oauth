package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/login/oauth/access_token", func(c *gin.Context) {
		code := c.Query("code")

		result, err := OAuthByCode(code)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(result.Scope)

		c.Redirect(http.StatusMovedPermanently, "http://127.0.0.1:5173?token="+result.AccessToken)
	})
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
