package twark

import (
  "fmt"
  "net/url"
  "strconv"
)

func FetchUser(screen_name string) {
  api := RandomApi()
  // fmt.Println(api)
  headers := Headers{
    "screen_name":         screen_name,
    "include_entities":    "true",
    "include_rts":         "true",
    "contributor_details": "true",
    "count":               "200",
  }

  ntweets := 0
  nresponses_failed := 0
  nresponses_empty := 0
  max_bad_responses := 10
  // responses is a list of integers that record how many tweets we got back from Twitter for each request
  // this is useful so that we can retry a couple times on empty responses, but not forever
  for nresponses_failed+nresponses_empty <= max_bad_responses {
    fmt.Println("Getting user timeline:", headers)
    tweets, err := api.GetUserTimeline(headers.Values())

    if err != nil {
      nresponses_failed++
      fmt.Println(err)
    }

    if len(tweets) == 0 {
      nresponses_empty++
    } else {
      last_tweet := tweets[len(tweets)-1]
      last_id, err := strconv.ParseUint(last_tweet.Id_str, 10, 64)
      if err != nil {
        panic(err)
      }
      // yes, explicit string() conversion is necessary
      headers["max_id"] = strconv.FormatUint(last_id-1, 10)
    }

    for _, tweet := range tweets {
      fmt.Println(tweet.Id_str, tweet.Created_at, tweet.Text)
      ntweets++
    }
  }
  fmt.Printf("Downloaded %d tweets for screen name: %q\n", ntweets, screen_name)
}

func PrintTimeline() {
  api := RandomApi()
  // fmt.Println(api)

  // searchResult, _ := api.GetSearch("golang", nil)

  // timeline, err := api.GetHomeTimeline()
  timeline, err := api.GetUserTimeline(url.Values{
    "screen_name": []string{"chbrown"},
    "count":       []string{"100"},
  })
  if err != nil {
    panic(err)
  }

  // searchResult, _ := api.GetSearch("golang", nil)
  for _, tweet := range timeline {
    // type Tweet struct:
    // Source        string
    // Id            int64
    // Retweeted     bool
    // Favorited     bool
    // User          TwitterUser
    // Truncated     bool
    // Text          string
    // Retweet_count int64
    // Id_str        string
    // Created_at    string
    // Entities      TwitterEntities
    // fmt.Println(tweet.Text)
    // tweet.User
    // tweet.Entities
    fmt.Println(
      tweet.Source,
      tweet.Id,
      tweet.Retweeted,
      tweet.Favorited,
      tweet.Truncated,
      tweet.Text,
      tweet.Retweet_count,
      tweet.Id_str,
      tweet.Created_at)
  }
}
