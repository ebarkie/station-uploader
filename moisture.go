// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

// cbToPercent converts soil moisture tension in centibars to a percentage.
func cbToPercent(soilType int, cb int) int {
	// This uses a simple linear scale based on depletion of plant
	// available water for each soil type.

	// Scheduling Irrigations: When and How Much Water to Apply.
	// Division of Agriculture and Natural Resources Publication 3396.
	// University of California Irrigation Program.
	// University of California, Davis. pp. 106.
	var depletedCbs = []int{
		30,  // Sand/Loamy Sand
		50,  // Sandy Loam
		130, // Loam/Silt Loam
		170, // Clay Loam/Clay
	}

	var dcb int
	if soilType < 0 {
		dcb = depletedCbs[0]
	} else if soilType >= len(depletedCbs) {
		dcb = depletedCbs[len(depletedCbs)-1]
	} else {
		dcb = depletedCbs[soilType]
	}

	if cb > dcb {
		return 0
	}
	return 100 - int((float32(cb)/float32(dcb))*100)
}
