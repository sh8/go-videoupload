// Package go-videoupload is wrapper for uploading video files easily
package videoupload

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//
type Video struct {
	Filename    string
	ContentType string
	Data        []byte
	Size        int
}

// Save video file with permission 0600
func (v *Video) Save(filename string) error {
	return ioutil.WriteFile(filename, v.Data, 0600)
}

// Get hash sum for creating unique file name
func (v *Video) GetHashSum() string {
	h := sha1.Sum(v.Data)
	return fmt.Sprintf("%s", h[:4])
}

// Check if the content type is the one of video
func okContentType(contentType string) bool {
	return contentType == "video/mp4" || contentType == "video/mpg" || contentType == "video/mpeg"
}

// Extruct video struct from http.Request
func Process(r *http.Request, field string) (*Video, error) {
	file, info, err := r.FormFile(field)

	if err != nil {
		return nil, err
	}

	contentType := info.Header.Get("Content-Type")

	if !okContentType(contentType) {
		return nil, errors.New(fmt.Sprintf("Wrong content type: %s", contentType))
	}

	bs, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	v := &Video{
		Filename:    info.Filename,
		ContentType: contentType,
		Data:        bs,
		Size:        len(bs),
	}

	return v, nil
}
