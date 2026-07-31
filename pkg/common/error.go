package common

const DeleteCacheError = "Failed to delete cache error"

// HandleError handles errors by panicking
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
