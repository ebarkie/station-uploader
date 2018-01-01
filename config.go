// Copyright (c) 2016-2018 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// config is the station uploader configuration.
type config struct {
	Station   ConfigStation    `yaml:"station"`
	Uploaders []ConfigUploader `yaml:"uploaders"`
}

// ConfigStation is the Davis Instruments weather server information.
type ConfigStation struct {
	Host string  `yaml:"host"`
	Lat  float64 `yaml:"lat"`
	Lon  float64 `yaml:"lon"`
	Type string  `yaml:"type"`
}

// ConfigUploader represents an APRS-IS, APRS-RF, or WU uploader.
type ConfigUploader struct {
	Name     string        `yaml:"name"`
	Type     string        `yaml:"type"`
	Interval time.Duration `yaml:"interval"`
	Dial     string        `yaml:"dial"`
	ID       string        `yaml:"id"`
	Password string        `yaml:"password"`
}

func readConfig(file string) (c config) {
	// Defaults
	for i := range c.Uploaders {
		c.Uploaders[i].Interval = 1800
	}

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		Error.Fatalf("Unable to read config file: %s", err.Error())
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		Error.Fatalf("Unable to parse config file: %s", err.Error())
	}

	return
}
