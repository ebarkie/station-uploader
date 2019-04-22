// Copyright (c) 2019 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/ebarkie/windy"
)

// WindyUploader handles Windy Stations uploads.
type WindyUploader struct{}

// Upload sends the received observations to Windy.com at the specified
// interval or every 5 minutes, whichever is longer.
func (WindyUploader) Upload(station ConfigStation, up ConfigUploader, uc upChan) {
	// Upload interval can not be shorter than 5 minutes.
	interval := up.Interval
	if interval < 300 {
		interval = 300
	}

	s := windy.Station{ID: up.ID, Key: up.Password}

	wx := &windy.Wx{}

	ok, _, er := stats(up.Name)
	t := time.NewTimer(0)
	for {
		o := <-uc

		// Only upload if interval has passed.
		select {
		case <-t.C:
		default:
			continue
		}

		// Build Weathercloud payload.
		wx.Bar(o.Loop.Bar.SeaLevel)
		wx.DewPoint(o.Loop.DewPoint)
		wx.OutHumidity(o.Loop.OutHumidity)
		wx.OutTemp(o.Loop.OutTemp)
		wx.RainRate(o.Loop.Rain.Rate)
		wx.UVIndex(o.Loop.UVIndex)
		if o.Loop.Wind.Cur.Speed > 0 {
			wx.WindDir(o.Loop.Wind.Cur.Dir)
		}
		wx.WindSpeed(float64(o.Loop.Wind.Cur.Speed))
		wx.WindGustSpeed(o.Loop.Wind.Gust.Last10MinSpeed)

		// Upload.
		Debug.Printf("%s request URL: %s", up.Name, s.Encode(wx))
		err := s.Upload(wx)
		if err != nil {
			Error.Printf("%s upload error: %s", up.Name, err.Error())
			er <- 1
		} else {
			Debug.Printf("%s upload successful", up.Name)
			ok <- 1
		}

		t.Reset(interval * time.Second)
	}
}
