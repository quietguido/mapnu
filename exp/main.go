package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/quietguido/mapnu/internal/database/psql"
	"github.com/quietguido/mapnu/pkg/assert"
)

const (
	astana_lat = 51.169392
	astana_lon = 71.449074
)

func main() {
	err := godotenv.Load("config.env")
	assert.ErrorNil(err, "failed to load config.env")

	// lg, err := zap.NewProduction()
	// assert.ErrorNil(err, "lg creation error")

	dbcon, err := psql.New(psql.Config{
		Addr:     os.Getenv("POSTGRES_HOST"), //change for local and docker
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DB:       os.Getenv("POSTGRES_DB"),
	})
	assert.ErrorNil(err, "failed db connection")

	//create 1000 elements in the dp
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 1000000; i++ {
		currentTime := time.Now()
		rand_time := currentTime.Unix() - rand.Int63n(1000000) + 2000000
		rand_lon := astana_lon - rand.Float64() + 0.5
		rand_lat := astana_lat - rand.Float64() + 0.5

		query := `
		INSERT INTO event (
		    name,
		    description,
		    created_by,
		    location_date,
		    organizer,
		    upvote,
		    downvote
		) VALUES (
		    'My Event Name',
		    'Event Description',
		    '6dd4095d-e156-43eb-ab94-7651e68ca989',
			ST_MakePointM($1, $2, $3), -- Example location (San Francisco)
		    'Organizer Name',
		    0,
		    0
		)
		`

		_, err := dbcon.Exec(query, []interface{}{
			rand_lon,
			rand_lat,
			rand_time,
		}...)

		assert.ErrorNil(err, "failed db connection")
	}
}
