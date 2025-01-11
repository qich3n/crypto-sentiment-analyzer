package services

import (
    "strings"
)

// CalculateSentiment calculates how many posts/tweets are "bullish" vs "bearish"
// In real usage, you'd use NLP or a more robust sentiment analysis approach.
func CalculateSentiment(redditData []string /*, twitterData []string*/) float64 {
    // Combine data
    combinedData := redditData
    // combinedData = append(combinedData, twitterData...)

    if len(combinedData) == 0 {
        return 50.0 // neutral if no data
    }

    bullishCount := 0
    totalCount := len(combinedData)

    // Simple naive approach: check for bullish keywords
    bullishKeywords := []string{"buy", "moon", "bullish", "long", "pump", "up"}
    // Bearish keywords could be used in a more advanced approach

    for _, post := range combinedData {
        for _, keyword := range bullishKeywords {
            if strings.Contains(post, keyword) {
                bullishCount++
                break
            }
        }
    }

    // Calculate the bullish sentiment percentage
    sentiment := (float64(bullishCount) / float64(totalCount)) * 100
    return sentiment
}
