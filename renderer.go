package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	apiroutes "web-dashboard/api-routes"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

/*
	type Renderer interface {
		Render(io.Writer, string, interface{}, Context) error
	}

	define interface Template with method Render

*/

func (t *Template) Render(wr io.Writer, name string, data interface{}, c echo.Context) error {
	err := t.templates.ExecuteTemplate(wr, name, data)

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

func (t *Template) Populate(name string) interface{} {
	// maps template names to api routes
	// return a callback?

	return apiroutes.GetEnumerateWindows()
}

func mustImportTemplates() *template.Template {
	templ, err := template.ParseGlob("templates/*") // Parses all .html files in the templates directory
	if err != nil {
		panic(err)
	}
	return templ
}

// TRASH
// populate := map[string]func() []interface{}{
// 	"open-tabs.static.html": func() []interface{} {
// 		a := pwsh.RunPwshCmd("./waterfoxTabs.ps1")
// 		return a
// 	},
// }

// The logic for retrieving data should be defined in this function,
// and retrievable from the echo.Context.
//
// e.g. tries to get populate function from a map[string]func() []interface{}
// interfaces can be nil...
// so return nil if no function for the template's populate fn request.path
// any filtering logic, etc can be in parameters.
// any other logic should be in middleware.
// good.
