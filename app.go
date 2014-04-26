package main
/*TODO: 1) construct links in the message; 2) extract repo descriptions*/

import (
  "fmt"
  "github.com/PuerkitoBio/goquery"
  "io/ioutil"
  "os"
  "encoding/json"
  "strings"
  "net/http"
  "net/url"
)

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func main() {
  dir, err := os.Getwd() 
  check(err)
  oldItems := []string{}
  if _, err := os.Stat(dir + "\\" + "old.json"); err == nil {
    fmt.Print("if")
    file, err := ioutil.ReadFile("old.json")
    check(err)
    json.Unmarshal(file, &oldItems)
  } else {
    fmt.Print("else")
  }
  fmt.Print(oldItems)

  doc, err := goquery.NewDocument("https://github.com/trending")
  check(err)
  currentItems := []string{}
  doc.Find(".repo-leaderboard-list-item").Each(func(i int, s *goquery.Selection) {
    repo := s.Find("h2").Find(".repository-name")
    currentItems = append(currentItems, repo.Text())
    fmt.Printf("repo: %s\n", repo.Text())
  })

  newItems := []string{}
  var isDuplicate bool
  for _, currentValue := range currentItems {
    isDuplicate = false
    for _, oldValue := range oldItems {
      if oldValue == currentValue {
        isDuplicate = true
        break
      }
    }
    if isDuplicate == false {
      newItems = append(newItems, currentValue)
    }
  }
  jsonItems, _ := json.Marshal(append(oldItems, newItems...))
  err = ioutil.WriteFile("old.json", jsonItems, 'l')
  check(err)

  message := strings.Join(newItems, "\n\n")
  _, err := http.PostForm("https://api.pushover.net/1/messages.json",
    url.Values{
      "token": {"aiCnXzdxjSYxjLdcba8aTmk7DVTKFZ"}, 
      "user": {"ho8BieZry6ec5Ngft93AalCENwqxEk"}, 
      "message": {message},
    })
}
