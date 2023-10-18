package file

import (
	"fmt"
	"io/ioutil"
	"os"
)

func Read(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("Не удалось открыть файл: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Не удалось прочитать файл: %v", err)
	}

	return string(data), nil
}
