package main

import (
	"goo"
	"net/http"
)

func main() {
	r := goo.Default()
	r.GET("/index", func(c *goo.Context) {
		c.HTML(http.StatusOK, "<h1>Hello World</h1>")
	})

	r.GET("/crash", func(c *goo.Context) {
		arr := []int{1, 2, 3}
		c.String(http.StatusOK, string(arr[4]))
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/foo", func(c *goo.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/foo/:name", func(c *goo.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.GET("/assets/*filepath", func(c *goo.Context) {
			c.JSON(http.StatusOK, goo.H{"filepath": c.Param("filepath")})
		})

		v2.POST("/login", func(c *goo.Context) {
			c.JSON(http.StatusOK, goo.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.Run(":8090")
}
