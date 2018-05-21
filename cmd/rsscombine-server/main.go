package main

import "github.com/rsscombine"
import "fmt"
import "time"
import "net/http"
import "log"
import "github.com/patrickmn/go-cache"
import "github.com/spf13/viper"
import "os"
import "strconv"
import "github.com/gorilla/feeds"

var feedCache = cache.New(3600*time.Second, 3600*time.Second)

func handler(w http.ResponseWriter, r *http.Request) {
  var combinedFeed *feeds.Feed
  cached, found := feedCache.Get("combinedFeed")
  if found {
    combinedFeed = cached.(*feeds.Feed)
  } else {
    combinedFeed = rsscombine.GetAtomFeed()
  }
  atom, _ := combinedFeed.ToAtom()
  feedCache.Set("combinedFeed", combinedFeed, cache.DefaultExpiration)
  fmt.Fprintf(w, atom)
  log.Printf("Rendered RSS with %v items", len(combinedFeed.Items))
}

func nullHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
  rsscombine.LoadConfig()
  cache_timeout_seconds := time.Duration(viper.GetInt("cache_timeout_seconds")) * time.Second
  feedCache = cache.New(cache_timeout_seconds, cache_timeout_seconds)
  herokuPort := os.Getenv("PORT")
  port := 0
  if herokuPort == "" {
    port = viper.GetInt("port")
  } else {
    port, _ = strconv.Atoi(herokuPort)
  }
  http.HandleFunc("/favicon.ico", nullHandler)
  http.HandleFunc("/", handler)
  serverTimeoutSeconds := time.Duration(viper.GetInt("server_timeout_seconds"))
  srv := &http.Server{
    Addr: fmt.Sprintf(":%v", port),
    ReadTimeout: serverTimeoutSeconds * time.Second,
    WriteTimeout: serverTimeoutSeconds * time.Second,
  }
  log.Printf("Listening on: http://localhost:%v/\n", port)
  log.Fatal(srv.ListenAndServe())
}
