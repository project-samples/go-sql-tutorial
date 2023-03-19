package handler

import (
	"encoding/json"
	"github.com/core-go/core"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"

	. "go-service/internal/model"
	. "go-service/internal/service"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) All(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.All(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, result)
}

func (h *UserHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	result, err := h.service.Load(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, result)
}

func (h *UserHandler) Insert(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}

	result, er2 := h.service.Insert(r.Context(), &user)
	if er2 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, result)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		http.Error(w, "Id not match", http.StatusBadRequest)
		return
	}

	result, er2 := h.service.Update(r.Context(), &user)
	if er2 != nil {
		http.Error(w, er2.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, result)
}

func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	ids := []string{"id"}

	var user User
	userType := reflect.TypeOf(user)
	_, jsonMap, _ := core.BuildMapField(userType)
	body, _ := core.BuildMapAndStruct(r, &user)
	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		http.Error(w, "Id not match", http.StatusBadRequest)
		return
	}
	json, er1 := core.BodyToJsonMap(r, user, body, ids, jsonMap)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}

	result, er2 := h.service.Patch(r.Context(), json)
	if er2 != nil {
		http.Error(w, er2.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, result)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	result, err := h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, result)
}

func respond(w http.ResponseWriter, result interface{}) {
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
