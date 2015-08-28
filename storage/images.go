package storage

import (
	"path"

	"github.com/Lunchr/luncher-api/db/model"
	"github.com/deiwin/imstor"
)

// Images is a wrapper around imstor.Storage that offers a type safe way
// to get paths for image sizes specific to this application. It also prepends
// the `images/` to all the returned paths.
type Images interface {
	ChecksumDataURL(string) (string, error)
	StoreDataURL(string) error
	PathsFor(checksum string) (*model.OfferImagePaths, error)
	HasChecksum(checksum string) (bool, error)
}

type images struct {
	imstor.Storage
	sizeStrings []string
}

const (
	large     = "large"
	thumbnail = "thumbnail"
)

func NewImages() Images {
	// We don't really care about the widths of the images, but double the height
	// seems like a reasonable limit
	sizes := []imstor.Size{
		imstor.Size{
			Name:   large,
			Width:  2 * 400,
			Height: 400,
		},
		imstor.Size{
			Name:   thumbnail,
			Width:  2 * 60,
			Height: 60,
		},
	}
	sizeStrings := make([]string, len(sizes))
	for i, size := range sizes {
		sizeStrings[i] = size.Name
	}
	formats := []imstor.Format{
		imstor.PNG2JPEG,
		imstor.JPEGFormat,
	}
	conf := imstor.NewConfig(sizes, formats)
	return images{
		Storage:     imstor.New(conf),
		sizeStrings: sizeStrings,
	}
}

// PathsFor returns a struct holding the image paths to the various sizes of images
// stored with the provided checksum. Returns nil for an empty checksum.
func (i images) PathsFor(checksum string) (*model.OfferImagePaths, error) {
	if checksum == "" {
		return nil, nil
	}
	largePath, err := i.pathForLarge(checksum)
	if err != nil {
		return nil, err
	}
	thumbnailPath, err := i.pathForThumbnail(checksum)
	if err != nil {
		return nil, err
	}
	return &model.OfferImagePaths{
		Large:     largePath,
		Thumbnail: thumbnailPath,
	}, nil
}

func (i images) HasChecksum(checksum string) (bool, error) {
	return i.HasSizesForChecksum(checksum, i.sizeStrings)
}

func (i images) pathForThumbnail(checksum string) (string, error) {
	return i.pathFor(checksum, thumbnail)
}

func (i images) pathForLarge(checksum string) (string, error) {
	return i.pathFor(checksum, large)
}

func (i images) pathFor(checksum, size string) (string, error) {
	p, err := i.PathForSize(checksum, size)
	if err != nil {
		return "", err
	}
	return path.Join("images", p), nil
}
