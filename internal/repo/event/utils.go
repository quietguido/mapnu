package event

import (
	"fmt"
)

func createST_GeomFromText(lon, lat, time float64) string {
	pointM := fmt.Sprintf(
		"ST_GeomFromText('POINTM(%f %f %f)', 4326)",
		lon,
		lat,
		time,
	)
	return pointM
}
