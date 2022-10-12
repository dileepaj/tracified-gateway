package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/stellar/go/support/log"
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
