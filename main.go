package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type ResBody struct {
	Token string `json:"access_token"`
}

var clientId = "35267af0118570d03009"
var clientSecret = "34f8bda5539ca6ddd308655563e67a6729aac9ca"

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/login/oauth/access_token", func(c *gin.Context) {
		code := c.Query("code")

		u, err := url.Parse("https://github.com/login/oauth/access_token")

		v := url.Values{}

		v.Set("client_id", clientId)
		v.Set("client_secret", clientSecret)
		v.Set("code", code)

		u.RawQuery = v.Encode()

		req, _ := http.NewRequest("POST", u.String(), nil)

		req.Header.Add("Accept", "application/json")

		resp, err := (&http.Client{}).Do(req)

		if err != nil {
			fmt.Println("Fatal error ", err.Error())
		}

		if resp != nil {
			defer resp.Body.Close()
		}

		body, err := io.ReadAll(resp.Body)

		jsonObj := ResBody{}
		json.Unmarshal(body, &jsonObj)

		c.Redirect(http.StatusMovedPermanently, "http://127.0.0.1:5173?token="+jsonObj.Token)
	})
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
