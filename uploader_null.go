// Copyright (c) 2016-2018 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import "time"

// NullUploader is a fake uploader used for testing.
type NullUploader struct{}

// Upload tests various functions like reading channels, logging, and
// saving stats but doesn't actually upload data anywhere.
func (NullUploader) Upload(_ ConfigStation, up ConfigUploader, uc upChan) {
	ok, _, _ := stats(up.Name)
	t := time.NewTimer(0)
	for {
		<-uc

		// Only upload if interval has passed.
		select {
		case <-t.C:
		default:
			continue
		}

		// Null upload.
		Info.Printf("%s upload fired", up.Name)
		ok <- 1

		t.Reset(up.Interval * time.Second)
	}

}
