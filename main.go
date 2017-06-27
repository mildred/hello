package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Write([]byte("Hello World!\n"))

	fmt.Fprintf(w, "os.Args: %#v\n", os.Args)

	addrs, err := net.DefaultResolver.LookupAddr(ctx, "consul.service.consul")
	if err != nil {
		fmt.Fprintf(w, "Consul DNS: consul.service.consul A error: %s\n", err)
	} else {
		for _, addr := range addrs {
			fmt.Fprintf(w, "Consul DNS: consul.service.consul A %s\n", addr)
		}
	}

	cname, srvs, err := net.DefaultResolver.LookupSRV(ctx, "consul", "tcp", "service.consul")
	if err != nil {
		fmt.Fprintf(w, "Consul DNS: _consul._tcp.service.consul SRV error: %s\n", err)
	} else {
		fmt.Fprintf(w, "Consul DNS: _consul._tcp.service.consul CNAME %v\n", cname)
		for _, srv := range srvs {
			fmt.Fprintf(w, "Consul DNS: _consul._tcp.service.consul SRV %v\n", srv)
		}
	}

	url := "http://consul.service.consul:8500/v1/catalog/services"
	fmt.Fprintf(w, "\nConsul: > GET %s\n", url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(w, "Consul: Error %s\n", err)
	} else {
		fmt.Fprintf(w, "Consul: < %s\n", res.Status)
		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintf(w, "Consul: Error %s\n", err)
		} else {
			fmt.Fprintf(w, "Consul services:\n%s\n", string(content))
		}
	}
	fmt.Fprintf(w, "\n\n")

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

	db_res := db.QueryRow("SELECT count(*) FROM counter")

	var count int
	db_res.Scan(&count)

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
