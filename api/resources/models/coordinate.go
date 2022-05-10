package models

type Coordinate struct {
	Longitude float32
	Latitude  float32
}

func NewCoordinate(longitude float32, latitude float32) Coordinate {
	return Coordinate{
		Longitude: longitude,
		Latitude:  latitude,
	}
}
