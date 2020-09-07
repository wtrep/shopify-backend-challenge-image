package image

import (
	"github.com/google/uuid"
	"os"
)

type Image struct {
	UUID       uuid.UUID
	Name       string
	Owner      string
	Extension  string
	Height     int32
	Length     int32
	Bucket     string
	BucketPath string
	Status     string
}

// Convert a CreateImageRequest into an Image object
func (i CreateImageRequest) toImage(owner string) Image {
	uuidToCreate := uuid.New()
	return Image{
		UUID:       uuidToCreate,
		Name:       i.Name,
		Owner:      owner,
		Extension:  i.Extension,
		Height:     i.Height,
		Length:     i.Length,
		Bucket:     os.Getenv("BUCKET"),
		BucketPath: uuidToCreate.String() + "." + i.Extension,
		Status:     "CREATED",
	}
}

// Convert an Image into a LinkedImageResponse object
func (i Image) toLinkedImageResponse(url string) LinkedImageResponse {
	return LinkedImageResponse{
		Uuid:      i.UUID.String(),
		Name:      i.Name,
		Url:       url,
		Owner:     i.Owner,
		Extension: i.Extension,
		Height:    i.Height,
		Length:    i.Length,
	}
}

// Convert an Image into an UnlinkedImageResponse object
func (i Image) toUnlinkedImageResponse() UnlinkedImageResponse {
	return UnlinkedImageResponse{
		Uuid:      i.UUID.String(),
		Name:      i.Name,
		Owner:     i.Owner,
		Extension: i.Extension,
		Height:    i.Height,
		Length:    i.Length,
	}
}

// Convert an Image into a CreateImageResponse object
func (i Image) toCreateImageResponse() CreateImageResponse {
	return CreateImageResponse{
		Uuid:      i.UUID.String(),
		Name:      i.Name,
		Owner:     i.Owner,
		Extension: i.Extension,
		Height:    i.Height,
		Length:    i.Length,
	}
}

// Convert an array of Image into an array of UnlinkedImageResponse object
func imagesToUnlinkedImagesReponse(images []Image) []UnlinkedImageResponse {
	response := make([]UnlinkedImageResponse, 0)
	for _, i := range images {
		response = append(response, i.toUnlinkedImageResponse())
	}
	return response
}
