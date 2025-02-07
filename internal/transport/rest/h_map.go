package rest

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	eventModel "github.com/quietguido/mapnu/internal/repo/event/model"
)

func (st *restH) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var createEvent eventModel.CreateEvent

	if err := JsonBodyDecoding(r, &createEvent); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := st.services.Event.Create(r.Context(), createEvent)
	if err != nil {
		st.lg.Error(err.Error())
		RespondWithError(w, http.StatusBadRequest, "bad request")
		return
	}
}

func (st *restH) GetEventByIdHandler(w http.ResponseWriter, r *http.Request) {
	eventId := r.PathValue("id")
	if eventId == "" {
		RespondWithError(w, http.StatusBadRequest, "bad request")
		return
	}
	err := uuid.Validate(eventId)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "not valid UUID")
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
	// Parse query parameters into the mapQuery struct
	var mapQuery eventModel.GetMapQueryParams
	if err := DecodeQuery(r, &mapQuery); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if mapQuery.FirstQuadLon == 0 || mapQuery.FirstQuadLat == 0 || mapQuery.SecondQuadLon == 0 || mapQuery.SecondQuadLat == 0 {
		RespondWithError(w, http.StatusBadRequest, "Missing or invalid quadrant parameters")
		return
	}

	fromTime, err := time.Parse(time.RFC3339, mapQuery.FromTime)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid 'fromtime' format. Expected RFC3339 (e.g., 2024-11-17T14:17:01Z).")
		return
	}

	toTime, err := time.Parse(time.RFC3339, mapQuery.ToTime)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid 'totime' format. Expected RFC3339 (e.g., 2024-11-17T14:17:01Z).")
		return
	}

	if toTime.Before(fromTime) {
		RespondWithError(w, http.StatusBadRequest, "'totime' must be after 'fromtime'.")
		return
	}

	mapQuery.FromTime = fromTime.Format(time.RFC3339)
	mapQuery.ToTime = toTime.Format(time.RFC3339)

	events, err := st.services.Event.GetMapForQuadrant(r.Context(), mapQuery)
	if err != nil {
		st.lg.Error(err.Error())
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}

	RespondWithJson(w, http.StatusOK, events)
}
