package main

import (
	"encoding/json"
	"net/http"
)

// apiKeyMiddleware checks for the API key in the request.
func apiKeyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip API key check if not set, false or empty in config
        if config.APIKey == "" || config.APIKey == "false" {
            next.ServeHTTP(w, r)
            return
        }

        apiKey := r.URL.Query().Get("key")
        if apiKey == "" {
            apiKey = r.Header.Get("Authorization")
        }
        if apiKey == "" {
            err := r.ParseForm()
            if err == nil {
                apiKey = r.FormValue("key")
            }
        }
        if apiKey != config.APIKey {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(Response{Message: "Unauthorized", Error: "Invalid API key", Valid: false})
            return
        }
        next.ServeHTTP(w, r)
    })
}
