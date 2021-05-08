package util

func IndexOf(strings []string, search string) int {

	for i, v := range strings {
		if v == search {
			return i
		}
	}
	return -1
}

func Remove(strings []string, search string) []string {

	result := []string{}
	for _, v := range strings {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}

func Contains(a []int, e int) bool {
	for _, v := range a {
		if e == v {
			return true
		}
	}
	return false
}
