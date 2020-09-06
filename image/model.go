package image

type CreateImageRequest struct {
	// name of the image
	Name string `json:"name,omitempty"`
	// extension of the image
	Extension string `json:"extension,omitempty"`
	Height    int32  `json:"height,omitempty"`
	Length    int32  `json:"length,omitempty"`
}

type CreateImageResponse struct {
	// unique id of the image
	Uuid string `json:"uuid,omitempty"`
	// name of the image
	Name string `json:"name,omitempty"`
	// owner of the image
	Owner string `json:"owner,omitempty"`
	// extension of the image
	Extension string `json:"extension,omitempty"`
	Height    int32  `json:"height,omitempty"`
	Length    int32  `json:"length,omitempty"`
}

type LinkedImageResponse struct {
	// unique id of the image
	Uuid string `json:"uuid,omitempty"`
	// name of the image
	Name string `json:"name,omitempty"`
	// url to the image
	Url string `json:"url,omitempty"`
	// owner of the image
	Owner string `json:"owner,omitempty"`
	// extension of the image
	Extension string `json:"extension,omitempty"`
	Height    int32  `json:"height,omitempty"`
	Length    int32  `json:"length,omitempty"`
}

type UnlinkedImageResponse struct {
	// unique id of the image
	Uuid string `json:"uuid,omitempty"`
	// name of the image
	Name string `json:"name,omitempty"`
	// owner of the image
	Owner string `json:"owner,omitempty"`
	// extension of the image
	Extension string `json:"extension,omitempty"`
	Height    int32  `json:"height,omitempty"`
	Length    int32  `json:"length,omitempty"`
}

type UnlinkedImagesResponse = []UnlinkedImageResponse
