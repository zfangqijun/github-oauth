package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

var (
	AppCommitSHA string
)

func main() {
	fmt.Printf("AppCommitSHA: %s\n", AppCommitSHA)

	r := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*.7z7z.cc"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		AllowWildcard:    true,
	}

	r.Use(cors.New(corsConfig))

	r.POST("/login/oauth/access_token", func(c *gin.Context) {
		code := c.Query("code")

		result, err := OAuthByCode(code)

		if err != nil {
			fmt.Println(err)
			c.Data(http.StatusBadRequest, "application/json; charset=utf-8", result)
		}

		token := gjson.Get(string(result), "access_token").String()

		if token != "" {
			scope := gjson.Get(string(result), "token_type").String()
			c.SetCookie("_gho", scope+" "+token, 60*60*24*7, "", "wado.local", false, true)
		}

		c.Data(http.StatusOK, "application/json; charset=utf-8", result)
	})

	// github api proxy

	target, _ := url.Parse("https://api.github.com")

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ModifyResponse = func(resp *http.Response) error {
		// 删除来自 GitHub 原响应报文的 Allow-Origin
		// 因为 Allow-Origin头 不允许有多个
		resp.Header.Del("Access-Control-Allow-Origin")
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		return
	}

	originalDirector := proxy.Director

	r.Any("/api/*path", func(c *gin.Context) {
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = target.Host
			req.URL.Path = strings.Replace(req.URL.Path, "/api", "", 1)
			req.Header.Set("Authorization", getTokenFromCookie(c))
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.Run(":5174")

}

func getTokenFromCookie(c *gin.Context) string {
	token, _ := c.Cookie("_gho")
	fmt.Println((token))
	return token
}
