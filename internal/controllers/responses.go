package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("Bad Request: %s", err)))
}

func SuccessResponse(w http.ResponseWriter, obj any) {
	marshaled, _ := json.Marshal(obj)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaled)
}

func CreatedResponse(w http.ResponseWriter, obj any) {
	marshaled, _ := json.Marshal(obj)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(marshaled)
}
