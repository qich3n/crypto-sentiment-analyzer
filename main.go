package main

import (
    "log"

    "github.com/qich3n/crypto-sentiment-analyzer/router"
)

func main() {
    // Create the Gin router from our router package
    r := router.SetupRouter()

    // Start the server on port 8080 (or any port you prefer)
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Could not run the server: %v", err)
    }
}
