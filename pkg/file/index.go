package file

import (
	"fmt"
	"io/ioutil"
	"os"
)

func Read(filename string) (string, error) {
	// Открываем файл для чтения
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("Не удалось открыть файл: %v", err)
	}
	defer file.Close()

	// Чтение всего файла в переменную data
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Не удалось прочитать файл: %v", err)
	}

	// Возвращаем содержимое файла в виде строки
	return string(data), nil
}
