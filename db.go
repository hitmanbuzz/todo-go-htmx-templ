package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"todo/components"
	"todo/utils"

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
			type TEXT,
			duration TEXT
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
	duration := r.PostFormValue("duration")

	if len(title) == 0 || len(title_type) == 0 {
		errString := "Title or Type is required"
		status := components.CreateStatus(false, errString)
		status.Render(r.Context(), w)
		return
	}

	t := utils.ParseTime(duration)
	result, err := d.db.Exec(`INSERT INTO todo (title, type, duration) VALUES (?, ?, ?)`, title, title_type, t.UTC())
	if err != nil {
		fmt.Fprintf(w, "Error inserting new todo")
		log.Println(err)
	}

	getID, _ := result.LastInsertId()
	status := components.CreateStatus(true, title)
	status.Render(r.Context(), w)

	log.Printf("ID: %d - Title: %s\n", getID, title)
}

func (d *Database) show_page(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "page doesn't exist")
		return
	}

	data, err := d.fetchPageData(pageInt)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "data not found")
		return
	}

	totalRows := d.GetTotalRows()
	totalPage := totalRows / utils.LIMIT

	if r.Header.Get("Hx-Request") == "true" {
		show := components.Show(data, int(pageInt), int(totalPage))
		show.Render(r.Context(), w)
		return
	}

	show := components.Show(data, int(pageInt), int(totalPage))
	RenderLayout(r.Context(), w, show)
}

func (d *Database) GetTotalRows() int64 {
	var rowsCount int64

	rows := d.db.QueryRow("SELECT COUNT(*) from todo")
	err := rows.Scan(&rowsCount)
	if err != nil {
		return 0
	}

	return rowsCount
}

func (d *Database) fetchPageData(page int64) ([]utils.PageData, error) {
	var page_data []utils.PageData
	offset := (page - 1) * utils.LIMIT

	rows, err := d.db.Query(`
		SELECT * FROM todo
		ORDER BY id
		LIMIT ? OFFSET ?
	`, utils.LIMIT, offset)

	if err != nil {
		return page_data, err
	}

	for rows.Next() {
		var data utils.PageData
		if err := rows.Scan(&data.Id, &data.Title, &data.Ttype, &data.Duration); err != nil {
			return page_data, err
		}

		page_data = append(page_data, data)
	}

	log.Println("data fetched from db")
	return page_data, nil
}
