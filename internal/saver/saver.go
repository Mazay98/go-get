package saver

import (
	"fmt"
	"go-get/internal/parser"
	"io/ioutil"
	"mime"
	"os"
)

func SaveResponseToFile(resp *parser.HttpResponse) (string, error) {
	fileTypes, err := mime.ExtensionsByType(resp.ContentType)

	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("go-get_file%s", fileTypes[0])
	i := 0

	for {
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			break
		}

		i++
		fileName = fmt.Sprintf("go-get_file (%d)%s", i, fileTypes[0])
	}

	err = ioutil.WriteFile(fileName, resp.Body, 0777)

	if err != nil {
		return "", err
	}

	return fileName, nil
}
