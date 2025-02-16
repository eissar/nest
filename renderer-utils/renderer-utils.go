package rendererutils

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	//apiroutes "web-dashboard/api-routes"

	"github.com/labstack/echo/v4"
)

// NOTE: Template rules:
// 1. ending in .html:  static template.
// 2. ending in .templ:  dynamic template
// 3. prefix ws, sse, ending in .templ:  template which retrieves data...
// using websockets or server-side events respectively.
// 4. no prefix, ending in .templ: template which retrieves data...
// dynamically using htmx.

type Template struct {
	Templates *template.Template
}

/*
	type Renderer interface {
		Render(io.Writer, string, interface{}, Context) error
	}

	define interface Template with method Render

*/

func (t *Template) Render(wr io.Writer, name string, data interface{}, c echo.Context) error {
	err := t.Templates.ExecuteTemplate(wr, name, data)

	// some error handling.
	if err != nil {
		// you can return here! error code 402 or something
		// TODO: use the c.Error with middleware
		msg := fmt.Sprintf("Error rendering template %s @Render error: %s", name, err.Error())
		fmt.Println("[ERROR]", msg)

		e := fmt.Sprintf("<p>%s</p>", msg)
		return c.String(http.StatusUnprocessableEntity, e)
	}
	return nil
}

/*
func (t *Template) Populate(name string) interface{} {
	// maps template names to api routes
	// return a callback?

	return apiroutes.GetEnumerateWindows()
}
*/

var Templates *template.Template

func MustImportTemplates() *template.Template {
	templ, err := template.ParseGlob("templates/*.templ") // Parses all files in the templates directory
	templ.ParseGlob("templates/*.html")                   // parse static templates
	if err != nil {
		panic(err)
	}
	return templ
}

type DynamicTemplatePopulateFunc func(c echo.Context, templateName string) interface{}

func InitRenderer() {
	Templates = MustImportTemplates()
}
