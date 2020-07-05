// Copyright (c) 2019 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"strconv"
	"time"

	"github.com/ebarkie/weatherlink/units"
	"github.com/ebarkie/windy"
)

// WindyUploader handles Windy Stations uploads.
type WindyUploader struct{}

// Upload sends the received observations to Windy at the specified
// interval or every 5 minutes, whichever is longer.
func (WindyUploader) Upload(station ConfigStation, up ConfigUploader, uc upChan) {
	// Upload interval can not be shorter than 5 minutes.
	interval := up.Interval
	if interval < 300 {
		interval = 300
	}

	r := &windy.Req{Key: up.Password}
	id, err := strconv.ParseInt(up.ID, 10, 32)
	if err != nil {
		Error.Printf("%s has an invalid station ID: %s, skipping uploads", up.Name, up.ID)
		return
	}

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

		// Build Windy payload.
		obs := windy.Obs{
			Station: int32(id),
		}

		obs.BaromIn = &o.Loop.Bar.SeaLevel
		dewpoint := units.Fahrenheit(o.Loop.DewPoint).Celsius()
		obs.DewPoint = &dewpoint
		rh := float64(o.Loop.OutHumidity)
		obs.RH = &rh
		obs.TempF = &o.Loop.OutTemp
		obs.RainIn = &o.Loop.Rain.Accum.LastHour
		obs.UV = &o.Loop.UVIndex
		if o.Archive.WindSpeedAvg > 0 {
			obs.WindDir = &o.Archive.WindDirPrevail
		}
		windSpeedMPH := float64(o.Archive.WindSpeedAvg)
		obs.WindSpeedMPH = &windSpeedMPH
		windGustMPH := float64(o.Archive.WindSpeedHi)
		obs.WindGustMPH = &windGustMPH

		r.Obss = []windy.Obs{obs}

		// Upload.
		Debug.Printf("%s request URL: %s, body: %s", up.Name, r.Encode(), r.Body())
		err := r.Upload()
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
