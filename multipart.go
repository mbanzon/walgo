package walgo

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type Form struct {
	value map[string][]string
	file  map[string][]*File
}

type File struct {
	name string
	data []byte
}

func GetLargeForm(r *http.Request) (f *Form, err error) {
	reader, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}

	return ReadLargeForm(reader)
}

func ReadLargeForm(r *multipart.Reader) (f *Form, err error) {
	// The following is (currently) a stripped down copy of the code in the
	// standard library to allow for large forms (beyond 10mb).
	form := &Form{make(map[string][]string), make(map[string][]*File)}

	for {
		p, err := r.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		name := p.FormName()
		if name == "" {
			continue
		}
		filename := p.FileName()

		var b bytes.Buffer

		if filename == "" {
			_, err := io.Copy(&b, p)
			if err != nil && err != io.EOF {
				return nil, err
			}
			form.value[name] = append(form.value[name], b.String())
			continue
		}

		ff := &File{
			name: filename,
		}
		_, err = io.Copy(&b, p)
		if err != nil && err != io.EOF {
			return nil, err
		}
		ff.data = b.Bytes()
		form.file[name] = append(form.file[name], ff)
	}

	return form, nil
}
func (f *Form) Value(name string) string {
	for k, v := range f.value {
		if k == name {
			if len(v) > 0 {
				return v[0]
			}
		}
	}

	return ""
}

func (f *Form) File(name string) (data []byte, filename string) {
	for k, v := range f.file {
		if k == name {
			if len(v) > 0 {
				f := v[0]
				return f.data, f.name
			}
		}
	}

	return nil, ""
}

func (f *Form) SaveTemporaryFile(name string) (filename string, err error) {
	data, filename := f.File(name)
	if data == nil {
		return "", fmt.Errorf("Unknown form file: %s", name)
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "engine_")
	if err != nil {
		return "", err
	}

	count, err := tmpFile.Write(data)
	if err != nil {
		return "", err
	}
	tmpFile.Close()

	if count < 1 {
		return "", fmt.Errorf("No bytes copied.")
	}

	return tmpFile.Name(), nil
}

func (f *Form) SaveTemporaryFileFromField(name string) (filename string, err error) {
	data := f.Value(name)
	if data == "" {
		return "", fmt.Errorf("Unknown form field: %s", name)
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "engine_")
	if err != nil {
		return "", err
	}

	count, err := tmpFile.WriteString(data)
	if err != nil {
		return "", err
	}
	tmpFile.Close()

	if count < 1 {
		return "", fmt.Errorf("No bytes copied.")
	}

	return tmpFile.Name(), nil
}
