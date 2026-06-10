package main

import "math"

func square(a float64) float64 {
	return a * a
}

func radiansFromDegrees(degrees float64) float64 {
	return 0.01745329251994329577 * degrees
}

func ReferenceHaversine(x0, y0, x1, y1, earthRadius float64) float64 {

	/* NOTE(casey): This is not meant to be a "good" way to calculate the Haversine distance.
	   Instead, it attempts to follow, as closely as possible, the formula used in the real-world
	   question on which these homework exercises are loosely based.
	*/

	lat1 := y0
	lat2 := y1
	lon1 := x0
	lon2 := x1

	dLat := radiansFromDegrees(lat2 - lat1)
	dLon := radiansFromDegrees(lon2 - lon1)
	lat1 = radiansFromDegrees(lat1)
	lat2 = radiansFromDegrees(lat2)

	a := square(math.Sin(dLat/2.0)) + math.Cos(lat1)*math.Cos(lat2)*square(math.Sin(dLon/2))
	c := 2.0 * math.Asin(math.Sqrt(a))

	result := earthRadius * c

	return result
}
