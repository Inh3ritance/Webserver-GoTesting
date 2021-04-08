package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
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

func encrypt(secretMessage string, key rsa.PublicKey) string {
    label := []byte("OAEP Encrypted")
    rng := rand.Reader
    ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(secretMessage), label)
    CheckError(err)
	fmt.Println("CipherText:", string(ciphertext))
    return base64.StdEncoding.EncodeToString(ciphertext)
}

func decrypt(cipherText string, privKey rsa.PrivateKey) string {
    ct, _ := base64.StdEncoding.DecodeString(cipherText)
    label := []byte("OAEP Encrypted")
    rng := rand.Reader
    plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
    CheckError(err)
    fmt.Println("Plaintext:", string(plaintext))
    return string(plaintext)
}

func CheckError(e error) {
	if e != nil {
		fmt.Println(e.Error)
	}
}

func main() {
	
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	CheckError(err)
	
	publickey := privatekey.PublicKey
	str := "hello"

	str = encrypt(str, publickey)
	str = decrypt(str, *privatekey)

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
