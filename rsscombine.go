package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/viper"
	"mvdan.cc/xurls"
)

func getUrlsFromFeedsUrl(feeds_url string) []string {
	log.Printf("Loading feed URLs from: %v", feeds_url)
	client := &http.Client{
		Timeout: time.Duration(viper.GetInt("client_timeout_seconds")) * time.Second,
	}
	response, err := client.Get(feeds_url)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			stringContents := string(contents)
			// TODO: this is a hack
			for _, exclude := range viper.GetStringSlice("feed_exclude_prefixes") {
				stringContents = strings.Replace(stringContents, exclude, "", -1)
			}
			feed_urls := xurls.Strict.FindAllString(stringContents, -1)
			return feed_urls
		}
	}
	return nil
}

func getUrls() []string {
	feeds_url := viper.GetString("feed_urls")
	if feeds_url != "" {
		return getUrlsFromFeedsUrl(feeds_url)
	}
	return viper.GetStringSlice("feeds")
}

func fetchUrl(url string, ch chan<- *gofeed.Feed) {
	log.Printf("Fetching URL: %v\n", url)
	fp := gofeed.NewParser()
	fp.Client = &http.Client{
		Timeout: time.Duration(viper.GetInt("client_timeout_seconds")) * time.Second,
	}
	feed, err := fp.ParseURL(url)
	if err == nil {
		ch <- feed
	} else {
		log.Printf("Error on URL: %v (%v)", url, err)
		ch <- nil
	}
}

func fetchUrls(urls []string) []*gofeed.Feed {
	allFeeds := make([]*gofeed.Feed, 0)
	ch := make(chan *gofeed.Feed)
	for _, url := range urls {
		go fetchUrl(url, ch)
	}
	for range urls {
		feed := <-ch
		if feed != nil {
			allFeeds = append(allFeeds, feed)
		}
	}
	return allFeeds
}

// TODO: there must be a shorter syntax for this
type byPublished []*gofeed.Feed

func (s byPublished) Len() int {
	return len(s)
}

func (s byPublished) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byPublished) Less(i, j int) bool {
	date1 := s[i].Items[0].PublishedParsed
	if date1 == nil {
		date1 = s[i].Items[0].UpdatedParsed
	}
	date2 := s[j].Items[0].PublishedParsed
	if date2 == nil {
		date2 = s[j].Items[0].UpdatedParsed
	}
	return date1.Before(*date2)
}

func getAuthor(feed *gofeed.Feed) string {
	if feed.Author != nil {
		return feed.Author.Name
	}
	if feed.Items[0].Author != nil {
		return feed.Items[0].Author.Name
	}
	log.Printf("Could not determine author for %v", feed.Link)
	return viper.GetString("default_author_name")
}

func combineallFeeds(allFeeds []*gofeed.Feed) *feeds.Feed {
	feed := &feeds.Feed{
		Title:       viper.GetString("title"),
		Link:        &feeds.Link{Href: viper.GetString("link")},
		Description: viper.GetString("description"),
		Author: &feeds.Author{
			Name:  viper.GetString("author_name"),
			Email: viper.GetString("author_email"),
		},
		Created: time.Now(),
	}
	sort.Sort(sort.Reverse(byPublished(allFeeds)))
	limit_per_feed := viper.GetInt("feed_limit_per_feed")
	seen := make(map[string]bool)
	for _, sourceFeed := range allFeeds {
		for _, item := range sourceFeed.Items[:limit_per_feed] {
			if seen[item.Link] {
				continue
			}
			created := item.PublishedParsed
			if created == nil {
				created = item.UpdatedParsed
			}
			feed.Items = append(feed.Items, &feeds.Item{
				Title:       item.Title,
				Link:        &feeds.Link{Href: item.Link},
				Description: item.Description,
				Author:      &feeds.Author{Name: getAuthor(sourceFeed)},
				Created:     *created,
				Content:     item.Content,
			})
			seen[item.Link] = true
		}
	}
	return feed
}

func GetAtomFeed() *feeds.Feed {
	urls := getUrls()
	allFeeds := fetchUrls(urls)
	combinedFeed := combineallFeeds(allFeeds)
	return combinedFeed
}

func LoadConfig() {
	viper.SetConfigName("rsscombine")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("RSSCOMBINE")
	viper.AutomaticEnv()
	viper.SetDefault("default_author_name", "Unknown Author")
	viper.SetDefault("client_timeout_seconds", "60")
	viper.SetDefault("feed_limit_per_feed", "20")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func main() {
	LoadConfig()
	bucket := viper.GetString("s3_bucket")
	filename := viper.GetString("s3_filename")
	combinedFeed := GetAtomFeed()
	atom, _ := combinedFeed.ToAtom()
	log.Printf("Rendered RSS with %v items", len(combinedFeed.Items))
	// if no S3 bucket is defined, simply print the feed on standard output
	if len(bucket) == 0 {
		fmt.Print(atom)
		return
	}
	// Upload the feed to S3
	sess, err := session.NewSession(&aws.Config{})
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filename),
		Body:        strings.NewReader(atom),
		ContentType: aws.String("text/xml"),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		log.Fatalf("Unable to upload %q to %q, %v", filename, bucket, err)
	}
	log.Printf("Successfully uploaded %q to %q\n", filename, bucket)
}
