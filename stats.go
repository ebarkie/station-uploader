// Copyright (c) 2016-2018 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import "time"

// stats sets up a new stat collector for an uplaoder.
func stats(name string) (ok, sk, er chan int) {
	ok = make(chan int)
	sk = make(chan int)
	er = make(chan int)

	go func() {
		var oks, skips, errors int

		t := time.NewTimer(logInterval)
		for {
			select {
			case i := <-ok:
				oks += i
			case i := <-sk:
				skips += i
			case i := <-er:
				errors += i
			case <-t.C:
				Info.Printf("%s upload stats ok=%d/skips=%d/errors=%d",
					name, oks, skips, errors)
				oks, skips, errors = 0, 0, 0
				t.Reset(logInterval)
			}
		}

	}()

	return
}
