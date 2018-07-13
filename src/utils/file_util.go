package utils

import "os"

func CreateFile(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Create(path)
	}
}
