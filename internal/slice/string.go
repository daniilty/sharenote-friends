package slice

const (
	NotFoundIndex = -1
)

func ContainsString(ss []string, str string) bool {
	for i := range ss {
		if ss[i] == str {
			return true
		}
	}

	return false
}

func RemoveString(ss []string, str string) []string {
	i := FindString(ss, str)
	if i == NotFoundIndex {
		return ss
	}

	return append(ss[:i], ss[i+1:]...)
}

func FindString(ss []string, str string) int {
	for i := range ss {
		if ss[i] == str {
			return i
		}
	}

	return NotFoundIndex
}
