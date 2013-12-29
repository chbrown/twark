package twark

import (
  "fmt"
  "github.com/ChimeraCoder/anaconda"
  "github.com/chbrown/twark/csv_with_header"
  "math/rand"
  "net/url"
  "time"
)

type Headers map[string]string

func (headers Headers) Values() (values url.Values) {
  values = make(url.Values, len(headers))
  for key, val := range headers {
    values[key] = []string{val}
  }
  return
}

func ChooseApi(accounts_filepath string) *anaconda.TwitterApi {
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
