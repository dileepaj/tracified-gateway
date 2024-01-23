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
	rstpromise := repo.CheckBatchExisitsInTransactions(cocstatus.BatchName, cocstatus.ProductName, cocstatus.TenantID)
	rst, err := rstpromise.Await()
	exist := rst.(bool)
	if err != nil {
		logrus.Error("Failed to check item availability", err.Error())
		result := apiModel.CommonBadResponse{
			Message: "Failed to check item availability",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
	if !exist {
		result := apiModel.CommonBadResponse{
			Message: "Item not available",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
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

func GetCOCTransferRequestbyPublicKeyandStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	publickey := r.URL.Query().Get("pubkey")
	if publickey == "" {
		json.NewEncoder(w).Encode(model.Error{Code: http.StatusBadRequest, Message: "Publickey  missing"})
		return
	}

	struserType := r.URL.Query().Get("usertype")
	userType, userTypeErr := strconv.Atoi(struserType)
	if userTypeErr != nil || userType < 1 || userType > 2 || struserType == "" {
		json.NewEncoder(w).Encode(model.Error{Code: http.StatusBadRequest, Message: "Invalid User Type"})
		return
	}
	strcocStaus := r.URL.Query().Get("coctype")
	cocStatus, converErr := strconv.Atoi(strcocStaus)
	if converErr != nil || cocStatus < 1 || cocStatus > 4 || strcocStaus == "" {
		json.NewEncoder(w).Encode(model.Error{Code: http.StatusBadRequest, Message: "Invalid COC Status"})
		return
	}
	logrus.Println("Publickkey:" + publickey)
	logrus.Println("status:", cocStatus)
	repo := dao.Connection{}
	rst, _ := repo.GetCOCTransfersbyPublickKeyandStatus(publickey, userType, cocStatus).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	DBResult, err1 := rst.([]model.COCState)
	if !err1 {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Code: 400, Message: "Unable to get Data"}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	result := apiModel.COCTransferResponse{
		Response: DBResult}
	json.NewEncoder(w).Encode(result)
	return
}
func UpdateCurrentOwnerOfCOC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	var cocUpdateOwner model.UpdateCOCOwner
	err := json.NewDecoder(r.Body).Decode(&cocUpdateOwner)
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
	rst, _ := repo.GetCOCPreviousCurrentOwner(cocUpdateOwner).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	DBResult, err1 := rst.(model.COCPreviousOwner)
	if !err1 {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Code: 400, Message: "COC Transfer request or COC owner does not exist"}
		json.NewEncoder(w).Encode(response)
		return
	}
	if DBResult.COCStatus == 4 {
		updateerr := repo.UpdateCOCOwner(cocUpdateOwner)
		if updateerr != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Error{Code: 400, Message: "Unable to update owner"}
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusOK)
		result := apiModel.CommonSuccessMessage{
			Message: "New Owner Updated Successfully"}
		json.NewEncoder(w).Encode(result)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Code: 400, Message: "Unable to update owner COC transfer not made yet"}
		json.NewEncoder(w).Encode(response)
		return
	}
}
