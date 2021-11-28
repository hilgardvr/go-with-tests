package iteration

func Repeat(c int, char string) string {
	var repeated string
	for i := 0; i < c; i++ {
		repeated += char
	}
	return repeated
}