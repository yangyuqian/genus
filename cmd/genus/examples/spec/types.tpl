type StringSlice []string

func (slice StringSlice) IsInclude(item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}

	return false
}
