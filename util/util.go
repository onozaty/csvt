package util

func Remove(strings []string, search string) []string {

	result := []string{}
	for _, v := range strings {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}
