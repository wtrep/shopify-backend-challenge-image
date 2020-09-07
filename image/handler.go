package image

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/wtrep/shopify-backend-challenge-image/common"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	tokenLifetime = 24 * time.Hour
)

type Handler struct {
	db *sql.DB
}

// Setup the routes and handle them
func SetupAndServeRoutes() {
	CheckEnvVariables()

	db, err := NewConnectionPool()
	if err != nil {
		panic(err)
	}
	handler := Handler{db: db}

	r := mux.NewRouter()
	r.HandleFunc("/image", handler.HandlePostImage).Methods("POST")
	r.HandleFunc("/image/{uuid}", handler.HandleGetImage).Methods("GET")
	r.HandleFunc("/image/{uuid}", handler.HandleDeleteImage).Methods("DELETE")
	r.HandleFunc("/images", handler.HandleGetImages).Methods("GET")
	r.HandleFunc("/upload/{uuid}", handler.HandlePostUpload).Methods("POST")
	r.HandleFunc("/healthz", HandleHealthzProbe)
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

// Ensure that all required environment variables are set
func CheckEnvVariables() {
	env := []string{"DB_IP", "DB_PASSWORD", "DB_USERNAME", "DB_NAME", "JWT_KEY", "BUCKET",
		"GOOGLE_APPLICATION_CREDENTIALS"}
	for _, e := range env {
		_, ok := os.LookupEnv(e)
		if !ok {
			panic("fatal: environment variable " + e + " is not set")
		}
	}
}

// Respond to the Kubernetes Readiness and Liveness probes
func HandleHealthzProbe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

// Handle the API request to create an image
func (h *Handler) HandlePostImage(w http.ResponseWriter, r *http.Request) {
	var request CreateImageRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		common.RespondWithError(w, &common.InvalidUUIDError)
		return
	}

	username, errResponse := handleJWT(r)
	if errResponse != nil {
		common.RespondWithError(w, errResponse)
		return
	}

	image := request.toImage(username)
	err = CreateImage(h.db, image)
	if err != nil {
		common.RespondWithError(w, &common.DatabaseInsertionError)
		return
	}

	err = json.NewEncoder(w).Encode(image.toCreateImageResponse())
	if err != nil {
		common.RespondWithError(w, &common.JSONEncoderError)
	}
}

// Handle the API request to get an image
func (h *Handler) HandleGetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	uuidToGet, err := uuid.Parse(vars["uuid"])
	if err != nil {
		common.RespondWithError(w, &common.InvalidUUIDError)
		return
	}

	username, errResponse := handleJWT(r)
	if errResponse != nil {
		common.RespondWithError(w, errResponse)
		return
	}

	image, err := GetImage(h.db, uuidToGet)
	if err != nil {
		common.RespondWithError(w, &common.ImageNotFoundError)
		return
	}

	if username == image.Owner && image.Status == "UPLOADED" {
		url, err := generateSignedURL(image.Bucket, image.BucketPath)
		if err != nil {
			common.RespondWithError(w, &common.URLGenerationError)
			return
		}

		response := image.toLinkedImageResponse(url)
		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			common.RespondWithError(w, &common.JSONEncoderError)
		}
	} else if image.Status == "CREATED" {
		common.RespondWithError(w, &common.ImageNotUploadedError)
	} else {
		common.RespondWithError(w, &common.UserPermissionDeniedError)
	}
}

// Handle the API request to upload an image
func (h *Handler) HandlePostUpload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	uuidToUpload, err := uuid.Parse(vars["uuid"])
	if err != nil {
		common.RespondWithError(w, &common.InvalidUUIDError)
		return
	}

	image, err := GetImage(h.db, uuidToUpload)
	if err != nil {
		common.RespondWithError(w, &common.ImageNotFoundError)
		return
	}

	username, detailedErr := handleJWT(r)
	if detailedErr != nil {
		common.RespondWithError(w, detailedErr)
		return
	}
	if image.Owner != username {
		common.RespondWithError(w, &common.WrongUserError)
		return
	}

	file, detailedErr := getImageFromForm(w, r)
	if detailedErr != nil {
		common.RespondWithError(w, detailedErr)
		return
	}
	defer file.Close()

	err = uploadToBucket(file, image.Bucket, image.BucketPath)
	if err != nil {
		common.RespondWithError(w, &common.FileUploadError)
		return
	}

	image.Status = "UPLOADED"
	err = UpdateImage(h.db, *image)
	if err != nil {
		common.RespondWithError(w, &common.DatabaseInsertionError)
		return
	}

	url, err := generateSignedURL(image.Bucket, image.BucketPath)
	if err != nil {
		common.RespondWithError(w, &common.URLGenerationError)
	}

	response := image.toLinkedImageResponse(url)
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		common.RespondWithError(w, &common.JSONEncoderError)
	}
}

// Handle the API request to delete an image
func (h *Handler) HandleDeleteImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	uuidToGet, err := uuid.Parse(vars["uuid"])
	if err != nil {
		common.RespondWithError(w, &common.InvalidUUIDError)
		return
	}

	username, errResponse := handleJWT(r)
	if errResponse != nil {
		common.RespondWithError(w, errResponse)
		return
	}

	image, err := GetImage(h.db, uuidToGet)
	if err != nil {
		common.RespondWithError(w, &common.ImageNotFoundError)
		return
	}

	if username != image.Owner {
		common.RespondWithError(w, &common.UserPermissionDeniedError)
		return
	}

	tx, err := DeleteImage(h.db, image.UUID)
	if err != nil {
		common.RespondWithError(w, &common.DBDeletionError)
		return
	}

	if image.Status == "UPLOADED" {
		err = deleteFile(image.Bucket, image.BucketPath)
		if err != nil {
			tx.Rollback()
			common.RespondWithError(w, &common.FileDeletionError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		common.RespondWithError(w, &common.DBDeletionError)
		return
	}

	response := image.toUnlinkedImageResponse()
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		common.RespondWithError(w, &common.JSONEncoderError)
	}
}

// Handle the API request to get all images owned by the user initiating the request
func (h *Handler) HandleGetImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username, errResponse := handleJWT(r)
	if errResponse != nil {
		common.RespondWithError(w, errResponse)
		return
	}

	images, err := GetImages(h.db, username)
	if err != nil {
		common.RespondWithError(w, &common.GetImagesDBError)
		return
	}

	response := imagesToUnlinkedImagesReponse(images)
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		common.RespondWithError(w, &common.JSONEncoderError)
	}
}

// Check the validity of the JWT and return the username related to the token
func handleJWT(r *http.Request) (string, *common.ErrorResponseError) {
	if r.Header["Key"] == nil {
		return "", &common.MissingTokenError
	}

	username, err := common.VerifyJWT(r.Header["Key"][0])
	if err != nil {
		return "", &common.InvalidTokenError
	}
	return username, nil
}

// Parse the multipart-form and return the file uploaded
func getImageFromForm(w http.ResponseWriter, r *http.Request) (multipart.File, *common.ErrorResponseError) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, &common.InvalidImageBodyError
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		return nil, &common.InvalidImageBodyError
	}
	return file, nil
}
