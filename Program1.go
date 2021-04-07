package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"net/http"
	"sync"
)

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./assets/index.html"))
	tmpl.Execute(w, nil)
}

func testServer(c echo.Context) error {
	return c.JSON(http.StatusOK, "server is up!")
}

func getWelcome(c echo.Context) error {
	return c.String(http.StatusOK, "Thanks for visiting the server")
}

func main() {

	// Create 2 WaitGroups for 2 GoRoutines
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// GoRoutine server
	go func() {
		fmt.Println("Running server...")
		server := echo.New()
		server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}))
		server.GET("/", testServer)
		server.GET("/getWelcome", getWelcome)
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
