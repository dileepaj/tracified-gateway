package codeGenerator

/*
Generate the code for meta data variable definitions
*/
func WriteMetaData() string {
	formulaIDInitiator := "\n\t" + `bytes32 formulaID;`
	expertPKInitiator := "\n\t" + `bytes32 expertPK;`
	tenetID := "\n\t" + `bytes32 tenetID;`

	initiatorString := formulaIDInitiator + expertPKInitiator + tenetID
	return initiatorString
}

/*
Generate the code for meta data variable setter
*/
func MetaDataSetter() string {
	functionSignature := "\n\t" + `function setMetadata(bytes32 _formulaID, bytes32 _expertPK, bytes32 _tenetID) public {`
	formulaIDSet := "\n\t\t" + `formulaID = _formulaID;`
	expertPkSet := "\n\t\t" + `expertPK = _expertPK;`
	tenetIDSet := "\n\t\t" + `tenetID = _tenetID;`
	endFunc := "\n\t" + `}` + "\n"

	setterString := functionSignature + formulaIDSet + expertPkSet + tenetIDSet + endFunc
	return setterString
}

/*
Generate the code for formula ID getter
*/
func MetaDataFormulaIDGetter() string {
	functionSign := "\n\t" + `function getFormulaID() public view returns (bytes32) {`
	returnCmd := "\n\t\t" + `return formulaID;`
	endFunc := "\n\t" + `}` + "\n"

	getterString := functionSign + returnCmd + endFunc
	return getterString
}

/*
Generate the code for expert pk getter
*/
func MetaDataExpertPKGetter() string {
	functionSign := "\n\t" + `function getExpertPK() public view returns (bytes32) {`
	returnCmd := "\n\t\t" + `return expertPK;`
	endFunc := "\n\t" + `}` + "\n"

	getterString := functionSign + returnCmd + endFunc
	return getterString
}

/*
Generate the code for tenet ID getter
*/
func MetaDataTenantIDGetter() string {
	functionSign := "\n\t" + `function getTenetID() public view returns (bytes32) {`
	returnCmd := "\n\t\t" + `return tenetID;`
	endFunc := "\n\t" + `}` + "\n"

	getterString := functionSign + returnCmd + endFunc
	return getterString
}
