package gofile

import (
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/gaigals/gotags"
	"github.com/spf13/cast"
)

const (
	tagKeyPath       = "path"
	tagKeyTmp        = "tmp"
	tagKeyContent    = "cnt"
	tagKeyPermission = "prem"
)

var (
	keyPath       = gotags.NewKey(tagKeyPath, false, true, nil, reflect.Struct)
	keyTmp        = gotags.NewKey(tagKeyTmp, true, false, nil, reflect.Struct)
	keyContent    = gotags.NewKey(tagKeyContent, false, false, nil, reflect.Struct)
	keyPermission = gotags.NewKey(tagKeyPermission, false, false, ValidatePrems, reflect.Struct)
)

var (
	tagSettings = gotags.NewTagSettingsDefault(
		"file",
		TagProcessor,
		keyPath,
		keyTmp,
		keyContent,
		keyPermission,
	)
)

type FileHandle struct {
	Name     string
	Path     string
	FullPath string
}

func ValidatePrems(value string) error {
	_, err := cast.ToUint32E(value)
	if err != nil {
		return fmt.Errorf("tag '%s' must be uint32(os.FileMod)",
			tagKeyPermission)
	}

	return nil
}

func TagProcessor(field gotags.FieldData) error {
	// Casting validation happens at ValidatePrems().
	prem := os.FileMode(cast.ToUint32(field.KeyValue(tagKeyPermission)))
	if prem == 0 {
		prem = os.ModePerm
	}

	fileHandle, err := newFile(
		field.KeyValue(tagKeyPath),
		field.KeyValue(tagKeyContent),
		prem,
		field.HasKey(tagKeyTmp),
	)
	if err != nil {
		return err
	}

	return field.ApplySelfValue(*fileHandle)
}

func NewFileHandleStruct(data any) error {
	_, err := tagSettings.FieldData(data)
	if err != nil {
		return err
	}

	return nil
}

func newFile(fullPath, content string, prems os.FileMode, isTemp bool) (*FileHandle, error) {
	if isTemp {
		return NewFileTemp(fullPath, content, prems)
	}

	return NewFile(fullPath, content, prems)
}

func NewFileTemp(fullPath, content string, prems os.FileMode) (*FileHandle, error) {
	return NewFile(path.Join(os.TempDir(), fullPath), content, prems)
}

func NewFile(fullPath, content string, prems os.FileMode) (*FileHandle, error) {
	fileHandle := FileHandle{FullPath: fullPath}
	fileHandle.parsePath()

	isValid, err := fileHandle.isPathValid()
	if err != nil {
		return nil, err
	}

	if !isValid {
		err := fileHandle.createPath(prems)
		if err != nil {
			return nil, err
		}
	}

	err = fileHandle.Write(content, prems)
	if err != nil {
		return nil, err
	}

	return &fileHandle, nil
}

func (fh *FileHandle) parsePath() {
	fh.Path, fh.Name = path.Split(fh.FullPath)
}

func (fh *FileHandle) isPathValid() (exist bool, err error) {
	_, err = os.Stat(fh.Path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func (fh *FileHandle) createPath(prems os.FileMode) error {
	return os.MkdirAll(fh.Path, prems)
}

func (fh *FileHandle) Write(content string, prems os.FileMode) error {
	flags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC

	// #nosec G304 Potential file inclusion via variable
	file, err := os.OpenFile(fh.FullPath, flags, prems)
	if err != nil {
		return fmt.Errorf("failed to open '%s', error: %v", fh.FullPath, err)
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("failed to write into '%s', error: %v", fh.FullPath, err)
	}

	return nil
}

func (fh *FileHandle) Read() (string, error) {
	// #nosec G304 Potential file inclusion via variable
	content, err := os.ReadFile(fh.FullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file '%s', error: %v", fh.FullPath, err)
	}

	return string(content), nil
}

func (fh *FileHandle) Remove() error {
	return os.Remove(fh.FullPath)
}
