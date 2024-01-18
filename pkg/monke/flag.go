package monke

import (
	"errors"
	"os"
	"path/filepath"
)

func IsDir(path string) (bool, error) {
	status, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := status.Mode(); {
	case mode.IsDir():
		return true, nil
	default:
		return false, nil
	}
}

func DirectoryValidator(value string) (string, error) {
	value, err := filepath.Abs(value)
	if err != nil {
		return "", err
	}
	res, err := IsDir(value)
	if err != nil {
		return "", err
	} else if !res {
		return "", errors.New(value + " is not a directory")
	}
	return value, nil
}

func ProtocolValidator(value string) (string, error) {
	switch value {
	case "http":
		fallthrough
	case "https":
		return value, nil
	default:
		return "", errors.New("unexpected value, expected one of [http, https]")
	}
}

func GitDirectoryValidator(path string) error {
	res, err := IsDir(path)
	if err != nil {
		return err
	} else if !res {
		return errors.New(path + " is not a directory")
	}
	if _, err := GitRevision(path); err != nil {
		return errors.New(path + " is not a valid git directory: " + err.Error())
	}
	return nil
}
