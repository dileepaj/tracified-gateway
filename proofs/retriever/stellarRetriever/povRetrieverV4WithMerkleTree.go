package stellarRetriever

func (poc *POCTreeV4) ConstructPOCMerkleTree() {
	poc.generatePOCV4WithMerkleTree()
	poc.updateSiblings()
}

func (poc *POCTreeV4) generatePOCV4WithMerkleTree() {
	// initialize tree
	if poc.Nodes == nil {
		poc.Nodes = make(map[string]*POCNode)
	}
	if poc.siblings == nil {
		poc.siblings = make(map[string][]string)
	}

	if poc.Level == 0 {
		poc.LastTxnHash = poc.TxnHash
	} else {
		if poc.LastTxnHash == "" {
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
			node.TrustLinks = append(node.TrustLinks, gtxe.TxnHash)
		}

	} else {
		poc.Nodes[gtxe.CurrentTxnHash] = &POCNode{
			Id:         ctxe.TxnHash,
			Data:       *ctxe,
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
			}
			if poc.Nodes[pgtxe.CurrentTxnHash] == nil {
				poc.Nodes[pgtxe.CurrentTxnHash] = &POCNode{
					Id:         pgtxe.TxnHash,
					Data:       *pgtxe,
					TrustLinks: []string{pgtxe.TxnHash},
				}
			}
			if !contains(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash) {
				poc.Nodes[pgtxe.CurrentTxnHash].Children = append(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash)
			}
		}
		poc.LastTxnHash = gtxe.PreviousTxnHash
		poc.generatePOCV4WithMerkleTree()
		break
	case "5":
		poc.LastTxnHash = gtxe.PreviousTxnHash
		// create backLink with splitParent
		p := ConcreteStellarTransaction{Txnhash: gtxe.PreviousTxnHash}
		pspgtxe, err2 := p.GetTransactionCollection()
		if err2 == nil {
			if poc.Nodes[pspgtxe.CurrentTxnHash] == nil {
				poc.Nodes[pspgtxe.CurrentTxnHash] = &POCNode{
					Id:         pspgtxe.TxnHash,
					Data:       *pspgtxe,
					TrustLinks: []string{pspgtxe.TxnHash},
				}
			}
			if !contains(poc.Nodes[pspgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash) {
				poc.Nodes[pspgtxe.CurrentTxnHash].Children = append(poc.Nodes[pspgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash)
			}
			if !contains(poc.Nodes[gtxe.CurrentTxnHash].Parents,pspgtxe.CurrentTxnHash) {
				poc.Nodes[gtxe.CurrentTxnHash].Parents = append(poc.Nodes[gtxe.CurrentTxnHash].Parents, pspgtxe.CurrentTxnHash)
			}
		}
		poc.generatePOCV4WithMerkleTree()
		break
	case "6":
		// maintain siblings globally
		_, exists := poc.siblings[gtxe.PreviousTxnHash]
		if !exists {
			poc.siblings[gtxe.PreviousTxnHash] = []string{gtxe.CurrentTxnHash}
		} else if !contains(poc.siblings[gtxe.PreviousTxnHash], gtxe.CurrentTxnHash) {
			poc.siblings[gtxe.PreviousTxnHash] = append(poc.siblings[gtxe.PreviousTxnHash], gtxe.CurrentTxnHash)
		}
		sp := ConcreteStellarTransaction{Txnhash: gtxe.PreviousTxnHash}
		// backlinks build transaction hash
		spgtxe1, err1 := sp.GetTransactionCollection()
		if err1 != nil {
			return
		}
		if err1 == nil {
			e := ConcreteStellarTransaction{Txnhash: spgtxe1.CurrentTxnHash}
			// current transaction hash of the backlinks build transaction hash
			spgtxe2, err := e.GetTransactionCollection()
			if err != nil {
				return
			}
			if poc.Nodes[spgtxe1.CurrentTxnHash] == nil {
				poc.Nodes[spgtxe1.CurrentTxnHash] = &POCNode{
					Id:         spgtxe1.CurrentTxnHash,
					Data:       *spgtxe2,
					TrustLinks: []string{spgtxe1.TxnHash},
				}
				poc.Nodes[spgtxe1.CurrentTxnHash].Data.TxnHash = spgtxe1.CurrentTxnHash
			}
			if !contains(poc.Nodes[gtxe.CurrentTxnHash].Parents, spgtxe1.CurrentTxnHash) {
				poc.Nodes[gtxe.CurrentTxnHash].Parents = append(poc.Nodes[gtxe.CurrentTxnHash].Parents, spgtxe1.CurrentTxnHash)
			}
			if !contains(poc.Nodes[spgtxe1.CurrentTxnHash].Children, gtxe.CurrentTxnHash) {
				poc.Nodes[spgtxe1.CurrentTxnHash].Children = append(poc.Nodes[spgtxe1.CurrentTxnHash].Children, gtxe.CurrentTxnHash)
			}
			p := ConcreteStellarTransaction{Txnhash: spgtxe1.PreviousTxnHash}
			spgtxe, err2 := p.GetTransactionCollection()
			if err2 == nil {
				if !contains(poc.Nodes[gtxe.CurrentTxnHash].Parents, spgtxe.CurrentTxnHash) {
					poc.Nodes[spgtxe1.CurrentTxnHash].Parents = append(poc.Nodes[spgtxe1.CurrentTxnHash].Parents, spgtxe.CurrentTxnHash)
				}
				if poc.Nodes[spgtxe.CurrentTxnHash] == nil {
					poc.Nodes[spgtxe.CurrentTxnHash] = &POCNode{
						Id:         spgtxe.TxnHash,
						Data:       *spgtxe,
						TrustLinks: []string{spgtxe.TxnHash},
					}
				}
				if !contains(poc.Nodes[spgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash) {
					poc.Nodes[spgtxe.CurrentTxnHash].Children = append(poc.Nodes[spgtxe.CurrentTxnHash].Children, spgtxe1.CurrentTxnHash)
				}
			}
		}
		poc.LastTxnHash = spgtxe1.PreviousTxnHash
		poc.generatePOCV4WithMerkleTree()
		if spgtxe1.TxnType == "8" {
			mergeParent := spgtxe1.CurrentTxnHash
			mergeHashes := []string{spgtxe1.MergeID, spgtxe1.PreviousTxnHash}
			for _, hash := range mergeHashes {
				p := ConcreteStellarTransaction{Txnhash: hash}
				pgtxe, _ := p.GetTransactionCollection()
				if !contains(poc.Nodes[mergeParent].Parents, spgtxe1.CurrentTxnHash) {
					poc.Nodes[mergeParent].Parents = append(poc.Nodes[mergeParent].Parents, pgtxe.CurrentTxnHash)
				}
				if poc.Nodes[pgtxe.CurrentTxnHash] == nil {
					poc.Nodes[pgtxe.CurrentTxnHash] = &POCNode{Id: pgtxe.CurrentTxnHash}
				}
				if !contains(poc.Nodes[pgtxe.CurrentTxnHash].Children, mergeParent) {
					poc.Nodes[pgtxe.CurrentTxnHash].Children = append(poc.Nodes[pgtxe.CurrentTxnHash].Children, mergeParent)
				}
			}
			for _, hash := range mergeHashes {
				poc.LastTxnHash = hash
				poc.generatePOCV4WithMerkleTree()
			}
		}
		break
	case "7":
		mergeParent := gtxe.CurrentTxnHash
		mergeHashes := []string{gtxe.MergeID}
		previousTxn := gtxe.PreviousTxnHash
		for true {
			p := ConcreteStellarTransaction{Txnhash: previousTxn}
			pgtxe, _ := p.GetTransactionCollection()
			if pgtxe.TxnType == "7" && pgtxe.MergeBlock != 0 {
				mergeHashes = append(mergeHashes, pgtxe.MergeID)
				previousTxn = pgtxe.PreviousTxnHash
			} else if pgtxe.TxnType == "7" && pgtxe.MergeBlock == 0 {
				mergeHashes = append(mergeHashes, previousTxn)
				break
			} else {
				mergeHashes = append(mergeHashes, previousTxn)
				break
			}
		}
		for _, hash := range mergeHashes {
			p := ConcreteStellarTransaction{Txnhash: hash}
			pgtxe, _ := p.GetTransactionCollection()
			if !contains(poc.Nodes[mergeParent].Parents, pgtxe.CurrentTxnHash) {
				poc.Nodes[mergeParent].Parents = append(poc.Nodes[mergeParent].Parents, pgtxe.CurrentTxnHash)
			}
			if poc.Nodes[pgtxe.CurrentTxnHash] == nil {
				poc.Nodes[pgtxe.CurrentTxnHash] = &POCNode{Id: pgtxe.CurrentTxnHash}
			}
			if !contains(poc.Nodes[pgtxe.CurrentTxnHash].Children, mergeParent) {
				poc.Nodes[pgtxe.CurrentTxnHash].Children = append(poc.Nodes[pgtxe.CurrentTxnHash].Children, mergeParent)
			}
		}
		for _, hash := range mergeHashes {
			poc.LastTxnHash = hash
			poc.generatePOCV4WithMerkleTree()
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
			}
			if poc.Nodes[pgtxe.CurrentTxnHash] == nil {
				poc.Nodes[pgtxe.CurrentTxnHash] = &POCNode{Id: pgtxe.CurrentTxnHash}
			}
			if !contains(poc.Nodes[pgtxe.CurrentTxnHash].Children, mergeParent) {
				poc.Nodes[pgtxe.CurrentTxnHash].Children = append(poc.Nodes[pgtxe.CurrentTxnHash].Children, mergeParent)
			}
		}
		for _, hash := range mergeHashes {
			poc.LastTxnHash = hash
			poc.generatePOCV4WithMerkleTree()
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
			}
			if poc.Nodes[pgtxe.CurrentTxnHash] == nil {
				poc.Nodes[pgtxe.CurrentTxnHash] = &POCNode{
					Id:         pgtxe.TxnHash,
					Data:       *pgtxe,
					TrustLinks: []string{pgtxe.TxnHash},
				}
			}
			if !contains(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash) {
				poc.Nodes[pgtxe.CurrentTxnHash].Children = append(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash)
			}
		}
		poc.LastTxnHash = gtxe.PreviousTxnHash
		poc.generatePOCV4()
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
			if poc.Nodes[pgtxe.CurrentTxnHash] == nil {
				poc.Nodes[pgtxe.CurrentTxnHash] = &POCNode{
					Id:         pgtxe.TxnHash,
					Data:       *pgtxe,
					TrustLinks: []string{pgtxe.TxnHash},
				}
			}
			if !contains(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash) {
				poc.Nodes[pgtxe.CurrentTxnHash].Children = append(poc.Nodes[pgtxe.CurrentTxnHash].Children, gtxe.CurrentTxnHash)
			}
		}
		poc.LastTxnHash = gtxe.PreviousTxnHash
		poc.generatePOCV4()
		break
	}
}