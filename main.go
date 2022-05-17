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
	r.POST("/login", func(c *koo.Context) {
		c.JSON(http.StatusOK, koo.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run("localhost:8080")
}