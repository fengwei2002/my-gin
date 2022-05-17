package main

import (
	"koo"
	"net/http"
)

func main() {
	r := koo.New()
	r.GET("/", func(c *koo.Context) {
		c.HTML(http.StatusOK, "<h1>Hello</h1>")
	})

	r.GET("/hello", func(c *koo.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *koo.Context) {
		// expect /hello?name=fengwei
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/static/*filepath", func(c *koo.Context) {
		// expect /static/background.css
		c.JSON(http.StatusOK, koo.H{"filepath": c.Param("filepath")})
	})

	r.Run("localhost:8080")
}