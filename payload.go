package walgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/url"
	"sync"
)

const (
	formUrlEncodedContentType = "application/x-www-form-urlencoded"
	octetStreamContentType    = "application/octet-stream"
)

type payload struct {
	contentType string
	data        []byte
}

type MultipartPayload struct {
	values map[string]string
	files  map[string]FormFile

	lock sync.Mutex
}

type FormFile struct {
	File io.ReadCloser
	Name string
}

func (p *payload) getContentType() (t string) {
	return p.contentType
}

func (p *payload) getData() (d []byte) {
	return p.data
}

func payloadFromValues(v url.Values) (p *payload) {
	return &payload{
		contentType: formUrlEncodedContentType,
		data:        []byte(v.Encode()),
	}
}

func createJsonPayload(v interface{}) (p *payload, err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	} else {
		return &payload{data: data}, nil
	}
}

func payloadFromRawData(d []byte) (p *payload) {
	return &payload{data: d, contentType: octetStreamContentType}
}

func payloadFromMultipart(m *MultipartPayload) (p *payload, err error) {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	defer writer.Close()

	for k, v := range m.values {
		err := writer.WriteField(k, v)
		if err != nil {
			return nil, err
		}
	}

	for k, v := range m.files {
		if file, err := writer.CreateFormFile(k, v.Name); err == nil {
			defer v.File.Close()
			_, err2 := io.Copy(file, v.File)
			if err2 != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	contentType := writer.FormDataContentType()
	writer.Close()

	return &payload{
		contentType: contentType,
		data:        buffer.Bytes(),
	}, nil
}

func (p *MultipartPayload) hasName(name string) bool {
	for k, _ := range p.values {
		if k == name {
			return true
		}
	}

	for k, _ := range p.files {
		if k == name {
			return true
		}
	}

	return false
}

func (p *MultipartPayload) Add(name, value string) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.values == nil {
		p.values = make(map[string]string)
	}

	if p.hasName(name) {
		return errors.New("Name has already been used.")
	}

	p.values[name] = value
	return nil
}

func (p *MultipartPayload) AddFile(name string, f FormFile) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.files == nil {
		p.files = make(map[string]FormFile)
	}

	if p.hasName(name) {
		return errors.New("Name has already been used.")
	}

	p.files[name] = f
	return nil
}
