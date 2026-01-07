package services

import "strings"

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

	// Crypto slang (mixed sentiment)
	"rekt":  -3,
	"fomo":  2,
	"fud":   -3,
	"lambo": 3,
	"bag":   -1,
}

// negationWords invert the following sentiment for a short window
var negationWords = map[string]bool{
	"not": true, "no": true, "never": true, "dont": true, "don't": true,
	"cant": true, "can't": true, "isnt": true, "isn't": true,
}

// intensifierWords amplify the following sentiment for a short window
var intensifierWords = map[string]float64{
	"very": 1.5, "extremely": 1.8, "super": 1.4, "so": 1.3, "really": 1.3,
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

		negateWindow := 0        // how many upcoming sentiment words to negate
		intensityWindow := 0     // how many upcoming sentiment words to intensify
		currentMultiplier := 1.0 // current intensity multiplier

		// Analyze each word in the post
		for _, word := range words {
			// Remove common punctuation and normalize basic contractions
			word = strings.Trim(word, ".,!?;:()[]{}\"'")

			// Handle negations (e.g., "not bullish", "don't buy")
			if negationWords[word] {
				negateWindow = 3 // affect the next few sentiment-bearing words
				continue
			}

			// Handle intensifiers (e.g., "very bullish", "extremely bearish")
			if mult, ok := intensifierWords[word]; ok {
				intensityWindow = 2
				currentMultiplier = mult
				continue
			}

			// Look up sentiment for this word
			if baseScore, exists := sentimentKeywords[word]; exists {
				score := baseScore

				// Apply negation if active
				if negateWindow > 0 {
					score = -score
					negateWindow--
				}

				// Apply intensity if active
				if intensityWindow > 0 {
					// simple integer scaling to keep scores manageable
					scaled := float64(score) * currentMultiplier
					if scaled > 0 {
						score = int(scaled + 0.5)
					} else {
						score = int(scaled - 0.5)
					}
					intensityWindow--
					if intensityWindow == 0 {
						currentMultiplier = 1
					}
				}

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
