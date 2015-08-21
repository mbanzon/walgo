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

type Payload interface {
	getContentType() (t string)
	getData() (d []byte)
}

type payloadImpl struct {
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

func (p *payloadImpl) getContentType() (t string) {
	return p.contentType
}

func (p *payloadImpl) getData() (d []byte) {
	return p.data
}

func PayloadFromValues(v url.Values) (p Payload) {
	return &payloadImpl{
		contentType: "application/x-www-form-urlencoded",
		data:        []byte(v.Encode()),
	}
}

func CreateJsonPayload(v interface{}) (p Payload, err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	} else {
		return &payloadImpl{data: data}, nil
	}
}

func PayloadFromMultipart(m MultipartPayload) (p Payload, err error) {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	defer writer.Close()

	for k, v := range m.values {
		if field, err := writer.CreateFormField(k); err == nil {
			_, err := field.Write([]byte(v))
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	for k, v := range m.files {
		if file, err := writer.CreateFormFile(k, v.Name); err == nil {
			defer v.File.Close()
			io.Copy(file, v.File)
		} else {
			return nil, err
		}
	}

	return &payloadImpl{
		contentType: writer.FormDataContentType(),
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
