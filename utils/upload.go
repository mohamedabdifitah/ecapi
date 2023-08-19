package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type ErrorResonse struct {
	StatusCode int
	Reason     error
}

func UploadFile(file *multipart.FileHeader) (*http.Response, *ErrorResonse) {
	src, err := file.Open()
	if err != nil {
		return nil, &ErrorResonse{StatusCode: 400, Reason: fmt.Errorf("File could not be opened")}
	}
	defer src.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreatePart(file.Header)
	if err != nil {
		return nil, &ErrorResonse{StatusCode: 500, Reason: fmt.Errorf("Error copying file , please try again")}
	}
	_, err = io.Copy(part, src)
	if err != nil {
		return nil, &ErrorResonse{StatusCode: 400, Reason: fmt.Errorf("this file is broken")}
	}
	err = writer.Close()
	if err != nil {
		return nil, &ErrorResonse{StatusCode: 500, Reason: fmt.Errorf("there is error , please try again")}
	}
	req, err := http.NewRequest("POST", os.Getenv("UPLOAD_URL"), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return nil, &ErrorResonse{StatusCode: 503, Reason: fmt.Errorf("there is error , please try again")}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &ErrorResonse{StatusCode: res.StatusCode, Reason: err}
	}
	return res, nil
}

func UploadPhotos(files []*multipart.FileHeader) (*http.Response, *ErrorResonse) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return nil, &ErrorResonse{StatusCode: 400, Reason: fmt.Errorf("File could not be opened")}
		}
		defer src.Close()
		part, err := writer.CreatePart(file.Header)
		if err != nil {
			return nil, &ErrorResonse{StatusCode: 500, Reason: fmt.Errorf("Error copying file , please try again")}
		}
		_, err = io.Copy(part, src)
		if err != nil {
			return nil, &ErrorResonse{StatusCode: 400, Reason: fmt.Errorf("this file is broken")}
		}
		_, err = io.Copy(part, src)
		if err != nil {
			return nil, &ErrorResonse{StatusCode: 400, Reason: fmt.Errorf("this file is broken")}
		}
	}
	err := writer.Close()
	if err != nil {
		return nil, &ErrorResonse{StatusCode: 500, Reason: fmt.Errorf("there is error , please try again")}
	}
	req, err := http.NewRequest("POST", os.Getenv("UPLOAD_URL")+"/photos", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return nil, &ErrorResonse{StatusCode: 503, Reason: fmt.Errorf("there is error , please try again")}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &ErrorResonse{StatusCode: res.StatusCode, Reason: err}
	}
	return res, nil
}
