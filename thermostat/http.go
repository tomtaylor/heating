package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/tylerb/graceful.v1"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type API struct {
	thermostat *Thermostat
}

type state struct {
	Temperature       float64 `json:"temperature"`
	TargetTemperature float64 `json:"target_temperature"`
	IsHeating         bool    `json:"is_heating"`
}

func NewAPI(thermostat *Thermostat) *API {
	return &API{
		thermostat: thermostat,
	}
}

func (a *API) RunLoop() {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/").HandlerFunc(a.getState)
	router.Methods("PUT").Path("/").HandlerFunc(a.setState)
	graceful.Run(":8080", 10*time.Second, router)
}

func (a *API) Stop() {

}

func (a *API) getState(w http.ResponseWriter, r *http.Request) {
	s := state{
		Temperature:       a.thermostat.Temperature(),
		TargetTemperature: a.thermostat.TargetTemperature(),
		IsHeating:         a.thermostat.IsOn(),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}

func (a *API) setState(w http.ResponseWriter, r *http.Request) {
	var s state
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &s); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	a.thermostat.SetTargetTemperature(s.TargetTemperature)
	s = state{
		Temperature:       a.thermostat.Temperature(),
		TargetTemperature: a.thermostat.TargetTemperature(),
		IsHeating:         a.thermostat.IsOn(),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}
}
