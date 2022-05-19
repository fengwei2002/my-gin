package main

import (
	"fmt"
	"html/template"
	"koo"
	"net/http"
	"time"
)

type Student struct {
	Name string
	Age int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Year(), t.Month(), t.Day()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := koo.New()
	r.Use(koo.Logger()) // global middleware
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	}) // 添加自定义模板渲染函数
	
	r.LoadHTMLGlob("templates/*") // 加载 templates 下的静态文件 存储在 engine 中 

	r.Static("/assets", "./static")
	stu1 := &Student{Name: "fengwei", Age: 20}
	stu2 := &Student{Name: "oldFengwei", Age: 22}

	r.GET("/", func(c *koo.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	}) // 根据名字渲染 html 文件 

	r.GET("/students", func(c *koo.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", koo.H{
			"title":  "koo",
			"stuArr": [2]*Student{stu1, stu2},
		})
	}) // 渲染 arr 文件到 html，使用 json 格式

	r.GET("/date", func(c *koo.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", koo.H{
			"title": "koo",
			"now":   time.Date(2022, 5, 19, 0, 0, 0, 0, time.UTC),
		})
	}) // 使用 json 将自定义函数的结果作为 html 渲染


	r.Run("localhost:8080")
}