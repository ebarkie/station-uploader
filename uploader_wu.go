// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
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
	w := wunderground.New(up.ID, up.Password)
	w.SoftwareType = "GoWunder 1337." + version
	w.Interval = up.Interval * time.Second

	ok, sk, er := stats(up.Name)
	t := time.NewTimer(0)
	for {
		o := <-uc

		// Only upload if interval has passed
		select {
		case <-t.C:
		default:
			continue
		}

		// Build Wunderground payload
		w.Wx.Barometer(o.Loop.Bar.SeaLevel)
		w.Wx.DailyRain(o.Loop.Rain.Accum.Today)
		w.Wx.DewPoint(o.Loop.DewPoint)
		w.Wx.OutdoorHumidity(o.Loop.OutHumidity)
		w.Wx.OutdoorTemperature(o.Loop.OutTemp)
		w.Wx.RainRate(o.Loop.Rain.Rate)
		for _, v := range o.Loop.SoilMoist {
			if v != nil {
				w.Wx.SoilMoisture(*v)
			}
		}
		for _, v := range o.Loop.SoilTemp {
			if v != nil {
				w.Wx.SoilTemperature(float64(*v))
			}
		}
		w.Wx.SolarRadiation(o.Loop.SolarRad)
		w.Wx.UVIndex(o.Loop.UVIndex)
		if o.Loop.Wind.Cur.Speed > 0 {
			w.Wx.WindDirection(o.Loop.Wind.Cur.Dir)
		}
		w.Wx.WindSpeed(float64(o.Loop.Wind.Cur.Speed))
		if o.Archive.WindSpeedHi > o.Loop.Wind.Cur.Speed {
			w.Wx.WindGustDirection(o.Archive.WindDirHi)
			w.Wx.WindGustSpeed(float64(o.Archive.WindSpeedHi))
		}
		if o.Loop.Wind.Gust.Last10MinSpeed > 0 {
			w.Wx.WindGustDirection10m(o.Loop.Wind.Gust.Last10MinDir)
		}
		w.Wx.WindGustSpeed10m(o.Loop.Wind.Gust.Last10MinSpeed)
		w.Wx.WindSpeedAverage2m(o.Loop.Wind.Avg.Last2MinSpeed)

		// Upload
		Debug.Printf("%s request URL: %s", up.Name, w.String())
		skipped, err := w.Upload()
		if err != nil {
			Error.Printf("%s upload error: %s", up.Name, err.Error())
			er <- 1
		} else if skipped {
			Debug.Printf("%s upload was unnecessary and skipped", up.Name)
			sk <- 1
		} else {
			Debug.Printf("%s upload successful", up.Name)
			ok <- 1
		}

		t.Reset(up.Interval * time.Second)
	}
}
