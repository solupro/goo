package main

import (
	"fmt"
	"goo"
	"html/template"
	"net/http"
	"time"
)

func main() {
	r := goo.Default()
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")
	r.GET("/index", func(c *goo.Context) {
		c.HTML(http.StatusOK, "index.html", goo.H{
			"name": c.Query("name"),
			"now":  time.Now(),
		})
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

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
