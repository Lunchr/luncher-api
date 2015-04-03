package storage

import (
	"path"

	"github.com/deiwin/imstor"
)

// Images is a wrapper around imstor.Storage that offers a type safe way
// to get paths for image sizes specific to this application. It also prepends
// the `images/` to all the returned paths.
type Images interface {
	ChecksumDataURL(string) (string, error)
	StoreDataURL(string) error
	PathForLarge(checksum string) (string, error)
}

type images struct {
	imstor.Storage
}

func NewImages() Images {
	sizes := []imstor.Size{
		imstor.Size{
			Name:   "large",
			Width:  800,
			Height: 400,
		},
	}
	formats := []imstor.Format{
		imstor.PNG2JPEG,
		imstor.JPEGFormat,
	}
	conf := imstor.NewConfig(sizes, formats)
	return images{
		imstor.New(conf),
	}
}

func (i images) PathForLarge(checksum string) (string, error) {
	return i.pathFor(checksum, "large")
}

func (i images) pathFor(checksum, size string) (string, error) {
	p, err := i.PathForSize(checksum, size)
	if err != nil {
		return "", err
	}
	return path.Join("images", p), nil
}
