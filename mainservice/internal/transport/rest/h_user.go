package rest

import (
	"net/http"

	"github.com/google/uuid"
	userModel "github.com/quietguido/mapnu/mainservice/internal/repo/user/model"
)

func (st *restH) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var createUser userModel.CreateUser

	if err := JsonBodyDecoding(r, &createUser); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := st.services.User.CreateUser(r.Context(), createUser)
	if err != nil {
		st.lg.Error(err.Error())
		RespondWithError(w, http.StatusBadRequest, "bad request")
		return
	}
}

func (st *restH) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")
	if userId == "" {
		RespondWithError(w, http.StatusBadRequest, "bad request")
		return
	}
	err := uuid.Validate(userId)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "not valid UUID")
		return
	}

	userModel, err := st.services.User.GetUserById(r.Context(), userId)
	if err != nil {
		st.lg.Error(err.Error())
		RespondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	RespondWithJson(w, http.StatusOK, userModel)
}
