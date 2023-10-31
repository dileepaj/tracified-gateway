package experthelpers

func WriteGetterMethods(metadataGetter string) string {
	commentForGetter := "\n\t" + `//get value and exponent`
	getterBody := "\n\t" + `function getValues() public view returns (int256, int256) {`
	getterBody = getterBody + "\n\t\t" + `return (result.value, result.exponent);`
	getterBody = getterBody + "\n\t" + `}` + "\n"
	contractBody := commentForGetter + getterBody
	contractBody += metadataGetter

	return contractBody
}
