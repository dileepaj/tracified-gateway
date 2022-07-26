package commons

//work as ternary operator for string
func ValidateStrings(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
