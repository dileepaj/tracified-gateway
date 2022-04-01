package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func InsertTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	key1, error := r.URL.Query()["Id"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'Id' is missing")
		return
	}

	key2, error := r.URL.Query()["Name"]

	if !error || len(key2[0]) < 1 {
		logrus.Error("Url Parameter 'Name' is missing")
		return
	}

	key3, error := r.URL.Query()["Address"]

	if !error || len(key3[0]) < 1 {
		logrus.Error("Url Parameter 'Address' is missing")
		return
	}

	key4, error := r.URL.Query()["Designation"]

	if !error || len(key4[0]) < 1 {
		logrus.Error("Url Parameter 'Designation' is missing")
		return
	}

	key5, error := r.URL.Query()["Specification"]

	if !error || len(key5[0]) < 1 {
		logrus.Error("Url Parameter 'Specification' is missing")
		return
	}

	Id := key1[0]
	Name := key2[0]
	Address := key3[0]
	Designation := key4[0]
	Specification := key5[0]

	var records = model.Testing{
		Id:            Id,
		Name:          Name,
		Address:       Address,
		Designation:   Designation,
		Specification: Specification,
	}
	object := dao.Connection{}
	err := object.InsertRecords(records)
	if err != nil {
		panic(err)
	}
}

func InsertModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var specs model.Testing
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&specs)
	if err != nil {
		panic(err)
	}
	var records = model.Testing{
		Id:            specs.Id,
		Name:          specs.Name,
		Address:       specs.Address,
		Designation:   specs.Designation,
		Specification: specs.Specification,
	}
	object := dao.Connection{}
	err1 := object.InsertRecords(records)
	if err1 != nil {
		panic(err)
	}

}

func GetTransactionsById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	key1, error1 := r.URL.Query()["Id"]

	if !error1 || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'Id' is missing")
		return
	}

	Id := key1[0]
	object := dao.Connection{}
	p := object.GetRecordsByID(Id)
	p.Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return data
	}).Catch(func(err error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Id Not Found in Gateway DataStore",
		}
		json.NewEncoder(w).Encode(result)
		return err
	})
	p.Await()
	// if err != nil {
	// 	panic(err)
	// }

	// result := apiModel.SubmitXDRSuccess{
	// 	Status: "Records retrieved successfully",
	// }
	// json.NewEncoder(w).Encode(result)

}

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var records model.Testing
	object := dao.Connection{}
	p := object.GetAllRecords(records)
	p.Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		logrus.Error(data)
		return data
	}).Catch(func(err error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Records Not Found in Gateway DataStore",
		}
		json.NewEncoder(w).Encode(result)
		return err
	})
	p.Await()
	// if err != nil {
	// 	panic(err)
	// }

	// json.NewEncoder(w).Encode(records)
}

func RemoveTransactionsById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logrus.Error("--------------------------------------------", r.URL.Query()["Id"])
	key1, error := r.URL.Query()["Id"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'Id' is missing")
		return
	}

	Id := key1[0]
	object := dao.Connection{}
	err := object.RemoveFromRecords(Id)
	if err != nil {
		panic(err)
	}
}

func RemoveAllTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var records model.Testing
	object := dao.Connection{}
	err := object.DeleteAllRecords(records)
	if err != nil {
		panic(err)
	}
}

func UpdateTransactionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	key1, error := r.URL.Query()["Id"]

	if !error || len(key1[0]) < 1 {
		logrus.Error("Url Parameter 'Id' is missing")
		return
	}

	key2, error := r.URL.Query()["Name"]

	if !error || len(key2[0]) < 1 {
		logrus.Error("Url Parameter 'Name' is missing")
		return
	}

	key3, error := r.URL.Query()["Address"]

	if !error || len(key3[0]) < 1 {
		logrus.Error("Url Parameter 'Address' is missing")
		return
	}

	key4, error := r.URL.Query()["Designation"]

	if !error || len(key4[0]) < 1 {
		logrus.Error("Url Parameter 'Designation' is missing")
		return
	}

	key5, error := r.URL.Query()["Specification"]

	if !error || len(key5[0]) < 1 {
		logrus.Error("Url Parameter 'Specification' is missing")
		return
	}

	Id := key1[0]
	Name := key2[0]
	Address := key3[0]
	Designation := key4[0]
	Specification := key5[0]

	var records = model.Testing{
		Id:            Id,
		Name:          Name,
		Address:       Address,
		Designation:   Designation,
		Specification: Specification,
	}
	object := dao.Connection{}
	_, err1 := object.GetRecordsByID(Id).Then(func(data interface{}) interface{} {
		selection := data.(model.Testing)
		err2 := object.UpdateRecordsById(selection, records)
		if err2 != nil {
			w.Header().Set("Content-Type", "application/json;")
			w.WriteHeader(http.StatusBadRequest)
			result := apiModel.SubmitXDRSuccess{
				Status: "Error when updating the selling status",
			}
			json.NewEncoder(w).Encode(result)
		} else {
			w.Header().Set("Content-Type", "application/json;")
			w.WriteHeader(http.StatusOK)
			result := apiModel.SubmitXDRSuccess{
				Status: "Records updated successfully",
			}
			json.NewEncoder(w).Encode(result)
		}
		return data
	}).Await()
	if err1 != nil {
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error when fetching the Record Details from Datastore or Record ID does not exists in the Datastore",
		}
		log.Println(err1)
		json.NewEncoder(w).Encode(result)
	}
}
