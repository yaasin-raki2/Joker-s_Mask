package jsm

import (
	"os"
)

func (j *Jsm) CreateDirIfNotExist(path string) error {
	const mod = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, mod)
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *Jsm) CreateFileIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}

		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}
	return nil
}
