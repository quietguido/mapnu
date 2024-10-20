package rest

import (
	"net/http"

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
