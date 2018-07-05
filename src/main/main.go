package main

import (
	"main/proofs/interpreter"
)

func main() {
	//create := executer.CreateAccount()
	// newRootHash := executer.InsertDataHash("E3FC18CB4776193F8AD15A947406DBYE", "SDL26B3CQN4AQHPV3MDRMUB5BXNMCQLHY3HVAD7ZOP4QACX2OL7V2IOW", "001", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee13311ee6e")

	// router := routes.NewRouter()

	// fmt.Println("Latest Root hash")
	// fmt.Println(newRootHash)

	// log.Fatal(http.ListenAndServe(":8030", router))

	interpreter.InterpretPOC("rooothash", "treee")

	// x, _ := json.Marshal(blockchainTree())
	// fmt.Println(string(x))

	// 	fmt.Println(len(x.Previous))

	// 	for i := 0; i < len(x.Previous); i++ {
	// 		if db.Previous[0] != nil {
	// 			if db.Current.TDPID == bc.Current.TDPID && db.Current.DataHash == bc.Current.DataHash {
	// 				//check whether there a re multiple pre?
	// 				//traverse through the length.array
	// 				compare(db.Previous[i], bc.Previous[i])
	// 				if seq != "" && db.Sequence == db.Sequence {
	// 					return true
	// 				} else {
	// 					return false
	// 				}
	// 			} else {
	// 				return false
	// 			}

	// 		} else {
	// 			return false
	// 		}
	// 	}

	// interpreter.InterpretPOC("rootHash", "treeObj")

}

// func blockchainTree() model.Node {
// 	cur1 := model.Current{"001", "cur1datahash", "cur1txn", "1000"}
// 	cur20 := model.Current{"0020", "cur20datahash", "cur20txn", "2000"}
// 	cur21 := model.Current{"0021", "cur21datahash", "cur21txn", "21000"}
// 	cur4 := model.Current{"004", "cur4datahash", "cur4txn", "4000"}
// 	nullCurr := model.Current{}

// 	preArr1 := []model.Node{}
// 	pre1 := model.Node{preArr1, cur1}

// 	preArr20 := []model.Node{pre1}
// 	pre20 := model.Node{preArr20, cur20}

// 	preArr21 := []model.Node{}
// 	pre21 := model.Node{preArr21, cur21}

// 	preArr3 := []model.Node{pre20, pre21}
// 	pre3 := model.Node{preArr3, nullCurr}

// 	preArr4 := []model.Node{pre3}
// 	pre4 := model.Node{preArr4, cur4}

// 	return pre4
// }
