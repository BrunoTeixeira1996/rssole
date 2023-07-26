package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/TheMightyGit/rssole/internal/rssole"
	"github.com/u-root/uio/cp"
)

const (
	defaultListenAddress     = "0.0.0.0:8090"
	defaultUpdateTimeSeconds = 300
)

type configFile struct {
	Config configSection `json:"config"`
}

type configSection struct {
	Listen        string `json:"listen"`
	UpdateSeconds int    `json:"update_seconds"`
}

func getFeedsFileConfigSection(filename string) configSection {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer jsonFile.Close()

	var c configFile
	d := json.NewDecoder(jsonFile)
	err = d.Decode(&c)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	return c.Config
}

func runningInGoKrazy() (string, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// copy required files to /perm/home/rssole/
	errFeeds := cp.Copy("/etc/rssole/feeds.json", "/perm/home/rssole/feeds.json")
	errReadCache := cp.Copy("/etc/rssole/readcache.json", "/perm/home/rssole/readcache.json")
	if errFeeds != nil || errReadCache != nil {
		return "", err
	}

	return currentPath, nil
}

func main() {
	var configFilename, configReadCacheFilename, currentGoKrazyPath string
	var err error

	flag.StringVar(&configFilename, "c", "feeds.json", "config filename")
	flag.StringVar(&configReadCacheFilename, "r", "readcache.json", "readcache location")
	gokrazyFlag := flag.Bool("gokrazy", false, "use this if you are using gokrazy")

	flag.Parse()

	if *gokrazyFlag {
		currentGoKrazyPath, err = runningInGoKrazy()

		// If there is an error we should exit
		if err != nil {
			log.Fatal(err)
		}

		// now we have the correct locations of these two files in /perm/home/rssole
		configFilename = currentGoKrazyPath + "/" + configFilename
		configReadCacheFilename = currentGoKrazyPath + "/" + configReadCacheFilename
	}

	cfg := getFeedsFileConfigSection(configFilename)

	if cfg.Listen == "" {
		cfg.Listen = defaultListenAddress
	}
	if cfg.UpdateSeconds == 0 {
		cfg.UpdateSeconds = defaultUpdateTimeSeconds
	}

	rssole.Start(configFilename, configReadCacheFilename, cfg.Listen, time.Duration(cfg.UpdateSeconds)*time.Second)
}
