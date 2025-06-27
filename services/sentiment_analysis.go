package services

import (
	"strings"
)

// Enhanced sentiment keywords with weighted scores
var sentimentKeywords = map[string]int{
	// Bullish keywords (positive scores)
	"buy": 3, "moon": 4, "bullish": 4, "long": 2, "pump": 3, "up": 2,
	"rocket": 4, "hodl": 3, "gain": 3, "profit": 3, "surge": 4,
	"rally": 3, "breakout": 4, "support": 2, "bull": 4, "rise": 3,
	"upward": 3, "green": 2, "optimistic": 3, "positive": 2, "strong": 2,
	"accumulate": 3, "investment": 2, "opportunity": 3, "recovery": 3,

	// Bearish keywords (negative scores)
	"sell": -3, "bear": -4, "bearish": -4, "short": -2, "dump": -4,
	"down": -2, "crash": -5, "drop": -3, "loss": -3, "dip": -2,
	"resistance": -2, "panic": -4, "correction": -2, "decline": -3,
	"fall": -3, "falling": -3, "red": -2, "negative": -2, "weak": -2,
	"exit": -2, "fear": -3, "bubble": -3, "overvalued": -3, "risky": -2,
}

// CalculateSentiment calculates sentiment score based on keyword analysis with weighted scoring
func CalculateSentiment(redditData []string) float64 {
	if len(redditData) == 0 {
		return 50.0 // neutral if no data
	}

	totalScore := 0
	postsWithSentiment := 0

	for _, post := range redditData {
		postScore := 0
		words := strings.Fields(strings.ToLower(post))

		// Analyze each word in the post
		for _, word := range words {
			// Remove common punctuation
			word = strings.Trim(word, ".,!?;:()[]{}\"'")

			if score, exists := sentimentKeywords[word]; exists {
				postScore += score
			}
		}

		// Cap individual post scores to prevent outliers from skewing results
		if postScore > 5 {
			postScore = 5
		} else if postScore < -5 {
			postScore = -5
		}

		// Only count posts that have some sentiment indicators
		if postScore != 0 {
			postsWithSentiment++
		}

		totalScore += postScore
	}

	// If no posts had sentiment indicators, return neutral
	if postsWithSentiment == 0 {
		return 50.0
	}

	// Calculate average sentiment per post
	avgSentiment := float64(totalScore) / float64(postsWithSentiment)

	// Convert to percentage (0-100) where:
	// -5 (max negative) = 0%
	// 0 (neutral) = 50%
	// +5 (max positive) = 100%
	percentage := 50 + (avgSentiment * 10)

	// Ensure bounds
	if percentage > 100 {
		percentage = 100
	} else if percentage < 0 {
		percentage = 0
	}

	return percentage
}
