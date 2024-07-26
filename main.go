package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
    var err error
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbHost := os.Getenv("DB_HOST")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    var user map[string]string
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    username, password := user["username"], user["password"]

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }

    _, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, string(hashedPassword))
    if err != nil {
        http.Error(w, "User already exists", http.StatusConflict)
        return
    }

    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "User registered successfully")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var user map[string]string
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    username, password := user["username"], user["password"]

    var storedPassword string
    err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Login successful")
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    username := vars["username"]

    var storedID string
    err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&storedID)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    userInfo := map[string]string{
        "username": username,
        "ID": storedID,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userInfo)
}

func main() {
    initDB()
    defer db.Close()

    r := mux.NewRouter()
    r.HandleFunc("/api/register", RegisterHandler).Methods("POST")
    r.HandleFunc("/api/login", LoginHandler).Methods("POST")
    r.HandleFunc("/api/user/{username}", GetUserHandler).Methods("GET")

    http.Handle("/", r)
    fmt.Println("Server started at :8090")
    log.Fatal(http.ListenAndServe(":8090", nil))
}
