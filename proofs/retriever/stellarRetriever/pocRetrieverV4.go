package stellarRetriever

import (
	"github.com/dileepaj/tracified-gateway/model"
)

type POCNode struct {
	Id      string
	Data    model.TransactionCollectionBody
	Parents []string
	Children []string
	Siblings []string
	TrustLinks []string
}

type POCTreeV4 struct {
	TxnHash string
	LastTxnHash string
	Level int
	Nodes map[string]*POCNode
	Siblings map[string][]string
	BackLinkParents []string
}

func (poc *POCTreeV4) ConstructPOC() {
	poc.generatePOCV4()
	poc.updateSiblings()
}

func (poc *POCTreeV4) generatePOCV4() {
	// initialize tree
	if poc.Nodes == nil {
		poc.Nodes = make(map[string]*POCNode)
	}
	if poc.Siblings == nil {
		poc.Siblings = make(map[string][]string)
	}

	if poc.Level == 0 {
		poc.LastTxnHash = poc.TxnHash
	} else {
		if poc.LastTxnHash == ""  {
			return
		}
	}
	poc.Level = poc.Level + 1
	d := ConcreteStellarTransaction{Txnhash: poc.LastTxnHash}
	gtxe, err := d.GetTransactionCollection()	
	if err != nil {
		return
	}
	e := ConcreteStellarTransaction{Txnhash: gtxe.CurrentTxnHash}
	ctxe, err := e.GetTransactionCollection()	
	if err != nil {
		return
	}

	// initialize the node
	if node, exists := poc.Nodes[gtxe.CurrentTxnHash]; exists {
		node.Id = ctxe.TxnHash
		node.Data = *ctxe

		if !contains(node.TrustLinks, gtxe.TxnHash) {
			node.TrustLinks =  append(node.TrustLinks, gtxe.TxnHash)
		}

	} else {
		poc.Nodes[gtxe.CurrentTxnHash] = &POCNode{
			Id: ctxe.TxnHash,
			Data: *ctxe,
			TrustLinks: []string{gtxe.TxnHash},
		}
	}
		
	// check type and add parents and children\
	switch gtxe.TxnType {
		case "0":

			return
		case "2":
			if gtxe.PreviousTxnHash == "" {
				return
			}
			p := ConcreteStellarTransaction{Txnhash: gtxe.PreviousTxnHash}
			pgtxe, err1 := p.GetTransactionCollection()
			if err1 == nil {
				if !contains(poc.Nodes[gtxe.CurrentTxnHash].Parents, pgtxe.CurrentTxnHash) {
					poc.Nodes[gtxe.CurrentTxnHash].Parents = append(poc.Nodes[gtxe.CurrentTxnHash].Parents, pgtxe.CurrentTxnHash)
					poc.BackLinkParents = append(poc.BackLinkParents, gtxe.PreviousTxnHash)
				}
				// if poc.Nodes[pgtxe.CurrentTxnHash] == nil {
				// 	poc.Nodes[pgtxe.CurrentTxnHash] = &POCNode{
				// 		Id: pgtxe.TxnHash,
				// 		Data: *pgtxe,
				// 		TrustLinks: []string{pgtxe.TxnHash},
				// 	}
				// }
				// if !contains(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash) {
				// 	poc.Nodes[pgtxe.CurrentTxnHash].Children = append(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash)
				// }
			}
			poc.LastTxnHash = gtxe.PreviousTxnHash
			// poc.generatePOCV4()
			break
		case "5":
			poc.LastTxnHash = gtxe.PreviousTxnHash
			// create backLink with splitParent
			p := ConcreteStellarTransaction{Txnhash: gtxe.PreviousTxnHash}
			pspgtxe, err2 := p.GetTransactionCollection()
			if err2 == nil {
				// if poc.Nodes[pspgtxe.CurrentTxnHash] == nil {
				// 	poc.Nodes[pspgtxe.CurrentTxnHash] = &POCNode{
				// 		Id: pspgtxe.TxnHash,
				// 		Data: *pspgtxe,
				// 		TrustLinks: []string{pspgtxe.TxnHash},
				// 	}
				// }
				// if !contains(poc.Nodes[pspgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash) {
				// 	poc.Nodes[pspgtxe.CurrentTxnHash].Children = append(poc.Nodes[pspgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash)
				// }
				if !contains(poc.Nodes[gtxe.CurrentTxnHash].Parents, pspgtxe.CurrentTxnHash) {
					poc.Nodes[gtxe.CurrentTxnHash].Parents = append(poc.Nodes[gtxe.CurrentTxnHash].Parents, pspgtxe.CurrentTxnHash)
					poc.BackLinkParents = append(poc.BackLinkParents, gtxe.PreviousTxnHash)
				}
			} 	
			// poc.generatePOCV4()
			break
		case "6":
			if gtxe.PreviousTxnHash == "" {
				return
			}
			p := ConcreteStellarTransaction{Txnhash: gtxe.PreviousTxnHash}
			pgtxe, err1 := p.GetTransactionCollection()
			if err1 == nil {
				if !contains(poc.Nodes[gtxe.CurrentTxnHash].Parents, pgtxe.CurrentTxnHash) {
					poc.Nodes[gtxe.CurrentTxnHash].Parents = append(poc.Nodes[gtxe.CurrentTxnHash].Parents, pgtxe.CurrentTxnHash)
					poc.BackLinkParents = append(poc.BackLinkParents, gtxe.PreviousTxnHash)
				}
				// if poc.Nodes[pgtxe.CurrentTxnHash] == nil {
				// 	poc.Nodes[pgtxe.CurrentTxnHash] = &POCNode{
				// 		Id:         pgtxe.TxnHash,
				// 		Data:       *pgtxe,
				// 		TrustLinks: []string{pgtxe.TxnHash},
				// 	}
				// }
				// if !contains(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash) {
				// 	poc.Nodes[pgtxe.CurrentTxnHash].Children = append(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash)
				// }
			}
	
			poc.LastTxnHash = gtxe.PreviousTxnHash
			// poc.generatePOCV4WithMerkleTree()
		case "7":
			mergeParent := gtxe.CurrentTxnHash
			mergeHashes := []string{gtxe.MergeID}
			previousTxn := gtxe.PreviousTxnHash
			for true {
				p := ConcreteStellarTransaction{Txnhash: previousTxn}
				pgtxe, _ := p.GetTransactionCollection()
				if pgtxe.TxnType == "7" {
					mergeHashes = append(mergeHashes, pgtxe.MergeID)
					previousTxn = pgtxe.PreviousTxnHash
				} else {
					mergeHashes = append(mergeHashes, pgtxe.MergeID, pgtxe.PreviousTxnHash)
					break
				}
			}
			for _, hash := range mergeHashes {
				p := ConcreteStellarTransaction{Txnhash: hash}
				pgtxe, _ := p.GetTransactionCollection()
				if !contains(poc.Nodes[mergeParent].Parents, pgtxe.CurrentTxnHash) {
					poc.Nodes[mergeParent].Parents = append(poc.Nodes[mergeParent].Parents,  pgtxe.CurrentTxnHash)
					poc.BackLinkParents = append(poc.BackLinkParents, hash)
				}
				// if poc.Nodes[pgtxe.CurrentTxnHash] == nil {
				// 	poc.Nodes[pgtxe.CurrentTxnHash] = &POCNode{Id: pgtxe.CurrentTxnHash}
				// }
				// if !contains(poc.Nodes[pgtxe.CurrentTxnHash].Children, mergeParent) {
				// 	poc.Nodes[pgtxe.CurrentTxnHash].Children = append(poc.Nodes[pgtxe.CurrentTxnHash].Children, mergeParent)	
				// }
			}
			for _, hash := range mergeHashes { 
				poc.LastTxnHash = hash
				poc.generatePOCV4()
			}
			break
		case "8":
			mergeParent := gtxe.CurrentTxnHash
			mergeHashes := []string{gtxe.MergeID, gtxe.PreviousTxnHash}
			for _, hash := range mergeHashes {
				p := ConcreteStellarTransaction{Txnhash: hash}
				pgtxe, _ := p.GetTransactionCollection()
				if !contains(poc.Nodes[mergeParent].Parents, pgtxe.CurrentTxnHash) {
					poc.Nodes[mergeParent].Parents = append(poc.Nodes[mergeParent].Parents, pgtxe.CurrentTxnHash)
					poc.BackLinkParents = append(poc.BackLinkParents, hash)
				}
				// if poc.Nodes[pgtxe.CurrentTxnHash] == nil {
				// 	poc.Nodes[pgtxe.CurrentTxnHash] = &POCNode{Id: pgtxe.CurrentTxnHash}
				// }
				// if !contains(poc.Nodes[pgtxe.CurrentTxnHash].Children, mergeParent) {
				// 	poc.Nodes[pgtxe.CurrentTxnHash].Children = append(poc.Nodes[pgtxe.CurrentTxnHash].Children, mergeParent)	
				// }
			}
			for _, hash := range mergeHashes { 
				poc.LastTxnHash = hash
				// poc.generatePOCV4()
			}
			break
		case "9":
			if gtxe.PreviousTxnHash == "" {
				return
			}
			p := ConcreteStellarTransaction{Txnhash: gtxe.PreviousTxnHash}
			pgtxe, err1 := p.GetTransactionCollection()
			if err1 == nil {
				if !contains(poc.Nodes[gtxe.CurrentTxnHash].Parents, pgtxe.CurrentTxnHash) {
					poc.Nodes[gtxe.CurrentTxnHash].Parents = append(poc.Nodes[gtxe.CurrentTxnHash].Parents, pgtxe.CurrentTxnHash)
					poc.BackLinkParents = append(poc.BackLinkParents, gtxe.PreviousTxnHash)
				}
				// if poc.Nodes[pgtxe.CurrentTxnHash] == nil {
				// 	poc.Nodes[pgtxe.CurrentTxnHash] = &POCNode{
				// 		Id: pgtxe.TxnHash,
				// 		Data: *pgtxe,
				// 		TrustLinks: []string{pgtxe.TxnHash},
				// 	}
				// }
				// if !contains(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash) {
				// 	poc.Nodes[pgtxe.CurrentTxnHash].Children = append(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash)
				// }
				}
			poc.LastTxnHash = gtxe.PreviousTxnHash
			// poc.generatePOCV4()
			break
		case "10":
			if gtxe.PreviousTxnHash == "" {
				return
			}
			p := ConcreteStellarTransaction{Txnhash: gtxe.PreviousTxnHash}
			pgtxe, err1 := p.GetTransactionCollection()
			if err1 == nil {
				if !contains(poc.Nodes[gtxe.CurrentTxnHash].Parents, pgtxe.CurrentTxnHash) {
					poc.Nodes[gtxe.CurrentTxnHash].Parents = append(poc.Nodes[gtxe.CurrentTxnHash].Parents, pgtxe.CurrentTxnHash)
				}
				// if poc.Nodes[pgtxe.CurrentTxnHash] == nil {
				// 	poc.Nodes[pgtxe.CurrentTxnHash] = &POCNode{
				// 		Id: pgtxe.TxnHash,
				// 		Data: *pgtxe,
				// 		TrustLinks: []string{pgtxe.TxnHash},
				// 	}
				// }
				// if !contains(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash) {
				// 	poc.Nodes[pgtxe.CurrentTxnHash].Children = append(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash)
				// }
			}
			//poc.LastTxnHash = gtxe.PreviousTxnHash
			poc.generatePOCV4()
			break
	}
}


func (poc *POCTreeV4) updateSiblings() {
	for _, v := range poc.Siblings {
		if len(v) < 2 {
			continue
		}
		for _, hash := range v {
			for _, shash := range v {
				if shash != hash {
					if poc.Nodes[hash].Siblings == nil {
						poc.Nodes[hash].Siblings = []string{shash}
					} else {
						if !contains(poc.Nodes[hash].Siblings, shash) {
							poc.Nodes[hash].Siblings = append(poc.Nodes[hash].Siblings, shash)
						}
					}
				}
			}
		}
	}
}


func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
