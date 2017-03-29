package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ardanlabs/cobalt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func app(db *sqlx.DB) *cobalt.Cobalt {
	c := cobalt.New(JSONEncoder{})

	// Recompile templates for each request.
	if os.Getenv("ENV") != "prod" {
		c.Templates.Development = true
	}

	// Define our custom template functions
	c.Templates.Funcs = template.FuncMap{
		"markdown": func(s string) template.HTML {
			unsafe := blackfriday.MarkdownCommon([]byte(s))
			html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
			return template.HTML(html)
		},
		"time": func(t time.Time) string {
			return t.Format(time.RubyDate)
		},
	}

	j := &JobHandlers{db: db}

	// Wire up job handlers to routes
	c.Get("/", j.Index)
	c.Get("/new", j.New)
	c.Post("/jobs", j.Create)
	c.Get("/jobs/:id", j.Show)

	// Serve static assets from the "public" directory
	c.ServeFiles("/public/*filepath", http.Dir("public"))

	return c
}

func main() {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	app(db).Run(":" + os.Getenv("PORT"))
}
