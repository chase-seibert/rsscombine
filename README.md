# RSS Combine

Combine multiple RSS feeds into a single feed, as a service.

## Quick Start

Here is how you run the service locally.

```bash
git clone git@github.com:chase-seibert/rsscombine.git
cd rsscombine
govendor sync
go run cmd/rsscombine-server/main.go
open http://localhost:8080
```

## Configuration

You can specify configuration options as either a YAML file, or as environment
variables. Environment variable names should be in all caps, prefixed by
`RSSCOMBINE_` with underscored included. The only exception is the `feeds`
items. You cannot specify that list of URLs via environment variables. Instead,
you can specify a `RSSCOMBINE_FEEDS_URL`. See "Feeds URL" bellow.

### Options

See the "Example YAML File" section for example defaults.

| YAML Name             | Environment Variable             | Description                                                                           |
|-----------------------|----------------------------------|---------------------------------------------------------------------------------------|
| title                 | RSSCOMBINE_TITLE                 | Title of the new RSS feed.                                                            |
| link                  | RSSCOMBINE_LINK                  | Link to the new RSS feed. Can be a webpage or the feed URL.                           |
| description           | RSSCOMBINE_DESCRIPTION           | Description of your new feed, shows in RSS readers.                                   |
| author_name           | RSSCOMBINE_AUTHOR_NAME           | Your full name, shows in RSS readers.                                                 |
| author_email          | RSSCOMBINE_AUTHOR_EMAIL          | Your email, shows in RSS readers.                                                     |
| port                  | PORT, RSSCOMBINE_AUTHOR_PORT     | Port to run the service on. For Heroku support, PORT environment variable supersedes. |
| cache_timeout_seconds | RSSCOMBINE_CACHE_TIMEOUT_SECONDS | Seconds to cache individual feeds in memory, as well as a feeds_url file.             |
| server_timeout_seconds | RSSCOMBINE_SERVER_TIMEOUT_SECONDS | Seconds to timeout calls to the combined RSS feed sever.             |
| client_timeout_seconds | RSSCOMBINE_CLIENT_TIMEOUT_SECONDS | Seconds to timeout call from the server to the individual RSS feeds.             |
| feeds                 |                                  | List of feeds to combine. Cannot be specified via environment variable.               |
| feed_urls             | RSSCOMBINE_FEED_URLS             | Optional: URL to parse feed URLs from. If set, this overrides the feeds setting.      |
| feed_exclude_prefixes | RSSCOMBINE_FEED_EXCLUDE_PREFIXES | Optional: list of URL prefixes to exclude from feed_urls parsing.                     |
| s3_bucket             | RSSCOMBINE_S3_BUCKET             | Optional: bucket name to use when uploading to S3. |
| s3_filename           | RSSCOMBINE_S3_FILENAME           | Optional: file name to use when uploading to S3. |

## Example YAML File

You can create a local `rsscombine.yml` file in this format:

```yaml
title: My Technical RSS Feed
link: http://wherethisfeedishosted.com/feed
description: This is a personal collection of technical RSS feeds.
author_name: John Doe
author_email: john@example.com
port: 8080
cache_timeout_seconds: 3600
feeds:
  - http://feeds.feedburner.com/TechCrunch
  - http://feeds.arstechnica.com/arstechnica/technology-lab
  - http://www.reddit.com/r/technology/.rss
  - http://rss.slashdot.org/slashdot/slashdotMainatom
```

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

### Web Server

You can run a live RSS server. The project contains an example `Procfile`
for Heroku.

```bash
heroku create
git push heroku master
heroku open
```

### Generate Static RSS File in S3

You can also generate a static RSS file and upload to S3. Combined with the
[Heroku Scheduler](https://elements.heroku.com/addons/scheduler), this option
will minimize your dyno hours and serve the RSS much faster.

For S3 uploads, you need to set the following as environment variables.

| Environment Variable             | Description                                                                           |
|----------------------------------|---------------------------------------------------------------------------------------|
| AWS_REGION                       | The AWS Region your bucket is in. |
| AWS_ACCESS_KEY_ID                | The AWS access key ID for your bucket. |
| AWS_SECRET_ACCESS_KEY            | The AWS secret access key for your bucket. |

See the Configuration section for setting `s3_bucket` and `s3_filename`.
