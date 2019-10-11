// Copyright (c) 2016 Eric Barkie. All rights reserved.
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

	"gitlab.com/ebarkie/weatherlink/data"
)

// Errors.
var (
	errBadResponse = errors.New("HTTP request returned non-OK status code")
)

type obs struct {
	Archive data.Archive
	Loop    data.Loop `json:"loop"`
}

// getLastArchive makes a HTTP GET call to the station at addr and retreives
// the most recent archive record no older than begin.
func getLastArchive(addr string, begin time.Time) (a data.Archive, err error) {
	var resp *http.Response
	resp, err = http.Get("http://" + addr + "/archive?begin=" + begin.Format(time.RFC3339))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errBadResponse
		return
	}

	var as []data.Archive
	err = json.NewDecoder(resp.Body).Decode(&as)
	if err != nil && len(as) > 0 {
		a = as[0]
	}

	return
}

// streamEvents makes a HTTP GET call to the station at addr and receives
// a continuous stream of archive and loop events using Server-sent events
// format.
func streamEvents(addr string, obss chan<- obs) (err error) {
	// Setup and initiate request
	var resp *http.Response
	resp, err = http.Get("http://" + addr + "/events")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Check HTTP status.
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP request returned non-OK status code: %d", resp.StatusCode)
		return
	}

	// Lightweight Server-sent event parser.
	reader := bufio.NewReader(resp.Body)

	// Archive records can take a while to arrive so initially seed with the
	// most recent one if it's not older than 9 minutes.  Ignore the error
	// response; if this fails it will come through the event stream later.
	var o obs
	o.Archive, _ = getLastArchive(addr, time.Now().Add(-9*time.Minute))

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
				err = json.Unmarshal([]byte(d), &o.Loop)
				obss <- o
			}
			if err != nil {
				return
			}
		}
	}
}
