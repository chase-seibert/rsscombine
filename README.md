# RSS Combine

Load RSS URLs from a hosted plain text file and combine them into one RSS feed.

## Quick Start

Here is how you run the service locally.

```bash
git clone git@github.com:chase-seibert/rsscombine.git
cd rsscombine
govendor sync
go run rsscombine.go
open http://localhost:8080
```

## Configuration

You can load the RSS URLs to combine via several methods.

## Local YAML File

You can create a local `rsscombine.yml` file in this format:

```yaml
title: My Technical RSS Feed
description: This is a personal collection of technical RSS feeds.
author_name: John Doe
author_email: john@example.com
port: 8080
cache_timeout_seconds: 3600
feeds:
  - http://feeds.feedburner.com/TechCrunch
  - http://feeds.arstechnica.com/arstechnica/technology-lab
  - http://www.reddit.com/r/technology/.rss
```

## Environment Variables

The above configuration value can also be specified as environment variables.
Each environment variable name should be in all caps, prefixed by `RSSCOMBINE_`
with underscored included.
For example, `title` can be loaded from `RSSCOMBINE_TITLE`, and
`cache_timeout_seconds` can be loaded from `RSSCOMBINE_CACHE_TIMEOUT_SECONDS`.

The only exception is the `feeds` items. For those values, you can specify
a `RSSCOMBINE_FEEDS_URL`. See bellow.

### Feeds URL

You can create a public file on the web, and RSS Combine can query that file and
parse out the URLs. This is especially useful for GitHub README files.

*Note: the file format does not matter, RSS Combine will pull any URL it can
find in the file.*

Example `README.md`:

```
This is a README with some URLs.

- TechCrunch http://feeds.feedburner.com/TechCrunch
- Ars Technica http://feeds.arstechnica.com/arstechnica/technology-lab
- Reddit http://www.reddit.com/r/technology/.rss
```

If that file is hosted at
`https://raw.githubusercontent.com/chase-seibert/rsscombine/master/examples/basic.md`, then you can
have RSS Combine load the file by defining the YAML key `feeds_url` or the
environment variable `RSSCOMBINE_FEEDS_URL` with that URL as the value. 

## Production

The project contains an example `Procfile` for Heroku.

```bash
heroku create
git push heroku master
heroku open
```
