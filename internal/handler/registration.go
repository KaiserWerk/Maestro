package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KaiserWerk/Maestro/internal/cache"

	"github.com/KaiserWerk/Maestro/internal/entity"
)

func (h *HttpHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	var reg entity.Registrant
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		http.Error(w, "unable to unmarshal JSON body", http.StatusBadRequest)
		return
	}

	err = h.MaestroCache.Register(reg.Id, reg.Address)
	if err != nil {
		if errors.Is(err, &cache.EntryExists{}) {
			http.Error(w, "unable to register address; ID already exists", http.StatusConflict)
		} else {
			http.Error(w, "unable to register address: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func (h *HttpHandler) DeregistrationHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing or empty id query parameter", http.StatusBadRequest)
		return
	}

	if ok := h.MaestroCache.Deregister(id); !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
