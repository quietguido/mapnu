package event

import (
	"fmt"
	"time"
)

// Helper function to generate a PostGIS POINTM WKT string
func createPointM(lon, lat float64, time time.Time) string {
	return fmt.Sprintf("POINTM(%f %f %d)", lon, lat, time.Unix())
}
