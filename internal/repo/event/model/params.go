package model

import (
	"log"
	"time"
)

type GetMapQueryParams struct {
	FirstQuadLon  float64 `form:"firstlon"`
	FirstQuadLat  float64 `form:"firstlat"`
	SecondQuadLon float64 `form:"secondlon"`
	SecondQuadLat float64 `form:"secondlat"`
	FromTime      string  `form:"fromtime"`
	ToTime        string  `form:"totime"`
}

func (st *GetMapQueryParams) GetFromTime() time.Time {
	parsedTime, err := time.Parse(time.RFC3339, st.FromTime)
	if err != nil {
		log.Printf("Failed to parse FromTime: %v", err)
		return time.Time{} // Return zero time if parsing fails
	}
	return parsedTime
}

func (st *GetMapQueryParams) GetToTime() time.Time {
	parsedTime, err := time.Parse(time.RFC3339, st.ToTime)
	if err != nil {
		log.Printf("Failed to parse ToTime: %v", err)
		return time.Time{} // Return zero time if parsing fails
	}
	return parsedTime
}
