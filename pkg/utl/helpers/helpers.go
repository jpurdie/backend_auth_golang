package helpers

func StringContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

const (
	LayoutISO = "2006-01-02"
	LayoutUS  = "January 2, 2006"
)
