package businessFacades

import (
	"encoding/json"
	"log"
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

// TODO: check get request
func GetPGPAccountByStellarPK(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	object := dao.Connection{}
	p := object.GetPGPAccountByStellarPK(vars["stellarPublicKey"])
	rst, err := p.Await()
	if err != nil {
		json.NewEncoder(w).Encode("failed to Get PGP account : " + err.Error())
		return
	}
	log.Println("Await response:", rst)
	json.NewEncoder(w).Encode(rst)

}
