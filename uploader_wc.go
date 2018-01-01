// Copyright (c) 2016-2018 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/ebarkie/weathercloud"
)

// WCUploader handles Weathercloud uploads.
type WCUploader struct{}

// Upload sends the received observations to Weathercloud at the specified
// interval or every 10 minutes, whichever is longer.
func (WCUploader) Upload(station ConfigStation, up ConfigUploader, uc upChan) {
	const minInterval = 600

	w := weathercloud.Device{WID: up.ID, Key: up.Password}
	w.SoftwareVersion = version

	wx := &weathercloud.Wx{}

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
		wx.DailyET(o.Loop.ET.Today)
		wx.DailyRain(o.Loop.Rain.Accum.Today)
		wx.DewPoint(o.Loop.DewPoint)
		wx.HeatIndex(o.Loop.HeatIndex)
		wx.OutHumidity(o.Loop.OutHumidity)
		wx.OutTemp(o.Loop.OutTemp)
		wx.RainRate(o.Loop.Rain.Rate)
		for _, v := range o.Loop.SoilMoist {
			if v != nil {
				wx.SoilMoist(*v)
			}
		}
		wx.SolarRad(o.Loop.SolarRad)
		wx.UVIndex(o.Loop.UVIndex)
		wx.WindChill(o.Loop.WindChill)
		if o.Loop.Wind.Cur.Speed > 0 {
			wx.WindDir(o.Loop.Wind.Cur.Dir)
		}
		wx.WindDirAvg(o.Archive.WindDirPrevail)
		wx.WindSpeed(float64(o.Loop.Wind.Cur.Speed))
		wx.WindGustSpeed(o.Loop.Wind.Gust.Last10MinSpeed)
		wx.WindSpeedAvg(o.Loop.Wind.Avg.Last10MinSpeed)

		// Upload.
		Debug.Printf("%s request URL: %s", up.Name, w.Encode(wx))
		err := w.Upload(wx)
		if err != nil {
			Error.Printf("%s upload error: %s", up.Name, err.Error())
			er <- 1
		} else {
			Debug.Printf("%s upload successful", up.Name)
			ok <- 1
		}

		// Upload interval can not be shorter than 10 minutes.
		if up.Interval > minInterval {
			t.Reset(up.Interval * time.Second)
		} else {
			t.Reset(minInterval * time.Second)
		}
	}
}
