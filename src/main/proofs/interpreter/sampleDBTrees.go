package interpreter

import (
	
	"main/model"
)

func dBTree() model.Node {
	cur1 := model.Current{"001", "E3FC18CB4776193F8AD15A947406DBYE", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee13311etxn"}

	cur2 := model.Current{"002", "E3FC18CB4776193F8AD15A947402000", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	cur3 := model.Current{"003", "E3FC18CB4776193F8AD15A947402000", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	cur4 := model.Current{"004", "E3FC18CB4776193F8AD15A947402000", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	preArr1 := []model.Node{}
	preArr2 := []model.Node{}
	preArr3 := []model.Node{}
	preArr4 := []model.Node{}

	pre1 := model.Node{preArr1, cur1, "1000500"}
	preArr2 = []model.Node{pre1}

	pre2 := model.Node{preArr2, cur2, "1000501"}
	preArr3 = []model.Node{pre2}

	pre3 := model.Node{preArr3, cur3, "1000502"}
	preArr4 = []model.Node{pre3}

	pre4 := model.Node{preArr4, cur4, "1000503"}
	return pre4
}

func dBTreeFalse() model.Node {
	cur1 := model.Current{"001", "E3FC18CB47761dssdsd93F8AD15A947406DBYE", "ERROR5c9845sdsdsdsdsdb218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee13311etxn"}

	cur2 := model.Current{"002", "E3FC18CB4776193F8AD15A947402000", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	cur3 := model.Current{"003", "E3FC18CB4776193F8AD15A947402000", "cda5c984ghjg5b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	cur4 := model.Current{"004", "E3FC18CB4776193F8AD15A947402000", "cda5c984ghjg5b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	preArr1 := []model.Node{}
	preArr2 := []model.Node{}
	preArr3 := []model.Node{}
	preArr4 := []model.Node{}

	pre1 := model.Node{preArr1, cur1, "1000500"}
	preArr2 = []model.Node{pre1}

	pre2 := model.Node{preArr2, cur2, "1000501"}
	preArr3 = []model.Node{pre2}

	pre3 := model.Node{preArr3, cur3, "1000502"}
	preArr4 = []model.Node{pre3}

	pre4 := model.Node{preArr4, cur4, "1000503"}
	return pre4
}

func dBTree2() model.Node {
	cur1 := model.Current{"001", "E3FC18CB4776193F8AD15A947406DBYE", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee13311etxn"}

	cur2 := model.Current{"002", "E3FC18CB4776193F8AD15A947402000", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	cur3 := model.Current{}

	cur4 := model.Current{"003", "E3FC18CB4776193F8AD15A947402000", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	preArr1 := []model.Node{}
	preArr2 := []model.Node{}
	preArr3 := []model.Node{}
	preArr4 := []model.Node{}

	pre1 := model.Node{preArr1, cur1, "1000500"}
	preArr2 = []model.Node{pre1}
	pre2 := model.Node{preArr2, cur2, "1000501"}

	preArr3 = []model.Node{pre2, pre1}

	pre3 := model.Node{preArr3, cur3, "1000502"}

	preArr4 = []model.Node{pre3}

	pre4 := model.Node{preArr4, cur4, "1000503"}	
	return pre4
}

func dBTree2False() model.Node {
	cur1 := model.Current{"001", "E3FC18CB4776193F8AD15A947406DBYE", "Error5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee13311etxn"}

	cur2 := model.Current{"002", "E3FC18CB4776193F8AghfghD15A947402000", "cda5c9845ghfhb218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	cur3 := model.Current{}

	cur4 := model.Current{"003", "E3FC18CB4776193F8AD15A947402000", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	preArr1 := []model.Node{}
	preArr2 := []model.Node{}
	preArr3 := []model.Node{}
	preArr4 := []model.Node{}

	pre1 := model.Node{preArr1, cur1, "1000500"}
	preArr2 = []model.Node{pre1}
	pre2 := model.Node{preArr2, cur2, "1000501"}

	preArr3 = []model.Node{pre2, pre1}

	pre3 := model.Node{preArr3, cur3, "1000502"}

	preArr4 = []model.Node{pre3}

	pre4 := model.Node{preArr4, cur4, "1000503"}	
	return pre4
}