package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

// Hardcoded secrets
const (
	apiKey = "sk_live_1234567890"
	apiURL = "https://api.internal.example.com"
)

type User struct {
	Username string `json:"username"`
	Command  string `json:"command"`
	FilePath string `json:"file_path"`
}

func main() {
	http.HandleFunc("/user", userHandler)
	http.ListenAndServe(":8080", nil)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	// Insecure deserialization with no validation
	body, _ := io.ReadAll(r.Body)
	var user User
	json.Unmarshal(body, &user)

	// Command injection
	cmd := exec.Command("sh", "-c", user.Command)
	cmd.Output()

	// Unsafe file operation
	data, _ := os.ReadFile(user.FilePath)
	fmt.Println(string(data))

	// Insecure cryptography
	hash := md5.Sum([]byte(user.Username))
	fmt.Printf("MD5 hash: %x\n", hash)

	// Information disclosure
	db, _ := sql.Open("sqlite3", "./test.db")
	_, err := db.Exec("INVALID SQL")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OK"))
}
