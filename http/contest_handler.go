package http

import (
	"encoding/json"
	"net/http"
	"rankit/errors"
	"rankit/rankit"

	"github.com/go-chi/chi/v5"
)

func (s *Server) CreateContest(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	var p rankit.CreateContestParam
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return nil, http.StatusBadRequest, errors.E(errors.Invalid, "invalid request body")
	}
	userID := s.sessionmgr.GetString(r.Context(), USER_ID_SESSION_KEY)
	p.CreatorID = userID

	contest, err := s.contestsvc.CreateContest(r.Context(), p)
	if err == nil {
		return contest, http.StatusCreated, nil
	}

	return
}

func (s *Server) AddItem(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	var p rankit.AddItemParam
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return nil, http.StatusBadRequest, errors.E(errors.Invalid, "invalid request body")
	}

	item, err := s.contestsvc.AddItem(r.Context(), p)
	if err == nil {
		return item, http.StatusCreated, nil
	}

	return
}

func (s *Server) GetContest(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	id := chi.URLParam(r, "contest_id")
	contest, err := s.contestsvc.GetContest(r.Context(), id)
	if err == nil {
		return contest, http.StatusOK, nil
	}
	return
}

func (s *Server) RecordVote(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	var p rankit.RecordVoteParam
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return nil, http.StatusBadRequest, errors.E(errors.Invalid, "invalid request body")
	}
	userID := s.sessionmgr.GetString(r.Context(), USER_ID_SESSION_KEY)
	p.VoterID = userID

	if err := s.contestsvc.RecordVote(r.Context(), p); err == nil {
		RespondMessage(w, http.StatusOK, "vote recorded")
	}
	return
}

func (s *Server) GetMatchUp(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	id := chi.URLParam(r, "contest_id")
	items, err := s.contestsvc.GetMatchUp(r.Context(), id)
	if err == nil {
		return items, http.StatusOK, nil
	}
	return
}

func (s *Server) GetItemEloHistory(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	id := chi.URLParam(r, "item_id")
	history, err := s.contestsvc.GetItemEloHistory(r.Context(), id)
	if err == nil {
		return history, http.StatusOK, nil
	}
	return
}
