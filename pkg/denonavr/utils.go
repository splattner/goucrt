package denonavr

func EqualValueList(a, b []ValueLists) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.Index != b[i].Index {
			return false
		}
	}
	return true
}
