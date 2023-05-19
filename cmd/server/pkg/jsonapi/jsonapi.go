// empyr - a reimagining of Vern Holford's Empyrean Challenge
// Copyright (C) 2023 Michael D Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

// Package jsonapi tries to implement JSONAPI.
package jsonapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Errors(w http.ResponseWriter, r *http.Request, status int, details interface{}) {
	var failure struct {
		Status string      `json:"status"`
		Errors interface{} `json:"errors"`
		Links  struct {
			Self string `json:"self"`
		} `json:"links"`
	}
	failure.Status = http.StatusText(status)
	failure.Links.Self = r.URL.Path
	failure.Errors = details

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(failure); err != nil {
		log.Printf("[http] error writing response: %+v\n", err)
	}
}

func Error(w http.ResponseWriter, r *http.Request, status int, errors ...error) {
	type errorSource struct {
		Pointer   string `json:"pointer,omitempty"`
		Parameter string `json:"parameter,omitempty"`
	}
	type errorObject struct {
		Id     string       `json:"id,omitempty"`
		Status string       `json:"status,omitempty"`
		Code   string       `json:"code,omitempty"`
		Title  string       `json:"title,omitempty"`
		Detail string       `json:"detail,omitempty"`
		Source *errorSource `json:"source,omitempty"`
	}
	var failure struct {
		Status string        `json:"status"`
		Errors []errorObject `json:"errors"`
		Links  struct {
			Self string `json:"self"`
		} `json:"links"`
	}
	failure.Status = http.StatusText(status)
	failure.Links.Self = r.URL.Path

	// the first error, by convention, is always the http status being returned
	failure.Errors = append(failure.Errors, errorObject{
		Status: fmt.Sprintf("%d", status),
		Detail: http.StatusText(status),
		Source: &errorSource{Parameter: r.URL.Path},
	})

	// then append any error details that the user provided
	for _, err := range errors {
		failure.Errors = append(failure.Errors, errorObject{
			Detail: fmt.Sprintf("%+v", err),
		})
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(failure); err != nil {
		log.Printf("[http] error writing response: %+v\n", err)
	}
}

func NoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func Ok(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)

	var success struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
		Links  struct {
			Self string `json:"self"`
		} `json:"links"`
	}
	success.Status = "ok"
	success.Data = data
	success.Links.Self = r.URL.Path

	if err := json.NewEncoder(w).Encode(success); err != nil {
		log.Printf("[http] error writing response: %+v\n", err)
	}
}
