package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!\n"))

	fmt.Fprintf(w, "os.Args: %#v\n", os.Args)

	dburl := fmt.Sprintf("%s://%s:%s@%s:%s/%s", os.Getenv("DB_ENGINE"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	fmt.Fprintf(w, "Connect to %s\n", dburl)
	db, err := sql.Open("postgres", dburl)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		return
	}
	defer db.Close()
	fmt.Fprintf(w, "Connected.\n")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS counter (count integer);")
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		return
	}

	_, err = db.Exec("INSERT INTO counter SELECT count(*) FROM counter")
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		return
	}

	res := db.QueryRow("SELECT count(*) FROM counter")

	var count int
	res.Scan(&count)

	fmt.Fprintf(w, "count: %#v\n", count)

	fmt.Fprintf(w, "\nEnvironment:\n")
	for _, env := range os.Environ() {
		fmt.Fprintf(w, "%s\n", env)
	}
	fmt.Fprintf(w, "\n")

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		return
	}
	fmt.Fprintf(w, "Generated on %s\n", hostname)
}

func main() {
	fmt.Println("hello world")
	http.HandleFunc("/", HelloWorld)
	log.Fatal(http.ListenAndServe(":80", nil))
}
