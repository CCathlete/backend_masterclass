package u

type StringSlice []string

func (s StringSlice) Contains(val string) bool {
	for _, innerVal := range s {
		if innerVal == val {
			return true
		}
	}
	return false
}
