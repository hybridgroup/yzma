package download

// GetModel downloads a model from the specified URL to the destination path.
func GetModel(url, dest string) error {
	return get(url, dest)
}
