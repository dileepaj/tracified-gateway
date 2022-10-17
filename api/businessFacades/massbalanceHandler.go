package businessFacades

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/nft/stellar/massbalance"
)

func SplitBatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Println("-------start now")
	object := dao.Connection{}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var obj model.Splits
	var subobj model.Destination
	err = json.Unmarshal(b, &obj)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println("-------data", obj)
	log.Println("array ", len(obj.Destination))
	er := object.InsertSplitData(obj)
	if er != nil {
		log.Println("Failed to save data")
	}

	for i := 0; i < len(obj.Destination); i++ {
		subobj = model.Destination{
			Source: obj.Destination[i].Source,
			Amount: obj.Destination[i].Amount,
		}

		log.Println("-------sub data", subobj)

		var result, err1 = massbalance.Split(subobj.Source, subobj.Amount, obj.NFTName, obj.Sender, obj.Issuer, obj.Limit)
		if err1 != nil {
			ErrorMessage := err1.Error()
			log.Println(w, ErrorMessage)
			return
		} else {
			var batchData = model.Batches{
				NFTName:       obj.NFTName,
				TXNHash:       result,
				CurrentOwner:  subobj.Source,
				PreviousOwner: obj.Sender,
			}

			var resultObj = object.BatchTrackingData(batchData)
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(resultObj)
			if err != nil {
				log.Println(err)
			}
			return
		}

	}
}

func MergeBatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Println("-------start now")
	object := dao.Connection{}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var obj model.Merges
	var subobj model.Destination
	err = json.Unmarshal(b, &obj)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println("-------data", obj)
	log.Println("array ", len(obj.Sender))

	er := object.InsertMergeData(obj)
	if er != nil {
		log.Println("Failed to save data")
	}

	for i := 0; i < len(obj.Sender); i++ {
		subobj = model.Destination{
			Source: obj.Sender[i].Source,
			Amount: obj.Sender[i].Amount,
		}

		log.Println("-------sub data", subobj)

		var result, err1 = massbalance.Merge(subobj.Source, subobj.Amount, obj.NFTName, obj.Destination, obj.Issuer, obj.Limit)
		if err1 != nil {
			ErrorMessage := err1.Error()
			log.Println(w, ErrorMessage)
			return
		} else {
			var batchData = model.Batches{
				NFTName:       obj.NFTName,
				TXNHash:       result,
				CurrentOwner:  obj.Destination,
				PreviousOwner: subobj.Source,
			}

			var resultObj = object.BatchTrackingData(batchData)
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(resultObj)
			if err != nil {
				log.Println(err)
			}
			return
		}

	}

}
