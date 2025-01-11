package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/qich3n/crypto-sentiment-analyzer/services"
)

// RenderIndex renders the index.html template
func RenderIndex(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", nil)
}

// GetSentimentAnalysis handles the GET request to fetch sentiment data
// Example usage: /api/analyze?coin=bitcoin
func GetSentimentAnalysis(c *gin.Context) {
    coin := c.Query("coin")
    if coin == "" {
        coin = "bitcoin" // default coin if no query provided
    }

    // Fetch data from Reddit
    redditData, err := services.FetchRedditData(coin)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // (Optional) Fetch data from Twitter
    // twitterData, _ := services.FetchTwitterData(coin)

    // Perform sentiment analysis (simple approach)
    sentimentPercentage := services.CalculateSentiment(redditData /*, twitterData*/)

    c.JSON(http.StatusOK, gin.H{
        "coin":               coin,
        "sentimentPercent":   sentimentPercentage,
        "sentimentDirection": directionLabel(sentimentPercentage),
    })
}

// directionLabel is a helper to return "Bullish" or "Bearish"
func directionLabel(percentage float64) string {
    if percentage >= 50 {
        return "Bullish"
    }
    return "Bearish"
}
