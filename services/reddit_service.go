package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// RedditAPIResponse represents the structure of Reddit API response
type RedditAPIResponse struct {
	Data struct {
		Children []struct {
			Data struct {
				Title     string `json:"title"`
				SelfText  string `json:"selftext"`
				Score     int    `json:"score"`
				Subreddit string `json:"subreddit"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// RedditService handles Reddit API interactions
type RedditService struct {
	clientID     string
	clientSecret string
	accessToken  string
	httpClient   *http.Client
}

// NewRedditService creates a new Reddit service with environment variables
func NewRedditService() (*RedditService, error) {
	clientID := strings.TrimSpace(os.Getenv("REDDIT_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("REDDIT_CLIENT_SECRET"))

	if clientID == "" || clientSecret == "" {
		return nil, errors.New("Reddit API credentials not found in environment variables")
	}

	return &RedditService{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{},
	}, nil
}

// authenticate gets an access token from Reddit
func (rs *RedditService) authenticate() error {
	authString := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", rs.clientID, rs.clientSecret)))

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(
		"POST",
		"https://www.reddit.com/api/v1/access_token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Basic "+authString)
	req.Header.Add("User-Agent", "CryptoSentimentAnalyzer/1.0")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := rs.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Reddit authentication failed: %s", string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	rs.accessToken = result.AccessToken
	return nil
}

// FetchRedditData fetches posts from Reddit with proper authentication
func (rs *RedditService) FetchRedditData(coin string) ([]string, error) {
	// Authenticate if we don't have a token
	if rs.accessToken == "" {
		if err := rs.authenticate(); err != nil {
			return nil, fmt.Errorf("failed to authenticate with Reddit: %v", err)
		}
	}

	// Search in r/CryptoCurrency
	url := fmt.Sprintf("https://oauth.reddit.com/r/CryptoCurrency/search.json?q=%s&restrict_sr=1&limit=100&sort=new",
		url.QueryEscape(coin))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Use OAuth token for authenticated requests
	req.Header.Set("Authorization", "Bearer "+rs.accessToken)
	req.Header.Set("User-Agent", "CryptoSentimentAnalyzer/1.0")

	resp, err := rs.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		// Token expired, re-authenticate
		if err := rs.authenticate(); err != nil {
			return nil, fmt.Errorf("failed to re-authenticate with Reddit: %v", err)
		}
		// Retry the request
		return rs.FetchRedditData(coin)
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Reddit API error (%d): %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var redditResp RedditAPIResponse
	if err := json.Unmarshal(body, &redditResp); err != nil {
		return nil, err
	}

	// Collect titles and self text for sentiment analysis
	var texts []string
	for _, child := range redditResp.Data.Children {
		title := strings.ToLower(child.Data.Title)
		selfText := strings.ToLower(child.Data.SelfText)

		// Combine title and self text for better sentiment analysis
		combined := title
		if selfText != "" && selfText != title {
			combined += " " + selfText
		}

		// Only include posts with some content
		if strings.TrimSpace(combined) != "" {
			texts = append(texts, combined)
		}
	}

	return texts, nil
}

// TestConnection tests the Reddit API connection
func (rs *RedditService) TestConnection() error {
	if err := rs.authenticate(); err != nil {
		return err
	}

	// Try a simple request
	_, err := rs.FetchRedditData("bitcoin")
	return err
}

// FetchRedditData is the original function signature for backward compatibility
func FetchRedditData(coin string) ([]string, error) {
	service, err := NewRedditService()
	if err != nil {
		return nil, err
	}

	return service.FetchRedditData(coin)
}
