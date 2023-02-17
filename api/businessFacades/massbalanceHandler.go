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

	var total = 0
	for i := 0; i < len(obj.Destination); i++ {
		amount, converterr := strconv.Atoi(obj.Destination[i].Amount)
		if converterr != nil {
			log.Println(converterr.Error())
			http.Error(w, "incorrect amount submitted", 500)
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
					log.Println(dberr.Error())
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
	//Move DB call after BC transaction occurs success fully
	er := object.InsertSplitData(obj)
	if er != nil {
		log.Println("Failed to save data")
	}

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

	var total = 0
	fmt.Println("arr size for merge:", len(obj.Sender))
	for i := 0; i < len(obj.Sender); i++ {
		amount, converterr := strconv.Atoi(obj.Sender[i].Amount)
		if converterr != nil {
			log.Println(converterr.Error())
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
		http.Error(w, "Addition of values provided is less than the limit", 500)
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
					http.Error(w, "Error occurred", 500)
					log.Println(dberr.Error())
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

	er := object.InsertMergeData(obj)
	if er != nil {
		log.Println("Failed to save data")
	}

	return
}

func Conversions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	object := dao.Connection{}
	var obj model.TokenCoversion
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&obj)
	if err != nil {
		panic(err)
	}
	result1, err1 := massbalance.SetConversion(obj.SellerSourceAccount, obj.BuyerSourceAccount, obj.ManageSellOffer, obj.ManageBuyOffer)
	if err1 != nil {
		ErrorMessage := err1.Error()
		log.Println(w, ErrorMessage)
		return
	} else {
		newobj := model.TokenCoversion{
			SellerSourceAccount: obj.SellerSourceAccount,
			BuyerSourceAccount:  obj.BuyerSourceAccount,
			ManageSellOffer:     obj.ManageSellOffer,
			ManageBuyOffer:      obj.ManageBuyOffer,
		}
		object.ConvertBatches(newobj)
		result2, err2 := massbalance.ConvertBatches(obj.SellerSourceAccount, obj.BuyerSourceAccount, obj.ManageSellOffer, obj.ManageBuyOffer)
		if err2 != nil {
			ErrorMessage := err2.Error()
			log.Println(w, ErrorMessage)
			return
		} else {
			newobj2 := model.TokenCoversion{
				SellerSourceAccount: obj.SellerSourceAccount,
				BuyerSourceAccount:  obj.BuyerSourceAccount,
				ManageSellOffer:     obj.ManageSellOffer,
				ManageBuyOffer:      obj.ManageBuyOffer,
			}
			object.ConvertBatches(newobj2)
			log.Println("Conversion succeeded with: ", result1, result2)
		}
	}

}

func SetRanges(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	dbConn := dao.Connection{}
	var obj model.MassBalancePayload
	var result []model.RangeSetResult
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&obj)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(obj.RatioData); i++ {
		var res model.RangeSetResult
		setOptionsrst, manageDataRst, err := massbalance.SetAccountLockLevel(obj.RatioData[i].ProductName, obj.RatioData[i].Userinput, obj.RatioData[i].LowerLimit, obj.RatioData[i].HigherLimit, obj.SingerAccount, obj.UserAccount)
		if err != nil {
			log.Println("err : ", err.Error())
		}
		if manageDataRst == "locked" {
			res.Status = manageDataRst
			res.ResultHash = setOptionsrst
			obj.RatioData[i].Result = "locked"
			obj.RatioData[i].ResultHash = setOptionsrst
		} else if manageDataRst == "pending" {
			res.Status = manageDataRst
			res.ResultHash = setOptionsrst
			obj.RatioData[i].Result = "pending"
			obj.RatioData[i].ResultHash = setOptionsrst
		} else {
			res.Status = "success"
			res.ResultHash = setOptionsrst
			obj.RatioData[i].Result = "success"
			obj.RatioData[i].ResultHash = setOptionsrst
		}
		result = append(result, res)

	}
	dbErr := dbConn.SetAccountRangeLevels(obj)
	if dbErr != nil {
		log.Println("failed to save transaction data: ", dbErr.Error())
	}
	log.Println("final result : ", result)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println(err)
	}
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	dbConn := dao.Connection{}
	var obj model.MassBalancePayload
	var result []model.RangeSetResult
	var newRangeData []model.Range
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&obj)
	if err != nil {
		panic(err)
	}

	//retrieve range levels from db
	ret, err1 := dbConn.GetAccountRangeLevels(obj.UserAccount.PublicKey).Then(func(data interface{}) interface{} {
		return data
	}).Await()

	if err1 != nil {
		panic(err)
	}

	rangeLevel := ret.(model.MassBalanceRangesDB)

	//unlock account
	_, err2 := massbalance.UnlockAccount(obj.SingerAccount, obj.UserAccount)

	if err2 != nil {
		panic(err2)
	}

	for i := 0; i < len(rangeLevel.RatioData); i++ {
		//check for failed transactions
		if rangeLevel.RatioData[i].Result != "success" {
			for j := 0; j < len(obj.RatioData); j++ {
				//redo the set options transaction
				if obj.RatioData[j].ProductName == rangeLevel.RatioData[i].ProductName && rangeLevel.RatioData[i].Result == "locked" {
					var res model.RangeSetResult
					setOptionsrst, manageDataRst, err := massbalance.SetAccountLockLevel(obj.RatioData[j].ProductName, obj.RatioData[j].Userinput, obj.RatioData[j].LowerLimit, obj.RatioData[j].HigherLimit, obj.SingerAccount, obj.UserAccount)
					if err != nil {
						log.Println("err : ", err.Error())
					}
					if manageDataRst == "locked" {
						res.Status = manageDataRst
						res.ResultHash = setOptionsrst
						obj.RatioData[j].Result = "locked"
						obj.RatioData[j].ResultHash = setOptionsrst
					} else if manageDataRst == "pending" {
						res.Status = manageDataRst
						res.ResultHash = setOptionsrst
						obj.RatioData[i].Result = "pending"
						obj.RatioData[i].ResultHash = setOptionsrst
					} else {
						res.Status = "success"
						res.ResultHash = setOptionsrst
						obj.RatioData[j].Result = "success"
						obj.RatioData[j].ResultHash = setOptionsrst
					}
					result = append(result, res)
					newRangeData = append(newRangeData, obj.RatioData[j])
					break
				}
			}
			if rangeLevel.RatioData[i].Result == "pending" {
				var res model.RangeSetResult
				setOptionsrst, manageDataRst, err := massbalance.SetAccountLockLevel(rangeLevel.RatioData[i].ProductName, rangeLevel.RatioData[i].Userinput, rangeLevel.RatioData[i].LowerLimit, rangeLevel.RatioData[i].HigherLimit, rangeLevel.SingerAccount, rangeLevel.UserAccount)
				if err != nil {
					log.Println("err : ", err.Error())
				}
				if manageDataRst == "locked" {
					res.Status = manageDataRst
					res.ResultHash = setOptionsrst
					rangeLevel.RatioData[i].Result = "locked"
					rangeLevel.RatioData[i].ResultHash = setOptionsrst
				} else if manageDataRst == "pending" {
					res.Status = manageDataRst
					res.ResultHash = setOptionsrst
					rangeLevel.RatioData[i].Result = "pending"
					rangeLevel.RatioData[i].ResultHash = setOptionsrst
				} else {
					res.Status = "success"
					res.ResultHash = setOptionsrst
					rangeLevel.RatioData[i].Result = "success"
					rangeLevel.RatioData[i].ResultHash = setOptionsrst
				}
				result = append(result, res)
				newRangeData = append(newRangeData, rangeLevel.RatioData[i])
			}

		} else {
			newRangeData = append(newRangeData, rangeLevel.RatioData[i])
		}
	}

	dbErr := dbConn.UpdateAccountRangeLevels(rangeLevel.Id, newRangeData)
	if dbErr != nil {
		log.Println("failed to save transaction data: ", dbErr.Error())
	}
	log.Println("final result : ", result)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println(err)
	}
}
