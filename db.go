package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"todo/components"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(source string) *Database {
	db, err := sql.Open("sqlite", source)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("database connected succesfully")

	return &Database{
		db: db,
	}
}

// create table
func (d *Database) Init() {
	d.db.SetMaxOpenConns(1)

	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS todo (
			id INTEGER PRIMARY KEY,
			title TEXT,
			type TEXT
		)
	`)

	if err != nil {
		d.db.Close()
		log.Fatal(err)
	}

	log.Println("table created successfully")
}

func (d *Database) create_todo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.PostFormValue("title")
	title_type := r.PostFormValue("type")

	if len(title) == 0 || len(title_type) == 0 {
		errString := "Title or Type is required"
		status := components.CreateStatus(false, errString)
		status.Render(r.Context(), w)
		return
	}

	result, err := d.db.Exec(`INSERT INTO todo (title, type) VALUES (?, ?)`, title, title_type)
	if err != nil {
		fmt.Fprintf(w, "Error inserting new todo")
		log.Println(err)
	}

	getID, _ := result.LastInsertId()
	status := components.CreateStatus(true, title)
	status.Render(r.Context(), w)

	log.Println("New Todo Insert ID:", getID)
}
