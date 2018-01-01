// Copyright (c) 2016-2018 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/ebarkie/wunderground"
)

// WUUploader handles Weather Underground uploads.
type WUUploader struct{}

// Upload sends the received observations to Weather Underground at the
// specified interval.  If an interval of 0 is used it will use RapidFire
// mode and send as rapidly as possible (usually every 2 seconds).
func (WUUploader) Upload(station ConfigStation, up ConfigUploader, uc upChan) {
	w := wunderground.Pws{ID: up.ID, Password: up.Password}
	w.SoftwareType = "GoWunder 1337." + version
	w.Interval = up.Interval * time.Second

	wx := &wunderground.Wx{}

	ok, sk, er := stats(up.Name)
	t := time.NewTimer(0)
	for {
		o := <-uc

		// Only upload if interval has passed.
		select {
		case <-t.C:
		default:
			continue
		}

		// Build Wunderground payload.
		wx.Bar(o.Loop.Bar.SeaLevel)
		wx.DailyRain(o.Loop.Rain.Accum.Today)
		wx.DewPoint(o.Loop.DewPoint)
		wx.OutHumidity(o.Loop.OutHumidity)
		wx.OutTemp(o.Loop.OutTemp)
		wx.RainRate(o.Loop.Rain.Rate)
		for _, v := range o.Loop.SoilMoist {
			if v != nil {
				wx.SoilMoist(*v)
			}
		}
		for _, v := range o.Loop.SoilTemp {
			if v != nil {
				wx.SoilTemp(float64(*v))
			}
		}
		wx.SolarRad(o.Loop.SolarRad)
		wx.UVIndex(o.Loop.UVIndex)
		if o.Loop.Wind.Cur.Speed > 0 {
			wx.WindDir(o.Loop.Wind.Cur.Dir)
		}
		wx.WindSpeed(float64(o.Loop.Wind.Cur.Speed))
		if o.Archive.WindSpeedHi > o.Loop.Wind.Cur.Speed {
			wx.WindGustDir(o.Archive.WindDirHi)
			wx.WindGustSpeed(float64(o.Archive.WindSpeedHi))
		}
		if o.Loop.Wind.Gust.Last10MinSpeed > 0 {
			wx.WindGustDir10m(o.Loop.Wind.Gust.Last10MinDir)
		}
		wx.WindGustSpeed10m(o.Loop.Wind.Gust.Last10MinSpeed)
		wx.WindSpeedAvg2m(o.Loop.Wind.Avg.Last2MinSpeed)

		// Upload.
		Debug.Printf("%s request URL: %s", up.Name, w.Encode(wx))
		err := w.Upload(wx)
		if err != nil {
			Error.Printf("%s upload error: %s", up.Name, err.Error())
			er <- 1
		} else if w.Skipped() {
			Debug.Printf("%s upload was unnecessary and skipped", up.Name)
			sk <- 1
		} else {
			Debug.Printf("%s upload successful", up.Name)
			ok <- 1
		}

		t.Reset(up.Interval * time.Second)
	}
}
