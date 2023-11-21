package businessFacades

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/sirupsen/logrus"
)

func COCTransferRequestInit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	var cocstatus model.COCState
	err := json.NewDecoder(r.Body).Decode(&cocstatus)
	if err != nil {
		logrus.Error("Error decoding data from payload : ", err)
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.CommonBadResponse{
			Message: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
	if cocstatus.COCStatus != model.COC_TRANSFER_ENABLED {
		logrus.Error("Invalid status:  ", err)
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.CommonBadResponse{
			Message: "Invalid COC status",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
	repo := dao.Connection{}
	saveErr := repo.InsertCOCStatus(cocstatus)
	if saveErr != nil {
		logrus.Error("failed to save COC transfer request", err)
		result := apiModel.CommonBadResponse{
			Message: "failed to save COC transfer request",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	result := apiModel.CommonSuccessMessage{
		Message: "Success"}
	json.NewEncoder(w).Encode(result)
	return
}

func UpdateCOCTransferStatus(w http.ResponseWriter, r *http.Request) {
	logger := utilities.NewCustomLogger()
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	var cocstatus model.UpdateCOCState
	err := json.NewDecoder(r.Body).Decode(&cocstatus)
	if err != nil {
		logrus.Error("Error decoding data from payload : ", err)
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.CommonBadResponse{
			Message: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
	repo := dao.Connection{}
	logrus.Println("ID", cocstatus.Id)
	rst, err := repo.GetCurrentCOCStatus(cocstatus.Id).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil || rst == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.Error{Code: http.StatusNoContent, Message: "Error while fetching data from Tracified %s"})
		logger.LogWriter("Error while fetching data : "+err.Error(), constants.ERROR)
		return
	}
	DBResult := rst.(model.COCStateDBResponse)
	// Define error messages for specific state transitions.
	errorMessages := map[string]string{
		"1_4": "User has not accepted or rejected COC transfer",
		"3_4": "Unable to transfer declined COC request",
		"3_2": "Invalid state transition! cannot accept a declined COC",
		"4_3": "Invalid state transition! Cannot reject already transferred COC",
		"4_2": "Invalid state transition! Cannot accept already transferred COC",
	}
	logrus.Println("DB REUSLT :", DBResult.COCStatus)
	logrus.Println("Requst:", cocstatus.COCStatus)
	if msg, ok := errorMessages[strconv.Itoa(int(DBResult.COCStatus))+"_"+strconv.Itoa(int(cocstatus.COCStatus))]; ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Error{Code: http.StatusNoContent, Message: "Invalid Status: " + msg})
		return
	}
	updateError := repo.UpdateCOCState(cocstatus)
	if updateError != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.Error{Code: http.StatusBadRequest, Message: "failed to save new state " + updateError.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	result := apiModel.CommonSuccessMessage{
		Message: "Update Successful"}
	json.NewEncoder(w).Encode(result)
	return
}
