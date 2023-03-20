package ethereum

import "encoding/hex"

func RemoveDuplicatesInAnArray(input []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range input {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func StringToHexString(value string) string {
	hx := hex.EncodeToString([]byte(value))
	return hx
}

func HexStringToString(value string) string {
	hx, _ := hex.DecodeString(value)
	return string(hx)
}