package main

import (
  "golang.org/x/net/html"
  "net/http"
  "os"
  "log"
  "os/exec"
  "bufio"
  "strings"
)

func main() {
    urls := crawl()
    crawlEach(urls)
}
func crawl() []string {
  resp, err := http.Get("https://golang.org/src/" + os.Args[1])
  if err != nil {
    log.Fatal("could not get package files")
  }

  defer resp.Body.Close()
  tizer := html.NewTokenizer(resp.Body)
  var urls []string
  loop:
  for {
    token := tizer.Next()
    switch {
    case token == html.ErrorToken:
      break loop
    case token == html.StartTagToken:
      found := tizer.Token()
      if found.Data != "a" {
        continue
      }
      for _, a := range found.Attr {
        if a.Key == "href" {
          urls = append(urls, a.Val)
        }
      }
    }
  }
  return urls
}
func crawlEach(urls []string) {
  for _ ,url := range urls {
    resp, err := http.Get("https://golang.org/src/" + os.Args[1] + "/" + url)
    if err != nil {
      log.Fatal("could not get package file")
    }
    defer resp.Body.Close()

    scanner := bufio.NewScanner(resp.Body)
    for scanner.Scan() {
      line := scanner.Text()
      if strings.Contains(line, os.Args[2]) {
        tizer := html.NewTokenizer(strings.NewReader(line))
        token := tizer.Next()

        switch {
        case token == html.ErrorToken:
          break
        case token == html.StartTagToken:
          found := tizer.Token()
          if found.Data != "span" {
            continue
          }

          for _, a := range found.Attr {
            if a.Key == "id" {
              cmd := exec.Command("open", (resp.Request.URL).String() + "#" + a.Val)
              cmd.Run()
            }
          }
        }
      }
    }
  }
}
