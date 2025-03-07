package utils

func AssertEqualStrings(a string, b string) {
	if a != b {
		panic("strings are not equal!")
	}
}

func AssertNotEqualStrings(a string, b string) {
	if a == b {
		panic("strings are equal!")
	}
}

func AssertEmptyString(a string) {
	if a != "" {
		panic("String is not empty!")
	}
}

func AssertNotEmptyString(a string) {
	if a == "" {
		panic("String is empty!")
	}
}
