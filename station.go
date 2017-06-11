// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ebarkie/weatherlink"
)

// Errors.
var (
	ErrBadResponse = errors.New("HTTP request returned non-OK status code")
)

type obs struct {
	Archive weatherlink.Archive
	Loop    weatherlink.Loop `json:"loop"`
}

// getLastArchive makes a HTTP GET call to the Davis station at serverAddress and
// retreives the most recent archive record no older than begin.
func getLastArchive(serverAddress string, begin time.Time) (a weatherlink.Archive, err error) {
	var resp *http.Response
	resp, err = http.Get("http://" + serverAddress + "/archive?begin=" + begin.Format(time.RFC3339))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = ErrBadResponse
		return
	}

	var as []weatherlink.Archive
	err = json.NewDecoder(resp.Body).Decode(&as)
	if (err != nil) && (len(as) > 0) {
		a = as[0]
	}

	return
}

// streamEvents makes a HTTP GET call to the Davis station at serverAddress and
// receives a continuous stream of archive and loop events using Server-sent events
// format.
func streamEvents(serverAddress string, obss chan<- obs) (err error) {
	// Setup and initiate request
	var resp *http.Response
	resp, err = http.Get("http://" + serverAddress + "/events")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP request returned non-OK status code: %d", resp.StatusCode)
		return
	}

	// Lightweight Server-sent event parser.
	reader := bufio.NewReader(resp.Body)

	// Archive records can take a while to arrive so initially seed it with the
	// most recent one if it's not older than 9 minutes.  It doesn't matter if it
	// errors since we'll just get it via the normal event stream later.
	var o obs
	o.Archive, _ = getLastArchive(serverAddress, time.Now().Add(-9*time.Minute))

	var e, s string
	for {
		s, err = reader.ReadString('\n')
		if err != nil {
			return
		}

		if strings.HasPrefix(s, "event:") {
			e = strings.TrimSpace(strings.TrimPrefix(s, "event:"))
		} else if strings.HasPrefix(s, "data:") {
			d := strings.TrimPrefix(s, "data:")
			switch e {
			case "archive":
				err = json.Unmarshal([]byte(d), &o.Archive)
			case "loop":
				err = json.Unmarshal([]byte(d), &o)
				obss <- o
			}
			if err != nil {
				return
			}
		}
	}
}
