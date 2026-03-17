package slack

import (
	"net/url"
	"path/filepath"
	"strings"
)

func parseFile(uri string) (string, error) {
	obj, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	splits := strings.Split(obj.Path, "/")
	name, err := url.PathUnescape(splits[2])
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(obj.Path)
	return name + ext, nil
}
