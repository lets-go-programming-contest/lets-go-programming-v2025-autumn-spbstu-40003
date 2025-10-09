package utils

import (
	"encoding/json"
	"errors"
	"os"
)

var errOutputFile = errors.New("cannot create or open output file")

func Save(valutes []JSONValute, path string) error {
	valutesJSON, err := json.MarshalIndent(valutes, "", " ")
	if err != nil {
		return errOutputFile
	}

	const permissionCode = 0o0644

	err = os.WriteFile(path, valutesJSON, permissionCode)
	if err != nil {
		return errOutputFile
	}

	return nil
}
