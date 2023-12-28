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
