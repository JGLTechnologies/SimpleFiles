package SimpleFiles

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type File struct {
	name string
	lock *sync.RWMutex
	perm os.FileMode
}

const defaultPerm = 0600

func (f *File) Read() ([]byte, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return os.ReadFile(f.name)
}

func (f *File) ReadJSON(v interface{}) error {
	f.lock.RLock()
	defer f.lock.RUnlock()
	data, err := os.ReadFile(f.name)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (f *File) ReadJSONAs[T any]() (T, error) {
	var result T
	f.lock.RLock()
	defer f.lock.RUnlock()
	data, err := os.ReadFile(f.name)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)
	return result, err
}

func (f *File) ReadXML(v interface{}) error {
	f.lock.RLock()
	defer f.lock.RUnlock()
	data, err := os.ReadFile(f.name)
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}

func (f *File) ReadXMLAs[T any]() (T, error) {
	var result T
	f.lock.RLock()
	defer f.lock.RUnlock()
	data, err := os.ReadFile(f.name)
	if err != nil {
		return result, err
	}
	err = xml.Unmarshal(data, &result)
	return result, err
}

func (f *File) ReadYAML(v interface{}) error {
	f.lock.RLock()
	defer f.lock.RUnlock()
	data, err := os.ReadFile(f.name)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, v)
}

func (f *File) ReadYAMLAs[T any]() (T, error) {
	var result T
	f.lock.RLock()
	defer f.lock.RUnlock()
	data, err := os.ReadFile(f.name)
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(data, &result)
	return result, err
}

func (f *File) ReadString() (string, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	data, err := os.ReadFile(f.name)
	return string(data), err
}

func (f *File) WriteJSON(v interface{}) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return writeToFile(f.name, data, f.perm)
}

func (f *File) WriteXML(v interface{}) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	data, err := xml.Marshal(v)
	if err != nil {
		return err
	}
	return writeToFile(f.name, data, f.perm)
}

func (f *File) WriteYAML(v interface{}) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return writeToFile(f.name, data, f.perm)
}

func (f *File) WriteString(s string) error {
	return f.Write([]byte(s))
}

func (f *File) Write(b []byte) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	return writeToFile(f.name, b, f.perm)
}

func (f *File) Append(b []byte) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	file, err := os.OpenFile(f.name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, f.perm)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(b); err != nil {
		return err
	}
	return file.Sync()
}

func (f *File) AppendString(s string) error {
	return f.Append([]byte(s))
}

func (f *File) Writer() (io.WriteCloser, error) {
	f.lock.Lock()
	file, err := os.OpenFile(f.name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.perm)
	if err != nil {
		f.lock.Unlock()
		return nil, err
	}
	return &lockedWriter{file: file, unlock: f.lock.Unlock}, nil
}

func (f *File) Rename(newName string) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	newPath := filepath.Clean(newName)
	if err := os.Rename(f.name, newPath); err != nil {
		return err
	}
	f.name = newPath
	return nil
}

func New(name string, perm *os.FileMode) (*File, error) {
	cleaned := filepath.Clean(name)
	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		return nil, err
	}
	var finalPerm os.FileMode = defaultPerm
	if perm != nil {
		finalPerm = *perm
	}
	lock := &sync.RWMutex{}
	file, err := os.OpenFile(absPath, os.O_CREATE, finalPerm)
	if err != nil {
		return nil, err
	}
	file.Close()
	return &File{name: absPath, lock: lock, perm: finalPerm}, nil
}

func writeToFile(name string, data []byte, perm os.FileMode) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(data); err != nil {
		return err
	}
	return file.Sync()
}

type lockedWriter struct {
	file   *os.File
	unlock func()
}

func (lw *lockedWriter) Write(p []byte) (int, error) {
	return lw.file.Write(p)
}

func (lw *lockedWriter) Close() error {
	err := lw.file.Sync()
	if closeErr := lw.file.Close(); err == nil {
		err = closeErr
	}
	lw.unlock()
	return err
}
