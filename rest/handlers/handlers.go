package handlers

import (
	"encoding/json"
	"net/http"
	"solid_go/rest/ent"
	"solid_go/rest/uc"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func GetHandler(getUc uc.GetUc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		name, ok := getUc(id)
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

func PostHandler(saveUc uc.SaveUc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var data ent.Data

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		saveUc(data.Id, data.Name)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	}
}
