package datasources

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Source -
type Source interface {
	Read(args ...string) ([]byte, error)
}

// NewSource -
func NewSource(alias, url string) (Source, error) {
	// u, err := ParseSourceURL(url)
	// if err != nil {
	// 	return nil, err
	// }
	// s := &Source{
	// 	Alias: alias,
	// 	URL:   u,
	// }
	return nil, nil
}

// ParseSourceURL -
func ParseSourceURL(value string) (*url.URL, error) {
	if value == "-" {
		value = "stdin://"
	}
	value = filepath.ToSlash(value)
	// handle absolute Windows paths
	volName := ""
	if volName = filepath.VolumeName(value); volName != "" {
		// handle UNCs
		if len(volName) > 2 {
			value = "file:" + value
		} else {
			value = "file:///" + value
		}
	}
	srcURL, err := url.Parse(value)
	if err != nil {
		return nil, err
	}

	if volName != "" {
		if strings.HasPrefix(srcURL.Path, "/") && srcURL.Path[2] == ':' {
			srcURL.Path = srcURL.Path[1:]
		}
	}

	if !srcURL.IsAbs() {
		srcURL, err = absURL(value)
		if err != nil {
			return nil, err
		}
	}
	return srcURL, nil
}

func absURL(value string) (*url.URL, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrapf(err, "can't get working directory")
	}
	urlCwd := filepath.ToSlash(cwd)
	baseURL := &url.URL{
		Scheme: "file",
		Path:   urlCwd + "/",
	}
	relURL := &url.URL{
		Path: value,
	}
	resolved := baseURL.ResolveReference(relURL)
	// deal with Windows drive letters
	if !strings.HasPrefix(urlCwd, "/") && resolved.Path[2] == ':' {
		resolved.Path = resolved.Path[1:]
	}
	return resolved, nil
}
