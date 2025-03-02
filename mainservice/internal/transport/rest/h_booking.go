package rest

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	bookingModel "github.com/quietguido/mapnu/mainservice/internal/repo/booking/model"
)

func (st *restH) CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
	var createBooking bookingModel.CreateBooking

	if err := JsonBodyDecoding(r, &createBooking); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	bookingId, err := st.services.Booking.Create(r.Context(), createBooking)
	if err != nil {
		st.lg.Error(err.Error())
		RespondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	response := map[string]any{
		"booking_id": bookingId,
		"message":    "Booking created successfully",
	}

	RespondWithJson(w, http.StatusOK, response)
}

func (st *restH) GetBookingByIdHandler(w http.ResponseWriter, r *http.Request) {
	bookingIdStr := r.PathValue("id")
	if bookingIdStr == "" {
		RespondWithError(w, http.StatusBadRequest, "Missing booking ID")
		return
	}

	bookingId, err := strconv.Atoi(bookingIdStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	booking, err := st.services.Booking.GetBookingById(r.Context(), bookingId)
	if err != nil {
		st.lg.Error(err.Error())
		RespondWithError(w, http.StatusNotFound, "Booking not found")
		return
	}

	RespondWithJson(w, http.StatusOK, booking)
}

func (st *restH) GetBookingsForUserHandler(w http.ResponseWriter, r *http.Request) { // change for token
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		RespondWithError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	bookings, err := st.services.Booking.GetBookingsForUser(r.Context(), userID)
	if err != nil {
		st.lg.Error(err.Error())
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve bookings")
		return
	}

	RespondWithJson(w, http.StatusOK, bookings)
}

func (st *restH) ChangeBookingStatusHandler(w http.ResponseWriter, r *http.Request) { // change for token
	var changeBookingStatus bookingModel.ChangeBookingStatus
	if err := JsonBodyDecoding(r, &changeBookingStatus); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := st.services.Booking.ChangeBookingStatus(r.Context(), changeBookingStatus)
	if err != nil {
		st.lg.Error(err.Error())
		RespondWithError(w, http.StatusInternalServerError, "Failed to change status")
		return
	}
}

func (st *restH) GetBookingApplicationsForOrganizer(w http.ResponseWriter, r *http.Request) { // change for token

}
