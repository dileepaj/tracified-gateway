package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/gorilla/mux"
)

func SavePGPAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var PGPResponse model.PGPAccount
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&PGPResponse)
	if err != nil {
		panic(err)
	}
	object := dao.Connection{}
	err1 := object.InsertPGPAccount(PGPResponse)
	if err1 != nil {
		panic(err1)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PGPResponse)
	return
}

func GetPGPAccountByStellarPK(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	object := dao.Connection{}
	p := object.GetPGPAccountByStellarPK(vars["stellarPublicKey"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.PGPAccount)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "StellarPK Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}
