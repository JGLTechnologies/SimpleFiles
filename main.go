package SimpleFiles

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

type FileWriter struct {
	File *File
}

func (fw *FileWriter) Write(p []byte) (n int, err error) {
	return len(p), fw.File.Write(p)
}

type File struct {
	name string
	Lock *sync.RWMutex
}

func (f *File) Read() ([]byte, error) {
	f.Lock.RLock()
	defer f.Lock.RUnlock()
	return ioutil.ReadFile(f.name)
}

func (f *File) ReadJSON(object interface{}) error {
	f.Lock.RLock()
	defer f.Lock.RUnlock()
	data, err := ioutil.ReadFile(f.name)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &object)
}

func (f *File) ReadXML(object interface{}) error {
	f.Lock.RLock()
	defer f.Lock.RUnlock()
	data, err := ioutil.ReadFile(f.name)
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, &object)
}

func (f *File) ReadString() (string, error) {
	f.Lock.RLock()
	defer f.Lock.RUnlock()
	data, err := ioutil.ReadFile(f.name)
	return string(data), err
}

func (f *File) WriteJSON(object interface{}) error {
	f.Lock.Lock()
	defer f.Lock.Unlock()
	data, jsonErr := json.Marshal(object)
	if jsonErr != nil {
		return jsonErr
	}
	writeErr := ioutil.WriteFile(f.name, data, 0644)
	return writeErr
}

func (f *File) WriteXML(object interface{}) error {
	f.Lock.Lock()
	defer f.Lock.Unlock()
	data, jsonErr := xml.Marshal(object)
	if jsonErr != nil {
		return jsonErr
	}
	writeErr := ioutil.WriteFile(f.name, data, 0644)
	return writeErr
}

func (f *File) WriteString(s string) error {
	f.Lock.Lock()
	defer f.Lock.Unlock()
	err := ioutil.WriteFile(f.name, []byte(s), 0644)
	return err
}

func (f *File) Write(b []byte) error {
	f.Lock.Lock()
	defer f.Lock.Unlock()
	err := ioutil.WriteFile(f.name, b, 0644)
	return err
}

func (f *File) Writer() io.Writer {
	return &FileWriter{File: f}
}

func New(name string) (*File, error) {
	lock := &sync.RWMutex{}
	file, err := os.OpenFile(name, os.O_CREATE, 0644)
	file.Close()
	if err != nil {
		return &File{}, err
	}
	return &File{name, lock}, nil
}
