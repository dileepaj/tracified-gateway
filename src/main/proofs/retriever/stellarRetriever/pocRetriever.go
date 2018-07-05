package stellarRetriever

import (
	"main/model"
)

func RetrievePOC(rootHash string) model.Node {
	poc := blockchainTree()
	return poc

}

func blockchainTree() model.Node {
	cur1 := model.Current{"001", "cur1datahash", "cur1txn", "1000"}
	cur20 := model.Current{"0020", "cur20datahash", "cur20txn", "2000"}
	cur21 := model.Current{"0021", "cur21datahash", "cur21txn", "21000"}
	cur4 := model.Current{"004", "cur4datahash", "cur4txn", "4000"}
	nullCurr := model.Current{}

	preArr1 := []model.Node{}
	pre1 := model.Node{preArr1, cur1}

	preArr20 := []model.Node{pre1}
	pre20 := model.Node{preArr20, cur20}

	preArr21 := []model.Node{}
	pre21 := model.Node{preArr21, cur21}

	preArr3 := []model.Node{pre20, pre21}
	pre3 := model.Node{preArr3, nullCurr}

	preArr4 := []model.Node{pre3}
	pre4 := model.Node{preArr4, cur4}

	return pre4
}
