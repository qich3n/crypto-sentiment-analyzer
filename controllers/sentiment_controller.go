package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qich3n/crypto-sentiment-analyzer/services"
)

// RenderIndex renders the index.html template
func RenderIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// HealthCheck returns the status of Reddit service
func HealthCheck(c *gin.Context) {
	// Check Reddit service
	redditService, err := services.NewRedditService()
	redditStatus := err == nil

	if redditStatus {
		// Quick connection test
		if testErr := redditService.TestConnection(); testErr != nil {
			redditStatus = false
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"services": gin.H{
			"reddit": redditStatus,
		},
		"data_source": "Reddit (r/CryptoCurrency)",
		"timestamp":   time.Now(),
	})
}

// GetSentimentAnalysis handles sentiment analysis for a given cryptocurrency
func GetSentimentAnalysis(c *gin.Context) {
	coin := c.Query("coin")
	if coin == "" {
		coin = "bitcoin" // default coin if no query provided
	}

	// Initialize Reddit service
	redditService, err := services.NewRedditService()
	if err != nil {
		log.Printf("Failed to initialize Reddit service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Reddit service unavailable: " + err.Error(),
		})
		return
	}

	// Fetch Reddit data
	redditData, err := redditService.FetchRedditData(coin)
	if err != nil {
		log.Printf("Reddit API error for %s: %v", coin, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch Reddit data: " + err.Error(),
		})
		return
	}

	// Check if we got any data
	if len(redditData) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"coin":               coin,
			"sentimentPercent":   50.0, // Neutral
			"sentimentDirection": "Neutral",
			"posts":              0,
			"message":            "No recent posts found for this cryptocurrency",
			"data_source":        "Reddit (r/CryptoCurrency)",
			"timestamp":          time.Now(),
		})
		return
	}

	// Calculate sentiment
	sentimentPercent := services.CalculateSentiment(redditData)

	// Prepare response
	response := gin.H{
		"coin":               coin,
		"sentimentPercent":   sentimentPercent,
		"sentimentDirection": directionLabel(sentimentPercent),
		"posts":              len(redditData),
		"data_source":        "Reddit (r/CryptoCurrency)",
		"timestamp":          time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// GetTrendingAnalysis analyzes sentiment for multiple popular cryptocurrencies
func GetTrendingAnalysis(c *gin.Context) {
	popularCoins := []string{"bitcoin", "ethereum", "solana", "cardano", "polkadot", "chainlink"}
	var trending []gin.H

	// Initialize Reddit service
	redditService, err := services.NewRedditService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Reddit service unavailable: " + err.Error(),
		})
		return
	}

	for _, coin := range popularCoins {
		// Get Reddit data
		redditData, err := redditService.FetchRedditData(coin)
		if err != nil {
			log.Printf("Reddit error for %s: %v", coin, err)
			continue // Skip this coin if there's an error
		}

		if len(redditData) > 0 {
			sentiment := services.CalculateSentiment(redditData)
			trending = append(trending, gin.H{
				"coin":      coin,
				"sentiment": sentiment,
				"direction": directionLabel(sentiment),
				"posts":     len(redditData),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"trending":    trending,
		"data_source": "Reddit (r/CryptoCurrency)",
		"timestamp":   time.Now(),
	})
}

// directionLabel helper function to return sentiment direction
func directionLabel(percentage float64) string {
	if percentage >= 60 {
		return "Bullish"
	} else if percentage <= 40 {
		return "Bearish"
	}
	return "Neutral"
}
