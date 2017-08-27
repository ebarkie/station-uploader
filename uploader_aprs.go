// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/ebarkie/aprs"
)

// APRSUploader handles APRS(-RF) and APRS-IS uploads.
type APRSUploader struct{}

// Upload sends the received observations to APRS at the specified
// interval.
func (APRSUploader) Upload(station ConfigStation, up ConfigUploader, uc upChan) {
	w := aprs.Wx{
		Lat:  station.Lat,
		Lon:  station.Lon,
		Type: station.Type,
	}

	ok, _, er := stats(up.Name)
	t := time.NewTimer(0)
	for {
		o := <-uc

		// Only upload if interval has passed
		select {
		case <-t.C:
		default:
			continue
		}

		// Build APRS text payload
		w.Clear()
		w.Altimeter(o.Loop.Bar.Altimeter)
		w.Humidity(o.Loop.OutHumidity)
		w.RainRate(o.Loop.Rain.Rate)
		w.RainLast24Hours(o.Loop.Rain.Accum.Last24Hours)
		w.RainToday(o.Loop.Rain.Accum.Today)
		w.SolarRadiation(o.Loop.SolarRad)
		w.Temperature(o.Loop.OutTemp)
		if o.Archive.WindSpeedAvg > 0 {
			w.WindDirection(o.Loop.Wind.Cur.Dir)
		} else {
			w.WindDirection(360)
		}
		w.WindSpeed(o.Archive.WindSpeedAvg)
		w.WindGust(o.Archive.WindSpeedHi)

		// Upload
		a := aprs.Frame{Text: w.String()}
		a.Src.FromString(up.ID)
		dial := func() error { return nil }
		switch strings.ToUpper(up.Type) {
		case "APRS":
			a.Dst = aprs.Address{Call: "APZ001"} // Experimental v0.0.1
			//a.Path = aprs.Path{aprs.Address{Call: "WIDE1", SSID: 1}, aprs.Address{Call: "WIDE2", SSID: 1}}
			dial = func() error { return a.SendKISS(up.Dial) }
		case "APRS-IS":
			a.Dst = aprs.Address{Call: "APRS"}
			a.Path = aprs.Path{aprs.Address{Call: "TCPIP", Repeated: true}}
			pw, err := strconv.Atoi(up.Password)
			if err != nil {
				pw = -1
			}
			dial = func() error { return a.SendIS(up.Dial, pw) }
		}
		Debug.Printf("%s frame: %s", up.Name, a)

		err := dial()
		if err != nil {
			Error.Printf("%s upload error: %s", up.Name, err.Error())
			er <- 1
		} else {
			Debug.Printf("%s upload successful", up.Name)
			ok <- 1
		}

		t.Reset(up.Interval * time.Second)
	}
}
