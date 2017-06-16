# Station Uploader

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](http://choosealicense.com/licenses/mit/)
[![Build Status](https://travis-ci.org/ebarkie/station-uploader.svg?branch=master)](https://travis-ci.org/ebarkie/station-uploader)

Retrieves observations from a Davis station and sends them to APRS
(TNC or IS (Citizen Weather Observer Program (CWOP))) and Weather
Underground.

## Installation

```
$ go get
$ go generate
$ go build
```

Refer to the [contrib](contrib) directory for a sample systemd service.

## Configuration

The configuration file is formatted as YAML.  It contains the Davis Instruments
weather station information and a series of uploaders to use.

```yaml
station:
  host: wx:8080
  lat: 35.7
  lon: -78.7
  type: DvsVP2+

uploaders:
  - name: APRS-IS
    type: aprs-is
    interval: 3600
    dial: tcp://rotate.aprs.net:14580
    id: N0CALL-13
    password: -1

  - name: APRS-RF
    type: aprs
    interval: 3600
    dial: direwolf:8001
    id: N0CALL-13

  - name: CWOP
    type: aprs-is
    interval: 300
    dial: tcp://cwop.aprs.net:14580
    id: aWnnnn

  - name: Wunderground
    type: wu
    interval: 0
    id: Kssssssnn
    password: deadbeef
```

## Usage

```
Usage of ./station-uploader:
  -config string
        station uploader config file (default "station-uploader.yaml")
  -debug
        enable debug mode

$ ./station-uploader
```

## License

Copyright (c) 2016-2017 Eric Barkie. All rights reserved.  
Use of this source code is governed by the MIT license
that can be found in the [LICENSE](LICENSE) file.
