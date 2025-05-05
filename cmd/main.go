//

package main // Declares the package name - 'main' is special in Go as it defines an executable program, not a library

import (
	"fmt"
	"html/template" // Package for generating HTML output with templates (similar to templating engines like Handlebars in JS)
	"io"            // Provides basic interfaces for I/O primitives (readers, writers, etc.)

	"github.com/labstack/echo/v4"            // Echo is a high-performance web framework for Go (like Express.js in the Node.js world)
	"github.com/labstack/echo/v4/middleware" // Echo's middleware components (similar to middleware in Express)
)

// Templates struct contains a pointer to template.Template
// In Go, structs are somewhat similar to JS objects, but they're strongly typed
// This struct will store our compiled templates
type Templates struct {
	templates *template.Template // The asterisk (*) denotes a pointer, which stores a memory address rather than a value directly
}

// Render is a method defined on the Templates struct
// In Go, methods are functions attached to a specific type (similar to methods in JS classes)
// This method implements Echo's Renderer interface, which requires a Render method with this signature
// The (t *Templates) part is called a "receiver" and defines which type this method belongs to
func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// ExecuteTemplate applies a named template to the specified data and writes the output to w
	// This is similar to template engines in JS where you render a template with data
	return t.templates.ExecuteTemplate(w, name, data)
}

// newTemplate is a constructor function that creates and returns a new Templates instance
// In Go, it's a convention to use "new" prefix for constructor functions (whereas JS uses "new" keyword)
func newTemplate() *Templates {
	return &Templates{ // The ampersand (&) creates a pointer to the struct we're creating
		// template.Must is a helper that panics if the template has an error (like throwing in JS)
		// ParseGlob loads all templates matching the pattern (similar to requiring multiple template files in JS)
		templates: template.Must(template.ParseGlob("views/*.tmpl")),
	}
}

// Count is a struct that holds a single integer field
// This is similar to creating a simple object in JS like { count: 0 }
// In Go, struct field names are capitalized if they need to be accessible outside the package
type Count struct {
	Count int // An integer field to store our counter value
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
			newContact("John Doe", "asd"),
			newContact("Jane Doe", "jane@gmail.com"),
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

// main function is the entry point for Go programs (similar to the main JS file that gets executed)
func main() {
	e := echo.New()            // Create a new instance of Echo (similar to creating an Express app in JS)
	e.Use(middleware.Logger()) // Add logging middleware (similar to app.use(logger()) in Express)

	// Initialize our counter - in Go, variables are strongly typed
	// Unlike JS where variables are dynamically typed
	count := Count{Count: 0} // Create a new Count struct and initialize the Count field to 0
	page := newPage()

	// Set the renderer for our Echo instance
	// In JS, you'd similarly configure your Express app to use a template engine
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error { // Anonymous function as route handler (like a callback in JS)
		count.Count++
		// This is similar to res.render('index', {count: count}) in Express
		return c.Render(200, "test1", count)
	})
	e.GET("/x", func(c echo.Context) error {
		fmt.Println("X route called")
		return c.String(200, "X route works!")
	})
	e.POST("/count", func(c echo.Context) error {
		count.Count++
		return c.Render(200, "count", count)
	})

	e.GET("/contacts", func(c echo.Context) error { // Anonymous function as route handler (like a callback in JS)
		// This is similar to res.render('index', {count: count}) in Express
		return c.Render(200, "index", page)
	})
	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.Data.hasEmail(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "Email already exists"
			fmt.Printf("Error thing: %+v\n", page)
			return c.Render(422, "form", formData)
		}

		page.Data.Contacts = append(page.Data.Contacts, newContact(name, email))
		fmt.Printf("Data passed to template: %+v\n", page)
		return c.Render(200, "display", page)
	})

	// e.Any("/*", func(c echo.Context) error {
	// 	fmt.Println("Catch-all route called for:", c.Request().URL.Path)
	// 	return c.String(200, "Caught by catch-all")
	// })

	// Start the server on port 42069
	// Note that in Go, unlike JS, errors are typically handled immediately rather than with promises/callbacks
	// e.Logger.Fatal will log the error and exit the program if e.Start returns an error
	// In JS you might do something like app.listen(42069).catch(err => console.error(err))
	e.Logger.Fatal(e.Start(":42069"))
}
