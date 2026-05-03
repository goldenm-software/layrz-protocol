package definitions

// Position defines the position data structure
type Position struct {
	// Is the latitude of the device
	Latitude *float64
	// Is the longitude of the device
	Longitude *float64
	// Is the altitude of the device
	Altitude *float64
	// Is the speed of the device
	Speed *float64
	// Is the direction of the device
	Direction *float64
	// Is the satellite count
	SatelliteCount *int
	// Is the HDOP value
	Hdop *float64
}
