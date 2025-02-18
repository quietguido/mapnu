package rest

import (
	"net/http"
	"strconv"
	"time"

	eventModel "github.com/quietguido/mapnu/internal/repo/event/model"
)

func (st *restH) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var createEvent eventModel.CreateEvent

	if err := JsonBodyDecoding(r, &createEvent); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	eventId, err := st.services.Event.Create(r.Context(), createEvent)
	if err != nil {
		st.lg.Error(err.Error())
		RespondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	response := map[string]any{
		"event_id": eventId,
		"message":  "Event created successfully",
	}

	RespondWithJson(w, http.StatusOK, response)
}

func (st *restH) GetEventByIdHandler(w http.ResponseWriter, r *http.Request) {
	eventIdStr := r.PathValue("id")
	if eventIdStr == "" {
		RespondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	eventModel, err := st.services.Event.GetEventById(r.Context(), eventId)
	if err != nil {
		st.lg.Error(err.Error())
		RespondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	RespondWithJson(w, http.StatusOK, eventModel)
}

func (st *restH) GetMapForQuadrantHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	firstLon, err := strconv.ParseFloat(query.Get("firstlon"), 64)
	if err != nil {
		http.Error(w, "Invalid firstlon parameter", http.StatusBadRequest)
		return
	}
	firstLat, err := strconv.ParseFloat(query.Get("firstlat"), 64)
	if err != nil {
		http.Error(w, "Invalid firstlat parameter", http.StatusBadRequest)
		return
	}
	secondLon, err := strconv.ParseFloat(query.Get("secondlon"), 64)
	if err != nil {
		http.Error(w, "Invalid secondlon parameter", http.StatusBadRequest)
		return
	}
	secondLat, err := strconv.ParseFloat(query.Get("secondlat"), 64)
	if err != nil {
		http.Error(w, "Invalid secondlat parameter", http.StatusBadRequest)
		return
	}

	dateStr := query.Get("date")
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		http.Error(w, "Invalid date parameter (must be RFC3339 format)", http.StatusBadRequest)
		return
	}

	queryParams := eventModel.GetMapQueryParams{
		FirstQuadLon:  firstLon,
		FirstQuadLat:  firstLat,
		SecondQuadLon: secondLon,
		SecondQuadLat: secondLat,
		Date:          date,
	}

	if queryParams.FirstQuadLon == 0 || queryParams.FirstQuadLat == 0 || queryParams.SecondQuadLon == 0 || queryParams.SecondQuadLat == 0 {
		RespondWithError(w, http.StatusBadRequest, "Missing or invalid quadrant parameters")
		return
	}

	events, err := st.services.Event.GetMapForQuadrant(r.Context(), queryParams)
	if err != nil {
		st.lg.Error(err.Error())
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}

	RespondWithJson(w, http.StatusOK, events)
}
