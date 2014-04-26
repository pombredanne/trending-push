package main

/*
TODO: 
- move pushover data to command line parameters
*/

import (
  "github.com/PuerkitoBio/goquery"
  "io/ioutil"
  "os"
  "encoding/json"
  "net/http"
  "net/url"
  "fmt"
  "runtime"
)

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func homeDir() string {
  if runtime.GOOS == "windows" {
    home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
    if home == "" {
      home = os.Getenv("USERPROFILE")
    }
    return home
  }
  return os.Getenv("HOME")
}

func main() {
  pushbulletToken := os.Args[1]
  pushbulletDevice:= os.Args[2]
  dir := homeDir() 
  oldItems := []string{}
  if _, err := os.Stat(dir + "/" + ".trending-push"); err == nil {
    file, err := ioutil.ReadFile(dir + "/" + ".trending-push")
    check(err)
    json.Unmarshal(file, &oldItems)
  }

  doc, err := goquery.NewDocument("https://github.com/trending")
  check(err)
  currentItems := []string{}
  currentDescriptions:= []string{}
  doc.Find(".repo-leaderboard-list-item").Each(func(i int, s *goquery.Selection) {
    repo := s.Find("h2").Find(".repository-name")
    description := s.Find(".repo-leaderboard-description")
    currentItems = append(currentItems, repo.Text())
    currentDescriptions = append(currentDescriptions, description.Text())
  })
  newItems := []string{}
  newDescriptions := []string{}
  var isDuplicate bool
  for currentIndex, currentValue := range currentItems {
    isDuplicate = false
    for _, oldValue := range oldItems {
      if oldValue == currentValue {
        isDuplicate = true
        break
      }
    }
    if isDuplicate == false {
      newItems = append(newItems, currentValue)
      newDescriptions = append(newDescriptions,
        currentDescriptions[currentIndex])
    }
  }
  jsonItems, _ := json.Marshal(append(oldItems, newItems...))
  err = ioutil.WriteFile(dir + "/" + ".trending-push", jsonItems, 0666)
  check(err)

  for index, value := range newItems {
    message := "https://github.com/" + value + "\n" + 
      newDescriptions[index]
    res, err := http.PostForm("https://" + pushbulletToken +
      "@api.pushbullet.com/api/pushes",
      url.Values{
        "device_iden": {pushbulletDevice}, 
        "type": {"note"},
        "title": {value},
        "body": {message},
      })
    check(err)
    fmt.Print(res)
    fmt.Print("\n\n")
  }
}
