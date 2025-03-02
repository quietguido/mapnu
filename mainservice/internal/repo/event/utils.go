package event

import (
	"fmt"
	"time"
)

// func createPoint(lon, lat float64) string {
// 	return fmt.Sprintf("ST_SetSRID(ST_Point(%f, %f), 4326)", lon, lat)
// }

func getPartition(t time.Time) string {
	return fmt.Sprintf("%s_%d_%02d_%02d", eventTable, t.Year(), t.Month(), t.Day())
}
