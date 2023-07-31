package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
	Config rssole.ConfigSection `json:"config"`
}

func getFeedsFileConfigSection(filename string) (rssole.ConfigSection, error) {
	var cfgFile configFile

	jsonFile, err := os.Open(filename)
	if err != nil {
		return cfgFile.Config, fmt.Errorf("error opening file: %w", err)
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)

	err = decoder.Decode(&cfgFile)
	if err != nil {
		return cfgFile.Config, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return cfgFile.Config, nil
}

func main() {
	var configFilename, configReadCacheFilename string

	flag.StringVar(&configFilename, "c", "feeds.json", "config filename")
	flag.StringVar(&configReadCacheFilename, "r", "readcache.json", "readcache location")
	gokrazyFlag := flag.Bool("gokrazy", false, "use this if you are using gokrazy")

	flag.Parse()

	if *gokrazyFlag {
		// copy required files to /perm/home/rssole/
		errFeeds := cp.Copy("/etc/rssole/feeds.json", "/perm/home/rssole/feeds.json")
		errReadCache := cp.Copy("/etc/rssole/readcache.json", "/perm/home/rssole/readcache.json")

		// If there is an error we should exit
		if errFeeds != nil {
			log.Fatal(errFeeds)
		} else if errReadCache != nil {
			log.Fatal(errReadCache)
		}
	}

	cfg, err := getFeedsFileConfigSection(configFilename)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Listen == "" {
		cfg.Listen = defaultListenAddress
	}

	if cfg.UpdateSeconds == 0 {
		cfg.UpdateSeconds = defaultUpdateTimeSeconds
	}

	err = rssole.Start(configFilename, configReadCacheFilename, cfg.Listen, time.Duration(cfg.UpdateSeconds)*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}
