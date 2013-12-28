package main

import (
  "errors"
  "flag"
  "fmt"
  "github.com/chbrown/twark"
  "os/user"
  "path"
)

// var verboseFlag = flag.Bool("verbose", false, "Print extra output")
// var versionFlag = flag.Bool("version", false, "Print version and exit")

var actionName string
var accountsFilepath string

func init() {
  u, err := user.Current()
  if err != nil {
    panic(err)
  }

  // Alternatively: if path[:2] == "~/" { path = strings.Replace(path, "~/", dir, 1) }
  defaultAccountsFilepath := path.Join(u.HomeDir, ".twitter")
  flag.StringVar(&actionName, "action", "work", "Twark action; one of: 'work'")
  flag.StringVar(&accountsFilepath, "accounts", defaultAccountsFilepath, "File containing table of OAuth credentials")
}

func main() {
  flag.Parse()

  if actionName == "work" {
    screen_names := flag.Args()
    api := twark.ChooseApi(accountsFilepath)
    for _, screen_name := range screen_names {
      fmt.Println("Fetching Twitter user:", screen_name)
      twark.FetchUser(api, screen_name)
    }
  } else {
    err := errors.New("Unrecognized action: " + actionName)
    panic(err)
  }
}
