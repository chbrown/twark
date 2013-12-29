package twark

import (
  "fmt"
  "github.com/ChimeraCoder/anaconda"
  "github.com/coopernurse/gorp"
  // "github.com/vmihailenco/pg" // -> http://godoc.org/github.com/vmihailenco/pg
  "log"
  "net/url"
  "strconv"
  "time"
  // _ "github.com/jbarham/gopgsqldriver"
  // pgsql "github.com/jbarham/pgsql.go"
  "database/sql"
  _ "github.com/lib/pq"
)

// vmihailenco's:
// db := pg.Connect(&pg.Options{
//   // User: "postgres",
//   Database: "twark",
// })
// defer db.Close()

type Task struct {
  Id                int64 // `db:"id"`
  Screen_name       string
  Last_updated      time.Time
  User_fetched      bool
  Backlog_exhausted bool
  Touched           time.Time
  Inserted          time.Time
}

type DB struct {
  sql.DB
}

func NewDB() *DB {
  // var db *DB
  // var err error
  db, err := sql.Open("postgres", "dbname=twark sslmode=disable") // user=postgres
  if err != nil {
    panic(err)
  }
  return &DB{*db}
}

func Zip(keys []string, vals []interface{}) map[string]interface{} {
  // assert len(keys) == len(vals)
  m := make(map[string]interface{}, len(keys))
  for i, key := range keys {
    m[key] = vals[i]
  }
  return m
}

func (db *DB) QueryMap(query string, args ...interface{}) (map[string]interface{}, error) {
  rows, err := db.Query(query, args...)
  if err != nil {
    return nil, err
  }

  columns, err := rows.Columns()
  if err != nil {
    return nil, err
  }

  if ok := rows.Next(); !ok {
    return nil, nil
  }

  cells := make([]interface{}, len(columns))
  pointers := make([]interface{}, len(columns))
  for i, _ := range columns {
    pointers[i] = &cells[i]
  }
  if err := rows.Scan(pointers...); err != nil {
    return nil, err
  }
  fmt.Println("C&P", cells, pointers)

  if err := rows.Err(); err != nil {
    return nil, err
  }

  return Zip(columns, cells), nil
}

func AddTask(screen_name string) error {
  db := NewDB()
  rows, err := db.Query("INSERT INTO tasks (screen_name) VALUES ($1) RETURNING id", screen_name)
  // rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
  // if err, ok := err.(*pq.Error), ok {
  // pgErr := err.(*pq.Error)
  // .Code.Name()
  // log.Fatal(pgErr)
  // fmt.Println("pq error:", err)
  // panic(err)
  // }
  if err != nil {
    return err
  }

  for rows.Next() {
    var id int
    if err := rows.Scan(&id); err != nil {
      return err
    }
    fmt.Printf("Added task #%d\n", id)
  }
  return rows.Err()
}

func WorkTasks(api *anaconda.TwitterApi) error {
  db := NewDB()
  // find an unsaturated task:
  res, err := db.QueryMap("SELECT id, screen_name, user_fetched FROM tasks WHERE NOT backlog_exhausted LIMIT 1")
  if err != nil {
    log.Fatal(err)
    return err
  }

  log.Println("Next task:", res)
  // var id int;
  // for rows.Next() {
  //   if err := rows.Scan(&id); err != nil {
  //     return err
  //   }
  //   fmt.Printf("Added task #%d\n", id)
  // }
  // if err := rows.Err(); err != nil {
  //   return err
  // }
  return nil

  screen_name := res["screen_name"].(string)

  // Options: https://dev.twitter.com/docs/api/1.1/get/statuses/user_timeline
  headers := Headers{
    // "user_id":      int,
    "screen_name":      screen_name,
    "include_entities": "true",
    // "since_id":           int,
    "count": "200",
    // "max_id":           int,
    "trim_user":           "true",
    "exclude_replies":     "false",
    "contributor_details": "true",
    "include_rts":         "true",
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

  return nil
}

func PrintTimeline(api *anaconda.TwitterApi) {
  // api := RandomApi()
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
