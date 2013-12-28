package twark

import (
  "fmt"
  "github.com/ChimeraCoder/anaconda"
  "github.com/chbrown/twark/csv_with_header"
  "math/rand"
  "net/url"
  "os/user"
  "path"
  "time"
)

func init() {
  // go's rand uses the same seed on each run by default (weird!)
  rand.Seed(time.Now().UTC().UnixNano())
}

func PrintTimeline() {
  api := RandomApi()
  // fmt.Println(api)

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

func RandomApi() *anaconda.TwitterApi {
  u, err := user.Current()
  if err != nil {
    panic(err)
  }

  // Alternatively: if path[:2] == "~/" { path = strings.Replace(path, "~/", dir, 1) }
  accounts_filepath := path.Join(u.HomeDir, ".twitter")

  accounts := csv_with_header.ReadAll(accounts_filepath)

  now_source := rand.NewSource(time.Now().UTC().UnixNano())
  now_rand := rand.New(now_source)
  account_index := now_rand.Intn(len(accounts))
  account := accounts[account_index]
  fmt.Printf("Using account #%d: %q\n", account_index, account["screen_name"])

  anaconda.SetConsumerKey(account["consumer_key"])
  anaconda.SetConsumerSecret(account["consumer_secret"])

  api := anaconda.NewTwitterApi(account["access_token"], account["access_token_secret"])

  return &api
}
