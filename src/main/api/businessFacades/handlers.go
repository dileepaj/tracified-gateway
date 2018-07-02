package businessfacades

import (
	"encoding/json"
	"fmt"
	"net/http"

	// . "main/model"
	"main/proofs/proofBuilder"

	"github.com/gorilla/mux"
)

func SaveDataHash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	result := proofBuilder.InsertTDP(vars["hash"], vars["secret"], vars["profileId"], vars["rootHash"])

	//test case
	// err1 := Error1{Code: 0, Message: "no root found"}
	// result := RootTree{Hash: "", Error: err1}

	//log the results
	fmt.Println(result, "result!!!\n")

	if result.Hash != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		switch result.Error.Code {
		case 0:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "No root"})
		case 1:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"})
		default:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"})
		}

	}

}

// func Index(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "Welcome!\n")
// }

// func TodoIndex(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(http.StatusOK)
// 	// if err := json.NewEncoder(w).Encode(todos); err != nil {
// 	// 	panic(err)
// 	// }
// }

// func TodoShow(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	var todoId int
// 	var err error
// 	if todoId, err = strconv.Atoi(vars["todoId"]); err != nil {
// 		panic(err)
// 	}
// 	todo := stellarexecuter.RepoFindTodo(todoId)
// 	if todo.Id > 0 {
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(http.StatusOK)
// 		if err := json.NewEncoder(w).Encode(todo); err != nil {
// 			panic(err)
// 		}
// 		return
// 	}

// 	// If we didn't find it, 404
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(http.StatusNotFound)
// 	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
// 		panic(err)
// 	}

// }

// /*
// Test with this curl command:

// curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/todos

// */
// func TodoCreate(w http.ResponseWriter, r *http.Request) {
// 	var todo proofBuilder.Todo
// 	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
// 	if err != nil {
// 		panic(err)
// 	}
// 	if err := r.Body.Close(); err != nil {
// 		panic(err)
// 	}
// 	if err := json.Unmarshal(body, &todo); err != nil {
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(422) // unprocessable entity
// 		if err := json.NewEncoder(w).Encode(err); err != nil {
// 			panic(err)
// 		}
// 	}

// 	t := stellarexecuter.RepoCreateTodo(todo)
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(http.StatusCreated)
// 	if err := json.NewEncoder(w).Encode(t); err != nil {
// 		panic(err)
// 	}
// }
