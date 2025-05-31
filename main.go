package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/qich3n/crypto-sentiment-analyzer/router"
	"github.com/qich3n/crypto-sentiment-analyzer/services"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found - using system environment variables")
	}

	// Validate required environment variables
	validateEnvironment()
}

func validateEnvironment() {
	required := []string{
		"REDDIT_CLIENT_ID",
		"REDDIT_CLIENT_SECRET",
	}

	missing := []string{}
	for _, env := range required {
		if os.Getenv(env) == "" {
			missing = append(missing, env)
		}
	}

	if len(missing) > 0 {
		log.Printf("⚠️  Warning: Missing required environment variables: %v", missing)
		log.Printf("📝 Please check your .env file or set these environment variables")
		log.Printf("🔗 Get Reddit API credentials from: https://www.reddit.com/prefs/apps")
	} else {
		log.Printf("✅ Environment variables loaded successfully")
	}
}

func testAPIConnection() {
	log.Printf("🔧 Testing Reddit API connection...")

	redditService, err := services.NewRedditService()
	if err != nil {
		log.Printf("❌ Reddit API initialization failed: %v", err)
		return
	}

	if err := redditService.TestConnection(); err != nil {
		log.Printf("❌ Reddit API connection failed: %v", err)
	} else {
		log.Printf("✅ Reddit API: Connected successfully")
	}
}

func main() {
	// Check if running in test mode
	if len(os.Args) > 1 && os.Args[1] == "test" {
		testAPIConnection()
		return
	}

	// Test API connection on startup
	testAPIConnection()

	// Create the Gin router
	r := router.SetupRouter()

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Crypto Sentiment Analyzer starting on port %s", port)
	log.Printf("📊 Data source: Reddit (r/CryptoCurrency)")
	log.Printf("🌐 Open http://localhost:%s to view the application", port)
	log.Printf("📖 Run with 'go run main.go test' to test Reddit API connection")

	// Start the server
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Could not start the server: %v", err)
	}
}
