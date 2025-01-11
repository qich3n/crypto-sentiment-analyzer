package services

import (
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "strings"
)

// For demonstration, we define a structure that might parse some of the JSON response
type RedditAPIResponse struct {
    Data struct {
        Children []struct {
            Data struct {
                Title string `json:"title"`
            } `json:"data"`
        } `json:"children"`
    } `json:"data"`
}

// FetchRedditData fetches posts from Reddit (subreddit of your choice: e.g., r/CryptoCurrency)
func FetchRedditData(coin string) ([]string, error) {
    // This is a public endpoint for demonstration. Actual usage might require OAuth with Reddit if necessary.
    // The user might have to pass a custom user agent, or set up a developer app in Reddit to avoid rate-limiting.
    url := fmt.Sprintf("https://www.reddit.com/r/CryptoCurrency/search.json?q=%s&restrict_sr=1", coin)

    // Make request
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    // Provide a User-Agent to avoid 429 error from Reddit
    req.Header.Set("User-Agent", "GoRedditClient/1.0")

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, errors.New("failed to fetch data from Reddit")
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var redditResp RedditAPIResponse
    if err := json.Unmarshal(body, &redditResp); err != nil {
        return nil, err
    }

    // Collect titles
    var titles []string
    for _, child := range redditResp.Data.Children {
        title := child.Data.Title
        // Filter out or process the title if needed
        titles = append(titles, strings.ToLower(title))
    }

    return titles, nil
}
