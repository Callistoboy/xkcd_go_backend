package rest

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"yadro.com/course/api/core"
)

func NewPingHandler(log *slog.Logger, pingers map[string]core.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		m := make(map[string]string, len(pingers))

		for name, pinger := range pingers {
			err := pinger.Ping(r.Context())
			var res string = "ok"
			if err != nil {
				log.Error("Ping failed", "error", err)
				res = "unavailable"
			} else {
				log.Debug("Ping success", "service", name)
			}
			m[name] = res
		}

		answer := core.PingAnswer{
			Replies: m,
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(answer); err != nil {
			log.Error("cannot encode", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewWordsHandler(log *slog.Logger, server core.Normalizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phrase := r.URL.Query().Get("phrase")
		log.Debug("Got phrase to normalize", "phrase", phrase)

		if phrase == "" {
			log.Error("Phrase is empty")
			http.Error(w, core.ErrBadArguments.Error(), http.StatusBadRequest)
			return
		}
		res, err := server.Norm(r.Context(), phrase)

		if err != nil {
			if status.Code(err) == codes.ResourceExhausted {
				log.Error("ResourceExhausted error")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			log.Error("Failed to norm", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Debug("Norm success", "result", res)

		response := map[string]any{
			"words": res,
			"total": len(res),
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(response); err != nil {
			log.Error("cannot encode", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewUpdateHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := updater.Update(r.Context()); err != nil {
			log.Error("error while update", "error", err)
			if errors.Is(err, core.ErrAlreadyExists) {
				http.Error(w, err.Error(), http.StatusAccepted)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func NewUpdateStatsHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := updater.Stats(r.Context())
		if err != nil {
			log.Error("Stats failed", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			log.Error("Cannot encode", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewUpdateStatusHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := updater.Status(r.Context())
		if err != nil {
			log.Error("Status failed", "error", err)
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(map[string]core.UpdateStatus{"status": resp}); err != nil {
			log.Error("Cannot encode", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}

func NewDropHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := updater.Drop(r.Context())
		if err != nil {
			log.Error("Drop failed", "error", err)
			w.WriteHeader(http.StatusBadGateway)

		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func NewSearchHandler(log *slog.Logger, searcher core.Searcher, words core.Normalizer) http.HandlerFunc {
	type Comics struct {
		ID  int    `json:"id"`
		URL string `json:"url"`
	}

	type ComicsReply struct {
		Comics []Comics `json:"comics"`
		Total  int      `json:"total"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		phrase := r.URL.Query().Get("phrase")
		limit := r.URL.Query().Get("limit")
		log.Debug("Got phrase", "phrase", phrase)
		log.Debug("Got limit", "limit", limit)

		if phrase == "" {
			log.Error("phrase is empty")
			http.Error(w, core.ErrBadArguments.Error(), http.StatusBadRequest)
			return
		}

		if limit == "" {
			limit = "10"
		}

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			log.Error("Failed to parse limit", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if limitInt < 0 {
			log.Error("wrong limit", "value", limit)
			http.Error(w, "bad limit", http.StatusBadRequest)
			return
		}

		comics, err := searcher.Search(r.Context(), phrase, limitInt)
		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				http.Error(w, "no comics found", http.StatusNotFound)
				return
			}
			log.Error("error while seaching", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}

		response := ComicsReply{
			Comics: make([]Comics, 0, len(comics)),
			Total:  len(comics),
		}

		for _, c := range comics {
			response.Comics = append(response.Comics, Comics{ID: c.ID, URL: c.URL})
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(response); err != nil {
			log.Error("cannot encode", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewSearchIndexHandler(log *slog.Logger, searcher core.Searcher, words core.Normalizer) http.HandlerFunc {
	type Comics struct {
		ID  int    `json:"id"`
		URL string `json:"url"`
	}

	type ComicsReply struct {
		Comics []Comics `json:"comics"`
		Total  int      `json:"total"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		phrase := r.URL.Query().Get("phrase")
		limit := r.URL.Query().Get("limit")
		log.Debug("Got phrase", "phrase", phrase)
		log.Debug("Got limit", "limit", limit)

		if phrase == "" {
			log.Error("phrase is empty")
			http.Error(w, core.ErrBadArguments.Error(), http.StatusBadRequest)
			return
		}

		if limit == "" {
			limit = "10"
		}

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			log.Error("Failed to parse limit", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if limitInt < 0 {
			log.Error("wrong limit", "value", limit)
			http.Error(w, "bad limit", http.StatusBadRequest)
			return
		}

		comics, err := searcher.SearchIndex(r.Context(), phrase, limitInt)
		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				http.Error(w, "no comics found", http.StatusNotFound)
				return
			}
			log.Error("error while seaching", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}

		response := ComicsReply{
			Comics: make([]Comics, 0, len(comics)),
			Total:  len(comics),
		}

		for _, c := range comics {
			response.Comics = append(response.Comics, Comics{ID: c.ID, URL: c.URL})
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(response); err != nil {
			log.Error("cannot encode", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type Authenticator interface {
	Login(user, password, sub string) (string, error)
}

type Login struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewLoginHandler(log *slog.Logger, auth Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var l Login
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			log.Error("could not decode login form", "error", err)
			http.Error(w, "could not parse login data", http.StatusBadRequest)
			return
		}
		token, err := auth.Login(l.Name, l.Password, "")
		if err != nil {
			log.Error("could not authenticate", "user", l.Name, "error", err)
			http.Error(w, core.ErrUserNotFound.Error(), http.StatusUnauthorized)
		}
		if _, err := w.Write([]byte(token)); err != nil {
			log.Error("failed to write reply", "error", err)
		}
	}
}
