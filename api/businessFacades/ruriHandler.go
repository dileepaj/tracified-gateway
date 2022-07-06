package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	ruri "github.com/dileepaj/tracified-gateway/nft/stellar/RURI"
	"github.com/sirupsen/logrus"
)

func SponsorBuyer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	key1, error := r.URL.Query()["publickey"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'publickey' is missing")
		return
	}

	key2, error := r.URL.Query()["nftName"]

	if !error || len(key2[0]) < 1 {
		logrus.Error("Url Parameter 'nftName' is missing")
		return
	}

	key3, error := r.URL.Query()["issuer"]

	if !error || len(key3[0]) < 1 {
		logrus.Error("Url Parameter 'issuer' is missing")
		return
	}
	publickey := key1[0]
	nftname := key2[0]
	issuer := key3[0]

	logrus.Println("values ------------ ", publickey, nftname, issuer)

	var txn, err = ruri.SponsorCreateAccount(publickey, nftname, issuer)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error when updating the buying status",
		}
		json.NewEncoder(w).Encode(result)
	} else {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusOK)
		result := model.XDRRuri{
			XDR: txn,
		}
		logrus.Println("XDR been passed to frontend : ", result)
		json.NewEncoder(w).Encode(result)
	}

}
