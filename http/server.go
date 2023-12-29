package http

import (
	"encoding/json"
	"net/http"

	"rankit/errors"
	"rankit/rankit"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"moul.io/chizap"
)

type Server struct {
	usersvc    rankit.UserService
	contestsvc rankit.ContestService
	sessionmgr *scs.SessionManager
	logger     *zap.Logger
}

func NewServer(
	usersvc rankit.UserService,
	contestsvc rankit.ContestService,
	sessionmgr *scs.SessionManager,
	logger *zap.Logger,
) *Server {
	return &Server{
		usersvc:    usersvc,
		contestsvc: contestsvc,
		sessionmgr: sessionmgr,
		logger:     logger,
	}
}

func (s *Server) Listen(addr string) error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(chizap.New(s.logger, &chizap.Opts{
		WithReferer:   false,
		WithUserAgent: false,
	}))
	r.Route("/v1/users", func(r chi.Router) {
		r.Post("/", s.toHTTPHandlerFunc(s.CreateUser))
		r.Post("/login", s.toHTTPHandlerFunc(s.AuthenticateUser))
		r.Post("/logout", s.toHTTPHandlerFunc(s.Logout))
	})
	r.Route("/v1/contests", func(r chi.Router) {
		r.Use(Auth(s.sessionmgr))
		r.Post("/", s.toHTTPHandlerFunc(s.CreateContest))
		r.Post("/{contest_id}/items", s.toHTTPHandlerFunc(s.AddItem))
		r.Get("/{contest_id}", s.toHTTPHandlerFunc(s.GetContest))
		r.Get("/{contest_id}/matchup", s.toHTTPHandlerFunc(s.GetMatchUp))
		r.Post("/{contest_id}/vote", s.toHTTPHandlerFunc(s.RecordVote))
		r.Get("/{contest_id}/items/{item_id}/history", s.toHTTPHandlerFunc(s.GetItemEloHistory))
	})

	return http.ListenAndServe(addr, s.sessionmgr.LoadAndSave(r))
}

func (s *Server) toHTTPHandlerFunc(handler func(w http.ResponseWriter, r *http.Request) (any, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, status, err := handler(w, r)
		if err != nil {
			switch err := err.(type) {
			case *errors.Error:
				RespondError(w, err)
			default:
				s.logger.Error("internal error", zap.Error(err))
				RespondMessage(w, http.StatusInternalServerError, "Internal Error")
			}
			return
		}
		if response != nil {
			RespondJSON(w, status, response)
		}
		if status >= 100 && status < 600 {
			w.WriteHeader(status)
		}
	}
}

func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if json.NewEncoder(w).Encode(data) != nil {
		http.Error(w, `{"message": "Internal Error"}`, http.StatusInternalServerError)
	}
}

func RespondMessage(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{"message": message})
}

func RespondError(w http.ResponseWriter, err *errors.Error) {
	switch err.Kind {
	case errors.NotFound:
		RespondMessage(w, http.StatusNotFound, err.Message)
	case errors.Invalid:
		var ve errors.ValidationErrors
		if errors.As(err, &ve) {
			RespondJSON(w, http.StatusBadRequest, map[string]any{
				"message":           err.Message,
				"validation_errors": ve,
			})
			return
		}
		RespondMessage(w, http.StatusBadRequest, err.Message)
	case errors.Unauthorized:
		RespondMessage(w, http.StatusUnauthorized, err.Message)
	case errors.Forbidden:
		RespondMessage(w, http.StatusForbidden, err.Message)
	default:
		RespondMessage(w, http.StatusInternalServerError, err.Message)
	}
}
