package main

import (
	"fmt"
	"net/http"
)

func main() {
    // Load configuration.
    err := loadConfig()
    if err != nil {
        logError("Error loading config:", err)
        return
    }

    // Initialize logger
    initLogger()

    // Run list update on program start
    err = updateList()
    if err != nil {
        logError("Error updating list on start:", err)
        return
    }

    // Load lists into memory at startup
    err = loadAllLists()
    if err != nil {
        logError("Error loading lists:", err)
        return
    }

    // Initialize cache and start garbage collection.
    InitCache()

    // Schedule list updates
    scheduleListUpdate()

    // Apply middleware
    handlerWithMiddleware := securityHeadersMiddleware(
        apiKeyMiddleware(
            fail2banMiddleware(
                rateLimitMiddleware(
                    http.HandlerFunc(handler),
                ),
            ),
        ),
    )

    http.Handle("/", handlerWithMiddleware)
    fmt.Println("Starting server on :80")
    if err := http.ListenAndServe(":80", nil); err != nil {
        logError("Error starting server:", err)
    }
}
