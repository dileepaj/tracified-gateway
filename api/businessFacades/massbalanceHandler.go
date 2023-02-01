package businessFacades

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
	//Move DB call after BC transaction occurs success fully
	er := object.InsertSplitData(obj)
	if er != nil {
		log.Println("Failed to save data")
	}
	var total = 0
	for i := 0; i < len(obj.Destination); i++ {
		amount, converterr := strconv.Atoi(obj.Destination[i].Amount)
		if converterr != nil {
			log.Fatal(converterr.Error())
			http.Error(w, "incorrect amount submittied", 500)
			return
		}
		total += amount
	}
	limit, limitError := strconv.Atoi(obj.Limit)
	if limitError != nil {
		http.Error(w, "Error occured", 500)
		return
	}
	if total > limit {
		http.Error(w, "Amounts provided for split cannot exceed limit", 500)
		return
	}
	if total < limit {
		http.Error(w, "Addition of values provided is less than the limit", 500)
		return
	}

	for i := 0; i < len(obj.Destination); i++ {
		log.Println("------------------------------- the element ", i)
		subobj = model.Destination{
			Source: obj.Destination[i].Source,
			Sign:   obj.Destination[i].Sign,
			Amount: obj.Destination[i].Amount,
		}

		log.Println("-------sub data", subobj)

		var result, err1 = massbalance.Split(subobj.Source, subobj.Sign, subobj.Amount, obj.NFTName, obj.Sender, obj.Issuer, obj.Limit)
		if err1 != nil {
			ErrorMessage := err1.Error()
			log.Println(w, ErrorMessage)
			return
		} else {
			var result1, err2 = massbalance.SplitPayment(subobj.Source, subobj.Sign, subobj.Amount, obj.NFTName, obj.Sender, obj.Issuer, obj.Limit)
			if err2 != nil {
				ErrorMessage := err2.Error()
				log.Println(w, ErrorMessage)
				return
			} else {

				var batchData = model.Batches{
					NFTName:         obj.NFTName,
					TXNHashTrust:    result,
					TXNHashTransfer: result1,
					CurrentOwner:    subobj.Source,
					PreviousOwner:   obj.Sender,
				}

				resultObj, dberr := object.BatchTrackingData(batchData)
				if dberr != nil {
					http.Error(w, "Error occured", 500)
					log.Fatal(dberr.Error())
					return
				}
				w.WriteHeader(http.StatusOK)
				err = json.NewEncoder(w).Encode(resultObj)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	log.Println("done with for loop")
	return
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
	var total = 0
	fmt.Println("arr size for merge:", len(obj.Sender))
	for i := 0; i < len(obj.Sender); i++ {
		amount, converterr := strconv.Atoi(obj.Sender[i].Amount)
		if converterr != nil {
			log.Fatal(converterr.Error())
			http.Error(w, "incorrect amount submittied", 500)
			return
		}
		total += amount
	}
	limit, limitError := strconv.Atoi(obj.Limit)
	if limitError != nil {
		http.Error(w, "Error occured", 500)
		return
	}
	if total > limit {
		http.Error(w, "Amounts provided for Merge cannot exceed limit", 500)
		return
	}
	if total < limit {
		http.Error(w, "Addition of values provided is less than thhe limit", 500)
		return
	}

	for i := 0; i < len(obj.Sender); i++ {
		log.Println("---------------element ", i)
		subobj = model.Destination{
			Source: obj.Sender[i].Source,
			Sign:   obj.Sender[i].Sign,
			Amount: obj.Sender[i].Amount,
		}

		log.Println("-------sub data", subobj)

		var result, err1 = massbalance.Merge(subobj.Source, subobj.Sign, subobj.Amount, obj.NFTName, obj.Destination, obj.Issuer, obj.Limit)
		if err1 != nil {
			ErrorMessage := err1.Error()
			log.Println(w, ErrorMessage)
			return
		} else {
			var result1, err2 = massbalance.TransferMerge(subobj.Source, subobj.Sign, subobj.Amount, obj.NFTName, obj.Destination, obj.Issuer, obj.Limit)
			if err2 != nil {
				ErrorMessage := err2.Error()
				log.Println(w, ErrorMessage)
				return
			} else {
				var batchData = model.Batches{
					NFTName:         obj.NFTName,
					TXNHashTrust:    result,
					TXNHashTransfer: result1,
					CurrentOwner:    obj.Destination,
					PreviousOwner:   subobj.Source,
				}

				resultObj, dberr := object.BatchTrackingData(batchData)
				if dberr != nil {
					http.Error(w, "Error occured", 500)
					log.Fatal(dberr.Error())
					return
				}
				w.WriteHeader(http.StatusOK)
				err = json.NewEncoder(w).Encode(resultObj)
				if err != nil {
					log.Println(err)
				}

			}
		}
	}
	log.Println("done with for loop")
	return
}

func Conversions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	object := dao.Connection{}
	var obj model.Conversions
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&obj)
	if err != nil {
		panic(err)
	}
	result1, err1 := massbalance.SetConversion(obj.Sender, obj.Amount, obj.SellAsset, obj.BuyAsset, obj.SellIssuer, obj.BuyIssuer, obj.Numerator, obj.Denominator)
	if err1 != nil {
		ErrorMessage := err1.Error()
		log.Println(w, ErrorMessage)
		return
	} else {
		newobj := model.Conversions{
			Sender:      obj.Sender,
			Amount:      obj.Amount,
			SellAsset:   obj.SellAsset,
			BuyAsset:    obj.BuyAsset,
			SellIssuer:  obj.SellIssuer,
			BuyIssuer:   obj.BuyIssuer,
			Numerator:   obj.Numerator,
			Denominator: obj.Denominator,
			TXNHash:     result1,
		}
		object.ConvertBatches(newobj)
		result2, err2 := massbalance.ConvertBatches(obj.Sender, obj.Amount, obj.SellAsset, obj.BuyAsset, obj.SellIssuer, obj.BuyIssuer, obj.Numerator, obj.Denominator)
		if err2 != nil {
			ErrorMessage := err2.Error()
			log.Println(w, ErrorMessage)
			return
		} else {
			newobj2 := model.Conversions{
				Sender:      obj.Sender,
				Amount:      obj.Amount,
				SellAsset:   obj.SellAsset,
				BuyAsset:    obj.BuyAsset,
				SellIssuer:  obj.SellIssuer,
				BuyIssuer:   obj.BuyIssuer,
				Numerator:   obj.Numerator,
				Denominator: obj.Denominator,
				TXNHash:     result2,
			}
			object.ConvertBatches(newobj2)
			log.Println("Conversion succeeded with: ", result1, result2)
		}
	}

}
