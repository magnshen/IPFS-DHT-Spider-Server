package main

import (
	"./controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"net/http"
	"os"
)
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  string
	)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Error()
	}
	c.Render(code, "error.html",map[string]interface{}{
		"message": msg,
	})
}

var ff *os.File

func main() {

	// Echo instance
	e := echo.New()
	// Middleware
	//e.Use(middleware.Recover())
	//e.Use(middleware.Logger())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))


	//static pages
	e.Static("/static", "static")

	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	//errot_pages
	e.HTTPErrorHandler = customHTTPErrorHandler
	//Routers
	e.GET("/",controllers.Index)
	e.GET("/api/getNews",controllers.GetNews)
	e.GET("/ipfsObject/:hash",controllers.GetIpfsObject)
	// Start server
	e.Logger.Fatal(e.Start(":80"))


}
