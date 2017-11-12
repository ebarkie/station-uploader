// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// Weather station uploader.
package main

//go:generate ./version.sh

import (
	"flag"
	"os"
	"strings"
	"time"
)

const (
	logInterval    = 30 * time.Minute
	minObsInterval = 2 * time.Second
)

type upChan chan obs

// upload is the main upload controller.
func upload(c config) {
	obss := make(chan obs)

	// Create channels and goroutines for configured uploaders.
	upChans := []upChan{}
	for _, up := range c.Uploaders {
		if u, exists := uploaders[strings.ToUpper(up.Type)]; exists {
			uc := make(upChan)
			go u.Upload(c.Station, up, uc)
			Info.Printf("%s(%s) uploader started", up.Name, strings.ToUpper(up.Type))
			upChans = append(upChans, uc)
		} else {
			Warn.Printf("Uploader type %s is unknown, ignoring", up.Type)
		}
	}

	// Stream observations from the station and send to the observations
	// channel.
	go func(obss chan obs) {
		for {
			err := streamEvents(c.Station.Host, obss) // Block, unless we hit an error
			if err != nil {
				Warn.Printf("Error retrieving observations: %s", err.Error())
			}
			time.Sleep(minObsInterval)
		}
	}(obss)

	// Read the observations channel, add tracked wind gust information, and
	// send to each configured uploader.
	go func(obss chan obs, upChans []upChan) {
		for {
			o := <-obss
			Debug.Printf("Received observation: %+v", o)

			for _, uc := range upChans {
				select {
				case uc <- o:
				default:
					// Skip due to prior upload being slow.  It might be nice
					// to collect this in the stats but that requires some
					// refactoring.
				}
			}
		}
	}(obss, upChans)

	// Block forever
	select {}
}

func main() {
	conf := flag.String("conf", "station-uploader.yaml", "station uploader config file")
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	if *debug {
		Debug.SetOutput(os.Stdout)
	}

	Info.Printf("Weather station uploader (version %s)", version)

	c := readConfig(*conf)
	upload(c)
}
