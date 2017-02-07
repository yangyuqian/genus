package genus

var GoHelperFuncs = map[string]interface{}{}

func StringWithDefault(v, d string) string {
	if v == "" {
		return d
	}

	return v
}

func BoolWithDefault(v, d bool) bool {
	return v || d
}
