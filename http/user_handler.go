package http

import (
	"encoding/json"
	"net/http"

	"rankit/errors"
	"rankit/rankit"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	var p rankit.CreateUserParam
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return nil, http.StatusBadRequest, errors.E(errors.Invalid, "invalid request body")
	}

	user, err := s.usersvc.CreateUser(r.Context(), p)
	if err == nil {
		s.sessionmgr.Put(r.Context(), USER_ID_SESSION_KEY, user.ID)
		return user, http.StatusCreated, nil
	}

	return
}

func (s *Server) AuthenticateUser(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	var p rankit.AuthenticateUserParam
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return nil, http.StatusBadRequest, errors.E(errors.Invalid, "invalid request body")
	}

	user, err := s.usersvc.Authenticate(r.Context(), p)
	if err == nil {
		s.sessionmgr.Put(r.Context(), USER_ID_SESSION_KEY, user.ID)
		return user, http.StatusOK, nil
	}

	return
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	userID := s.sessionmgr.GetString(r.Context(), USER_ID_SESSION_KEY)
	if userID == "" {
		return nil, http.StatusUnauthorized, errors.E(errors.Unauthorized, "unauthorized")
	}

	user, err := s.usersvc.GetUser(r.Context(), userID)
	if err == nil {
		return user, http.StatusOK, nil
	}

	return
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	s.sessionmgr.Remove(r.Context(), USER_ID_SESSION_KEY)
	RespondMessage(w, http.StatusOK, "logged out")
	return
}
