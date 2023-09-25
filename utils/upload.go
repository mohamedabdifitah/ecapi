package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type ErrorResonse struct {
	StatusCode int
	Reason     error
}

func UploadFiles(base string, files ...*multipart.FileHeader) ([]string, *ErrorResonse) {
	if base == "" {
		base = filepath.Base("/assets")
	}
	var response []string
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return response, &ErrorResonse{StatusCode: 400, Reason: fmt.Errorf("File could not be opened")}
		}
		defer src.Close()
		if err = os.MkdirAll(filepath.Dir(base), 0750); err != nil {
			return response, &ErrorResonse{StatusCode: 500, Reason: err}
		}
		var extension string = strings.Split(file.Filename, ".")[1]
		var name string = base + "/" + "hello" + "." + extension
		out, err := os.Create(name)
		if err != nil {
			return response, &ErrorResonse{StatusCode: 500, Reason: err}
		}
		defer out.Close()

		_, err = io.Copy(out, src)
		if err != nil {
			return response, &ErrorResonse{StatusCode: 500, Reason: err}
		}
		response = append(response, out.Name())
	}

	return response, nil
}
