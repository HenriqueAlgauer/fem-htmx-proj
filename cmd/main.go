package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Contact struct {
	Name  string
	Email string
}

func newContact(name string, email string) Contact {
	return Contact{
		Name:  name,
		Email: email,
	}
}

type Contacts = []Contact

type Data struct {
	Contacts Contacts
}

func (d *Data) hasEmail(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

func newData() Data {
	return Data{
		Contacts: []Contact{
			newContact("Jhon", "jd@gmail.com"),
			newContact("Clara", "cd@gmail.com"),
		},
	}
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	page := newPage()
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.Data.hasEmail(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "Email alredy exists"

			return c.Render(422, "form", formData)
		}
		contact := newContact(name, email)
		page.Data.Contacts = append(page.Data.Contacts, contact)

		c.Render(200, "form", newFormData())
		return c.Render(200, "oob-contact", contact)
	})

	e.Logger.Fatal(e.Start((":42069")))
}

// package main

// import (
// 	"html/template"
// 	"io"
// 	"net/http"
// 	"strconv"

// 	"github.com/labstack/echo/v4"
// 	"github.com/labstack/echo/v4/middleware"
// )

// type Templates struct {
// 	templates *template.Template
// }

// func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
// 	return t.templates.ExecuteTemplate(w, name, data)
// }

// func NewTemplates() *Templates {
// 	return &Templates{
// 		templates: template.Must(template.ParseGlob("views/*.html")),
// 	}
// }

// type Block struct {
//     Id int
// }

// type Blocks struct {
//     Start int
//     Next int
//     More bool
//     Blocks []Block
// }

// func main() {
// 	e := echo.New()
//     e.Renderer = NewTemplates()
//     e.Use(middleware.Logger())

//     e.GET("/blocks", func(c echo.Context) error {
//         startStr := c.QueryParam("start")
//         start, err := strconv.Atoi(startStr)
//         if err != nil {
//             start = 0
//         }

//         blocks := []Block{}
//         for i := start; i < start + 10; i++ {
//             blocks = append(blocks, Block{Id: i})
//         }

//         template := "blocks"
//         if start == 0 {
//             template = "blocks-index"
//         }
//         return c.Render(http.StatusOK, template, Blocks{
//             Start: start,
//             Next: start + 10,
//             More: start + 10 < 100,
//             Blocks: blocks,
//         });
//     });

//     e.Logger.Fatal(e.Start(":42069"))
// }
