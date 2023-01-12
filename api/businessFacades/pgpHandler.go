package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func SavePGPKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var keypair model.RSAKeyPair

	err := json.NewDecoder(r.Body).Decode(&keypair)
	if err != nil {
		log.Error("Failed to decode data.")
		return
	}
	dbcon := dao.Connection{}
	err1 := dbcon.InsertRSAKeyPair((keypair))
	if err1 != nil {
		log.Error("Failed to save data")
	}
	w.WriteHeader(http.StatusOK)
	response := model.RSAKeySaveSuccess{Message: "RSA keypair Saved Successfully"}
	json.NewEncoder(w).Encode(response)
}

func GetRSAPublicKeyBySHA256PK(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	object := dao.Connection{}
	p := object.GetRSAPublicKeyBySHA256PK(vars["sha256pk"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.RSAPublickey)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "SHA256 Public Key Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}

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

