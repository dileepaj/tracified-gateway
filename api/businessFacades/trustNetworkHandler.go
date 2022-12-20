package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveTrustNetworkUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var trustNetworkUser model.TrustNetWorkUser
	err := json.NewDecoder(r.Body).Decode(&trustNetworkUser)
	if err != nil {
		log.Error("Failed to decode data.: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		rst := model.RSAKeySaveSuccess{Message: "Failed to Decode User data"}
		json.NewEncoder(w).Encode(rst)
		return
	}
	dbcon := dao.Connection{}
	response, err1 := dbcon.SaveTrustNetworkUser(trustNetworkUser)
	if err1 != nil {
		log.Error("Failed to save data")
		w.WriteHeader(http.StatusBadRequest)
		rst := model.RSAKeySaveSuccess{Message: "Failed to save User"}
		json.NewEncoder(w).Encode(rst)
		return
	}
	w.WriteHeader(http.StatusOK)
	rst := model.RSAKeySaveSuccess{Message: response}
	json.NewEncoder(w).Encode(rst)
}

func GetTrustNetWorkUserbyID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	dbcon := dao.Connection{}
	objID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rst := model.Error{Message: "Invalid User ID"}
		json.NewEncoder(w).Encode(rst)
		return
	}
	p := dbcon.GetTrustNetWorkUserbyID(objID)
	p.Then(func(data interface{}) interface{} {
		result := data.(model.LoggedInTrustNetworkUser)
		return result
	}).Catch(func(error error) error {
		return error
	})
	result, err1 := p.Await()
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "User does not exist"}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func GetTrustNetWorkUserbyEncryptedPW(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	dbcon := dao.Connection{}
	objID := (vars["password"])
	fmt.Println("------pw ", objID)
	p := dbcon.GetTrustNetWorkUserbyEncryptedPW(objID)
	p.Then(func(data interface{}) interface{} {
		result := data.(model.LoggedInTrustNetworkUser)
		return result
	}).Catch(func(error error) error {
		return error
	})
	result, err1 := p.Await()
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "User does not exist"}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func EndorseTrustNetworkUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var updateRequest model.AcceptUserEndorsment
	err := json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Invalid ID"}
		json.NewEncoder(w).Encode(response)
		return
	}
	dbcon := dao.Connection{}
	p := dbcon.GetTrustNetworkUserEndorsment(updateRequest.EndorseePKHash)
	p.Then(func(data interface{}) interface{} {
		result := data.(model.TrustNetWorkUser)
		return result
	}).Catch(func(error error) error {
		return error
	})
	rst, err1 := p.Await()
	if rst == nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Endorse does not exist"}
		json.NewEncoder(w).Encode(response)
		return
	}
	toUpdate := rst.(model.TrustNetWorkUser)
	var endorsmentExistFlag = false
	for _, item := range toUpdate.Endorsments {
		p := dbcon.GetTrustNetWorkUserbyID(updateRequest.EndorsmentData.UserID)
		p.Then(func(data interface{}) interface{} {
			result := data.(model.LoggedInTrustNetworkUser)
			return result
		}).Catch(func(error error) error {
			return error
		})
		_, err1 := p.Await()
		//If the endorser does not exist an error message will be sent
		if err1 != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Error{Message: "Endorser does not exist"}
			json.NewEncoder(w).Encode(response)
			return
		}
		if item.UserID == updateRequest.EndorsmentData.UserID {
			endorsmentExistFlag = true
			break
		}
	}
	//if the same endorser is attempting to endorse the same account once again an error message will be sent
	if endorsmentExistFlag {
		w.WriteHeader(http.StatusBadRequest)
		rst := model.Error{Message: "The Endorser has already endorsed this account"}
		json.NewEncoder(w).Encode(rst)
		return
	}
	toUpdate.Endorsments = append(toUpdate.Endorsments, updateRequest.EndorsmentData)

	if err1 != nil || rst == nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Failed to get user endorsment"}
		json.NewEncoder(w).Encode(response)
		return
	}
	updateError := dbcon.UpdateTrustNetworkUserEndorsment(updateRequest.EndorseePKHash, toUpdate)
	if updateError != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Failed to Update User Endorsment"}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := model.EndorsmentUpdateSuccess{Message: "Endorsment Added Successfully"}
	json.NewEncoder(w).Encode(response)

}

func ValidateTrustNetworkUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var UserLoginRequest model.UserLogin
	err := json.NewDecoder(r.Body).Decode(&UserLoginRequest)
	if err != nil {
		log.Error("Invalid login request ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Invalid login request"}
		json.NewEncoder(w).Encode(response)
		return
	}
	dbcon := dao.Connection{}
	p := dbcon.ValidateTrustNetworkUser(UserLoginRequest.Email, UserLoginRequest.Password)
	p.Then(func(data interface{}) interface{} {
		result := data.(model.LoggedInTrustNetworkUser)
		return result
	}).Catch(func(error error) error {
		return error
	})
	result, err := p.Await()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("Invalid login request DB: ", err.Error())
		response := model.Error{Message: "incorrect username or password"}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

}

func GetTrustNetworkUserEndorsmentCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	dbcon := dao.Connection{}
	objID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rst := model.Error{Message: "Invalid User ID"}
		json.NewEncoder(w).Encode(rst)
		return
	}
	p := dbcon.GetTrustNetWorkUserbyID(objID)
	p.Then(func(data interface{}) interface{} {
		result := data.(model.LoggedInTrustNetworkUser)
		return result
	}).Catch(func(error error) error {
		return error
	})
	result, err1 := p.Await()
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "User does not exist"}
		json.NewEncoder(w).Encode(response)
		return
	}
	userEndorsments := result.(model.LoggedInTrustNetworkUser)
	var userEndorsmentCountTracker model.TrustNetworkUserEndorsmentCount
	userEndorsmentCountTracker.Totalendorsements = 0
	userEndorsmentCountTracker.FullEndorsements = 0
	userEndorsmentCountTracker.PartialEndorsements = 0
	for _, item := range userEndorsments.Endorsments {
		if item.EndorsmentsStatus == "accepted-full" {
			userEndorsmentCountTracker.FullEndorsements += 1
		} else if item.EndorsmentsStatus == "accepted-partial" {
			userEndorsmentCountTracker.PartialEndorsements += 1
		}
		userEndorsmentCountTracker.Totalendorsements += 1
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userEndorsmentCountTracker)
}
func GetAllTrustNetworkUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	dbcon := dao.Connection{}
	p := dbcon.GetAllTrustNetworkUsers()
	p.Then(func(data interface{}) interface{} {
		result := data.([]model.LoggedInTrustNetworkUser)
		return result
	}).Catch(func(error error) error {
		return error
	})
	result, err1 := p.Await()
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Error retriving Trust network user information."}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var UserLoginPwdReset model.UpdateTrustNetworkUserPassword
	err := json.NewDecoder(r.Body).Decode(&UserLoginPwdReset)
	if err != nil {
		log.Error("Invalid login request ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Invalid login request"}
		json.NewEncoder(w).Encode(response)
		return
	}

	pwdCheckRst, pwderr := services.EncodeTrustnetworkResetPassword(UserLoginPwdReset.PasswordResetCode)
	if pwderr != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Reset password validation faliure"}
		json.NewEncoder(w).Encode(response)
		return
	}

	dbcon := dao.Connection{}
	p := dbcon.GetTrustNetworkUserResetPassword(UserLoginPwdReset.Email, pwdCheckRst)
	p.Then(func(data interface{}) interface{} {
		if data == nil {
			return nil
		}
		result := data.(model.TrustNetWorkUser)
		return result
	}).Catch(func(error error) error {
		return error
	})
	result, err := p.Await()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("Incorrect reset password code")
		response := model.Error{Message: "Incorrect reset password code"}
		json.NewEncoder(w).Encode(response)
		return
	}
	updateObj := result.(model.TrustNetWorkUser)
	if commons.Decrypt(updateObj.PasswordResetCode) != commons.Decrypt(pwdCheckRst) {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("Incorrect reset password code")
		response := model.Error{Message: "Incorrect reset password code2"}
		json.NewEncoder(w).Encode(response)
		return
	}

	updateError := dbcon.UpdateTrustNetworkUserPassword(updateObj.PGPPKHash, updateObj, UserLoginPwdReset.Password)
	if updateError != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Failed to Update User Endorsment"}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := model.SuccessResponse{Message: "Updated Password"}
	json.NewEncoder(w).Encode(response)
}
func Resetpassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var ResetPasswordRequest model.ResetPsswordRequest
	err := json.NewDecoder(r.Body).Decode(&ResetPasswordRequest)
	if err != nil {
		log.Error("Invalid login request ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Invalid Reset request"}
		json.NewEncoder(w).Encode(response)
		return
	}
	dbcon := dao.Connection{}
	p := dbcon.GetTrustNetWorkUserbyEmail(ResetPasswordRequest.Email)
	p.Then(func(data interface{}) interface{} {
		result := data.(model.TrustNetWorkUser)
		return result
	}).Catch(func(error error) error {
		return error
	})
	user, usererr := p.Await()
	if usererr != nil || user == nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Email provided does not exist"}
		json.NewEncoder(w).Encode(response)
		return
	}
	userData := user.(model.TrustNetWorkUser)
	newpwd, pwderr := services.GeneratePassword()
	if pwderr != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Error sending email"}
		json.NewEncoder(w).Encode(response)
		return
	}
	rest, encodeError := services.EncodeTrustnetworkResetPassword(newpwd)
	if encodeError != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Error encoding password"}
		json.NewEncoder(w).Encode(response)
		return
	}

	updateError := dbcon.UpdateTrustNetworkResetUserPassword(userData.PGPPKHash, userData, rest)
	if updateError != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Update fail"}
		json.NewEncoder(w).Encode(response)
		return
	}

	emailErr := services.SendEmail(newpwd, ResetPasswordRequest.Email)
	if emailErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "failed to send Email"}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := model.SuccessResponse{Message: "Email has been sent"}
	json.NewEncoder(w).Encode(response)
}
