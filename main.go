package main

func main() {
	db := NewDatabase("app.db")
	db.Init()
	defer db.db.Close()

	server := NewServer(db)
	server.Run()
}
