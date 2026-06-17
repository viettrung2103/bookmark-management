package common

// HandleError handles errors by panicking
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
