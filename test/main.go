package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Connect through the proxy on port 5433
	connStr := "host=localhost port=5433 user=postgres password=postgres dbname=pgrelay_test sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("open:", err)
	}
	defer db.Close()

	// Ping
	if err := db.Ping(); err != nil {
		log.Fatal("ping:", err)
	}
	fmt.Println("connected through proxy")

	// Create a test table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS proxy_test (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		value INT
	)`)
	if err != nil {
		log.Fatal("create table:", err)
	}
	fmt.Println("created table proxy_test")

	// Insert some rows
	for i := 0; i < 5; i++ {
		_, err := db.Exec("INSERT INTO proxy_test (name, value) VALUES ($1, $2)",
			fmt.Sprintf("item_%d", i), i*10)
		if err != nil {
			log.Fatal("insert:", err)
		}
	}
	fmt.Println("inserted 5 rows")

	// Query rows back
	rows, err := db.Query("SELECT id, name, value FROM proxy_test ORDER BY id")
	if err != nil {
		log.Fatal("query:", err)
	}
	defer rows.Close()

	fmt.Println("rows:")
	for rows.Next() {
		var id, value int
		var name string
		if err := rows.Scan(&id, &name, &value); err != nil {
			log.Fatal("scan:", err)
		}
		fmt.Printf("  id=%d name=%s value=%d\n", id, name, value)
	}

	// Clean up
	_, err = db.Exec("DROP TABLE proxy_test")
	if err != nil {
		log.Fatal("drop table:", err)
	}
	fmt.Println("dropped table proxy_test")

	fmt.Println("all good — proxy works!")
}
