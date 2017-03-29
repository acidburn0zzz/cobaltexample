package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ardanlabs/cobalt"
	"github.com/jmoiron/sqlx"
)

// Job is a job posting provided by a company.
type Job struct {
	ID          int       // Provided by the db
	Title       string    // Name of the job
	Company     string    // Name of the company
	Description string    // Markdown text description of the job
	Contact     string    // Name of who posted it
	Created     time.Time // When it was created
	Updated     time.Time // When it was last updated
}

// JobHandlers holds all of the Job related handlers.
type JobHandlers struct {
	db *sqlx.DB
}

// Index shows a list of jobs
func (j *JobHandlers) Index(ctx *cobalt.Context) {
	var jobs []Job
	if err := j.db.Select(&jobs, "SELECT * FROM jobs"); err != nil {
		log.Print(err)
		ctx.ServeStatus(http.StatusInternalServerError)
		return
	}

	ctx.ServeHTML("jobs/index", jobs)
}

// New shows a form for creating a job posting
func (j *JobHandlers) New(ctx *cobalt.Context) {
	ctx.ServeHTML("jobs/new", nil)
}

// Create accepts a POST and creates a new Job
func (j *JobHandlers) Create(ctx *cobalt.Context) {
	now := time.Now()
	job := Job{
		Title:       ctx.Request.FormValue("title"),
		Company:     ctx.Request.FormValue("company"),
		Description: ctx.Request.FormValue("description"),
		Contact:     ctx.Request.FormValue("contact"),
		Created:     now,
		Updated:     now,
	}

	stmt, err := j.db.PrepareNamed(`INSERT INTO jobs (title, company, description, contact, created, updated)
		VALUES (:title, :company, :description, :contact, :created, :updated)
		RETURNING id`)

	if err != nil {
		log.Print(err)
		ctx.ServeStatus(http.StatusInternalServerError)
		return
	}

	var id int
	if err := stmt.Get(&id, job); err != nil {
		log.Print(err)
		ctx.ServeStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect("/jobs/"+strconv.Itoa(id), http.StatusFound)
}

// Show looks for a particular job and shows its full details.
func (j *JobHandlers) Show(ctx *cobalt.Context) {
	id, err := strconv.Atoi(ctx.ParamValue("id"))

	if err != nil {
		log.Print(err)
		ctx.ServeStatus(http.StatusBadRequest)
		return
	}

	var job Job

	if err := j.db.Get(&job, "SELECT * FROM jobs WHERE id = $1", id); err != nil {
		log.Print(err)
		ctx.ServeStatus(http.StatusNotFound)
		return
	}

	ctx.ServeHTML("jobs/show", job)
}
