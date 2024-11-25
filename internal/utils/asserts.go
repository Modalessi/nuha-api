package utils

func Assert(e error, msg string) {
	if e != nil {
		panic(msg + "\n" + e.Error())
	}
}

func AssertOn(c bool, msg string) {
	if !c {
		panic(msg)
	}
}
