package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Response represents the structure of our JSON response.
type Response struct {
    Email   string            `json:"email"`
    Valid   bool              `json:"valid"`
    Message string            `json:"message"`
    Error   string            `json:"error,omitempty"`
    Cached  bool              `json:"cached"`
    Checks  map[string]bool   `json:"checks,omitempty"`
}

// handler handles the incoming requests and returns validation results.
func handler(w http.ResponseWriter, r *http.Request) {
    
    if r.URL.Path == "/favicon.ico" {
        http.NotFound(w, r)
        return
    }

    defer func() {
        if r := recover(); r != nil {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(Response{Message: "Internal Server Error", Error: fmt.Sprintf("%v", r), Valid: false})
        }
    }()

    var email string
    var err error
    noCache := false

    // Determine email from request.
    switch r.Method {
    case http.MethodGet:
        // Get email from path.
        parts := strings.Split(r.URL.Path, "/")
        if len(parts) > 1 && parts[1] != "" {
            email = parts[1]
        } else {
            // Get email from query param.
            email = r.URL.Query().Get("email")
        }
        // Check for no-cache query param.
        noCache = r.URL.Query().Get("nocache") == "true"
        
    case http.MethodPost:
        // Get email from request body.
        err = r.ParseForm()
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(Response{Email: email, Valid: false, Message: "Bad Request", Error: "Error parsing form"})
            return
        }
        email = r.FormValue("email")
        // Check for no-cache form param.
        noCache = r.FormValue("nocache") == "true"
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(Response{Email: email, Valid: false, Message: "Method Not Allowed"})
        return
    }

    // Check for no-cache header.
    if !noCache {
        noCache = r.Header.Get("Cache-Control") == "no-cache"
    }

    if email == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(Response{Email: email, Valid: false, Message: "Missing email parameter", Error: "No email provided in request"})
        return
    }

    // Validate email with all checks.
    isValid, message, errMessage, cached, checks := validateEmail(email, noCache)
    if isValid {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(Response{Email: email, Valid: true, Message: message, Cached: cached, Checks:  checks})
    } else {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(Response{Email: email, Valid: false, Message: message, Error: errMessage, Cached: cached, Checks:  checks})
        trackFailedLogin(r)
    }
}
