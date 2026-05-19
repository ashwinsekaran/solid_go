package handlers

import (
	"encoding/json"
	"net/http"
	"solid_go/rest/ent"
	"solid_go/rest/uc"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// GetHandler returns an httprouter.Handle that retrieves a record by ID.
// It closes over getUc (a use-case function) — the handler knows nothing about
// the storage layer; it only speaks HTTP and delegates work to the use-case.
func GetHandler(getUc uc.GetUc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Parse and validate the path parameter — respond 400 if it's not an int.
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		name, ok := getUc(id) // delegate to business logic
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Record not found"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(name))
	}
}

// PostHandler returns an httprouter.Handle that creates a new record.
// JSON decoding happens here (HTTP concern); saving happens in saveUc (business concern).
func PostHandler(saveUc uc.SaveUc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var data ent.Data

		// Decode the JSON body — respond 400 on malformed input.
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		saveUc(data.Id, data.Name) // delegate persistence to the use-case

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	}
}
