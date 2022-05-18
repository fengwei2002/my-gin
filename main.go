package main

import (
	"koo"
	"log"
	"net/http"
	"time"
)

func main() {
	r := koo.New()
	r.Use(koo.Logger()) // global middleware
	r.GET("/index", func(c *koo.Context) {
		c.HTML(http.StatusOK, "<h1>Hello</h1>")
	})
	v1 := r.Group("/v1") // *RouterGroup
	{
		v1.GET("/", func(c *koo.Context) {
			c.HTML(http.StatusOK, "<h1>KOO From V1</h1>")
		})

		v1.GET("/hello", func(c *koo.Context) {
			// expect /hello?name=fengwei
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // use v2 group middleware
	{
		v2.GET("/hello/:name", func(c *koo.Context) {
			// expect /hello/fengwei
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *koo.Context) {
			c.JSON(http.StatusOK, koo.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.Run("localhost:8080")
}


func onlyForV2() koo.HandlerFunc {
	return func(c *koo.Context) {
		// start time
		t := time.Now()
		// end
		log.Printf("> [%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}