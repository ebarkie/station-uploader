// Copyright (c) 2016-2018 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

// Uploader is an interface for uploading observation data to various
// places like ham radio APRS, Internet APRS, or various Internet REST
// services.
type Uploader interface {
	Upload(ConfigStation, ConfigUploader, upChan)
}

// Registered uploaders
var uploaders = map[string]Uploader{
	"APRS":    APRSUploader{},
	"APRS-IS": APRSUploader{},
	"NULL":    NullUploader{},
	"WC":      WCUploader{},
	"WU":      WUUploader{},
}
