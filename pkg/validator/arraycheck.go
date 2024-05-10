package validator

func IsInArray(str string, arr []string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}
