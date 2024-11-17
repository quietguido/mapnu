package model

import "time"

type GetMapQueryParams struct {
	FirstQuadLon  float64   `form:"firstlon"`
	FirstQuadLat  float64   `form:"firstlat"`
	SecondQuadLon float64   `form:"secondlon"`
	SecondQuadLat float64   `form:"secondlat"`
	FromTime      time.Time `form:"fromtime"`
	ToTime        time.Time `form:"totime"`
}
