![badge](./badge.svg) ![workflow status](https://github.com/TheMightyGit/rssole/actions/workflows/build.yml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/TheMightyGit/rssole)](https://goreportcard.com/report/github.com/TheMightyGit/rssole)

# rssole

An absolutely no frills RSS Reader inspired by the late Google Reader. Runs on
your local machine or local network serving your RSS feeds via a clean
responsive web interface.

![Screenshot 2023-08-10 at 14 21 53](https://github.com/TheMightyGit/rssole/assets/888751/a44ae604-72a4-4e92-8ed7-5580663eaf0c)

A single executable with a single config file that can largely be configured
within the web UI.

Its greatest feature is the lack of excess features. It tries to do a simple
job well and not get in the way.

## Background

I really miss Google Reader, and I really like simplicity. So I made this
non-SaaS ode to Google Reader so I can triage my incoming information in one
place with one interface in a way I like. At heart this is a very self serving
project solely based around my needs, and because of that it's something I use
constantly. Hopefully it's of use to some other people, or you can build upon
it (MIT license, do what you want to it - make it comfortable for you).

## Advantages/Limitations

I see these as advantages (so they are unlikely to be added as features), but
some may see them as limitations...

- Only shows what's in the feed currently, does not store stories beyond their
  lifetime in the feed.
- Doesn't try to fetch anything from the linked page, only shows info present
  in the feed. The aim is not to keep you inside the RRS reader, if you want
  more then follow the link to the origin site.
- No bookmarks/favourites - you can already do this in the browser.
- It's not multi-user, there is no login or security protection. It's not
  intended as a SaaS product, it's just for you on your local machine or
  network. But you can stick an authenticating HTTP proxy in front of it if you
  wish.

## Pre-Built Binaries and Packages

Check out the [Releases](https://github.com/TheMightyGit/rssole/releases/)
section in github, there should be a good selection of pre-built binaries
and packages for various platforms.

## Installing via Brew

```console
$ brew install themightygit/rssole/rssole
```

## Installing via Go

You can install the binary with go install:

```console
$ go install github.com/TheMightyGit/rssole/cmd/rssole@latest
```

## Building

NOTE: You can ignore the `Makefile`, that's really just a helper for me during
development.

To build for your local architecture/OS...

```console
$ go build ./cmd/...
```

It should also cross build for all the usual golang targets fine as well (as no
CGO is used)...

```console
$ GOOS=linux GOARCH=amd64 go build ./cmd/...
$ GOOS=linux GOARCH=arm64 go build ./cmd/...
$ GOOS=darwin GOARCH=amd64 go build ./cmd/...
$ GOOS=darwin GOARCH=arm64 go build ./cmd/...
$ GOOS=windows GOARCH=amd64 go build ./cmd/...
$ GOOS=windows GOARCH=arm64 go build ./cmd/...
```

...but I only regularly test on `darwin/amd64` and `linux/amd64`.
I've seen it run on `windows/amd64`, but it's not something I try regularly.

### Smallest Binary

Go binaries can be a tad chunky, so if you're really space constrained then...

```console
$ go build -ldflags "-s -w" ./cmd/...
$ upx rssole
```

## Running

### Command Line

If you built locally then it should be in the current directory:

```console
$ ./rssole
```

If you used `go install` or brew then it should be on your path already:

```console
$ rssole
```

### GUI

Double click on the file, I guess.

If your system has restrictions on which binaries it will run then try
compiling locally instead of using the pre-built binaries.

## Now read your feeds with your browser

Now open your browser on `<hostname/ip>:8090` e.g. http://localhost:8090

## Network Options

By default it binds to `0.0.0.0:8090`, so it will be available on all network
adaptors on your host. You can change this in the `rssole.json` config file.

I run rssole within a private network so this is good enough for me so that I
can run it once but access it from all my devices. If you run this on an alien
network then someone else can mess with the UI (there's no protection at all on
it) - change the `listen` value in `rssole.json` to `127.0.0.1:8090` if you
only want it to serve locally.

If you want to protect rssole behind a username and password or encryption
(because you want rssole wide open on the net so you can use it from anywhere)
then you'll need a web proxy that can be configured to sit in front of it to
provide that protection. I'm highly unlikely to add username/password or
encryption directly to rssole as I don't need it. Maybe someone will create a
docker image that autoconfigures all of that... maybe that someone is you?

## Config

### Arguments

```console
$ ./rssole -h
Usage of ./rssole:
  -c string
        config filename (default "rssole.json")
  -r string
        readcache location (default "rssole_readcache.json")
```

### `rssole.json`

There are two types of feed definition...

- Regular RSS URLs.
- Scrape from website (for those pesky sites that have no RSS feed).
  - Scraping uses css selectors and is not well documented yet.

Use `category` to group similar feeds together.

```json
{
  "config": {
    "listen": "0.0.0.0:8090",
    "update_seconds": 300
  },
  "feeds": [
    {"url":"https://github.com/TheMightyGit/rssole/releases.atom", "category":"Github Releases"},
    {"url":"https://news.ycombinator.com/rss", "category":"Nerd"},
    {"url":"http://feeds.bbci.co.uk/news/rss.xml", "category":"News"},
    {
      "url":"https://www.pcgamer.com/uk/news/", "category":"Games",
      "name":"PCGamer News",
      "scrape": {
        "urls": [
          "https://www.pcgamer.com/uk/news/",
          "https://www.pcgamer.com/uk/news/page/2/",
          "https://www.pcgamer.com/uk/news/page/3/"
        ],
        "item": ".listingResult",
        "title": ".article-name",
        "link": ".article-link"
      }
    }
  ]
}
```

## Usage with gokrazy

Add the forked repo to the gokrazy instace

``` bash
$ gok -i brun0-pi add ~/Desktop/personal/rssole/cmd/rssole/
```

Then use this in the gokrazy Config

``` json
  "github.com/TheMightyGit/rssole/cmd/rssole": {
    "CommandLineFlags": [
      "-gokrazy",
      "-c=/perm/home/rssole/feeds.json",
      "-r=/perm/home/rssole/readcache.json"
    ],
    "ExtraFilePaths": {
      "/etc/rssole/feeds.json": "/home/brun0/gokrazy/brun0-pi/rssole/feeds.json",
      "/etc/rssole/readcache.json": "/home/brun0/gokrazy/brun0-pi/rssole/readcache.json"
    }
  }
```

After this just update the gokrazy instance and that should word out of the box

``` bash
gok update -i brun0-pi
```

## Key Dependencies

I haven't had to implement anything actually difficult, I just do a bit of
plumbing. All the difficult stuff has been done for me by these projects...

- github.com/mmcdole/gofeed - for reading all sorts of RSS formats.
- github.com/andybalholm/cascadia - for css selectors during website scrapes.
- github.com/JohannesKaufmann/html-to-markdown/v2 to convert HTML into Markdown
  (thus sanitizing and simplifying it).
- github.com/gomarkdown/markdown to render content markdown back to HTML.
- github.com/k3a/html2text - for making a plain text summary of html.
- HTMX - for the javascript anti-framework (and a backend engineers delight).
- Bootstrap 5 - for HTML niceness simply because I know it slightly better than
  the alternatives.
