package twark

import (
  "fmt"
  "strings"
)

type StringMap map[string]string

func (strmap StringMap) String() string {
  // custom string representation for string->string maps
  pairs := make([]string, len(strmap))
  // is there really no easy way to iterate a map's keys, values, and their indices all together?
  i := 0
  for key, val := range strmap {
    // pairs := make([]string, 0, len(m))
    // pairs = append(pairs, fmt.Sprintf("%s: %s", key, val))
    pairs[i] = fmt.Sprintf("%s: %s", key, val)
    i++
  }
  return fmt.Sprint("{ ", strings.Join(pairs, ", "), " }")
}
