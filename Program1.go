package main

import (
	"fmt"
	"github.com/labstack/echo"
	"html/template"
	"net/http"
	"sync"
)

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./assets/index.html"))
	tmpl.Execute(w, nil)
}

func testServer(c echo.Context) error {
	return c.String(http.StatusOK, "server is up!")
}

func main() {

	// Create 2 WaitGroups for 2 GoRoutines
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// GoRoutine server
	go func() {
		fmt.Println("Running server...")
		server := echo.New()
		server.GET("/", testServer)
		server.Start(":8000")
		wg.Done()
	}()

	// GoRoutine WebHost
	go func() {
		fmt.Println("Running webpage...")
		http.HandleFunc("/", index)
		http.ListenAndServe(":3000", nil)
		wg.Done()
	}()

	wg.Wait()
}
