package v1

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (api *API) GetUsers(w http.ResponseWriter, r *http.Request) error {

	users, err := api.userService.List(r.Context())
	if err != nil {
		return err
	}

	resp := make([]UserResponse, len(users))

	for i, u := range users {
		resp[i] = *ToUserResponse(u)
	}

	return api.JSON(w, http.StatusOK, resp)
}

func (api *API) GetUserByID(w http.ResponseWriter, r *http.Request) error {

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return NewAPIError(http.StatusBadRequest, ErrInvalidID)
	}

	user, err := api.userService.Get(r.Context(), id)
	if err != nil {
		return err
	}

	return api.JSON(w, http.StatusOK, ToUserResponse(user))
}

func (api *API) CreateUser(w http.ResponseWriter, r *http.Request) error {

	input := &CreateUserRequest{}

	err := decodeIntoStruct(r, input)
	if err != nil {
		return NewAPIError(http.StatusBadRequest, err)
	}

	user, err := input.ToUser()
	if err != nil {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid user data: %s", err))
	}

	err = api.userService.Create(r.Context(), user)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

func (api *API) UpdateUser(w http.ResponseWriter, r *http.Request) error {

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return NewAPIError(http.StatusBadRequest, ErrInvalidID)
	}

	input := &UpdateUserRequest{}

	err = decodeIntoStruct(r, input)
	if err != nil {
		return NewAPIError(http.StatusBadRequest, err)
	}

	user, err := input.ToUser()
	if err != nil {
		return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid user data: %s", err))
	}

	user.ID = id

	err = api.userService.Update(r.Context(), user)
	if err != nil {
		return err
	}

	return api.JSON(w, http.StatusOK, ToUserResponse(user))
}

func (api *API) DeleteUser(w http.ResponseWriter, r *http.Request) error {

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return NewAPIError(http.StatusBadRequest, ErrInvalidID)
	}

	err = api.userService.Delete(r.Context(), id)
	if err != nil {
		return err
	}

	return nil
}
