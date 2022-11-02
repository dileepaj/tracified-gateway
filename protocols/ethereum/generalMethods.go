package ethereum

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