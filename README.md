# Shopify Backend Challenge - Image Microservice
This is a simple GO microservice that allows the upload and consumption of images through a REST API. The following actions are supported :
 - Create a DB Record for an image
 - Get details about a specified image including a temporary download link
 - Upload an image to cloud storage
 - Delete an image from the Cloud Storage including its database records
 - Get details of all images owned by the authenticated user

For a detailed documentation on how to query the API, please visit the [SwaggerHup API Page](https://app.swaggerhub.com/apis-docs/wtrep/shopify-images-repo/1.0.0).

## Details about the microservice
### File structure
The main logic of the microservice can be found in the image package. The common package is shared between this microservice and the [auth microservice](https://github.com/wtrep/shopify-backend-challenge-auth).

### Authentification
This microservice rely on the [Shopify Backend Challenge Image Microservice](https://github.com/wtrep/shopify-backend-challenge-auth) to authenticate users and to
provide JWT to users. Since the signing key is shared as Kubernetes Secret between the two microservices, the Image Microservice can authenticate users without
relying on an active sessions database.

### Database
The microservice needs to have access to a MySQL database. You can find the Terraform code for a GCP Cloud SQL instance in the [main repository](https://github.com/wtrep/shopify-backend-challenge/tree/master/terraform/cloud_sql).

### Cloud Storage
The images are hosted on a GCP Cloud Storage Bucket. The microservice needs to have access to a GCP service account that allows write access to the repository and the permission to generate temporary download links. An example can be found in the [main repository.](https://github.com/wtrep/shopify-backend-challenge/tree/master/terraform/bucket)

### Docker Image and Kubernetes
The microservice is packaged into a Dockerimage to allow deployment into a Kubernetes Cluster. You can also download the built image directly from [Docker Hub](https://hub.docker.com/r/wtrep/shopify-backend-challenge-image)

## Environment variables
The following environment variables need to be set for the microservice to work : \\
| Environment variable           | Description                                                                                                                            |
| -------------------------------|:--------------------------------------------------------------------------------------------------------------------------------------:|
| DB_IP                          | IP Address of the database                                                                                                             |
| DB_USERNAME                    | Username to access the MySQL DB                                                                                                        |
| DB_PASSWORD                    | Password to access the MySQL DB                                                                                                        |
| DB_NAME                        | Name of the MySQL database                                                                                                             |
| JWT_KEY                        | Private key to verify JWT Tokens. Must be the same as the [auth microservice](https://github.com/wtrep/shopify-backend-challenge-auth) |
| BUCKET                         | Name of the GCP Bucket where to upload the images                                                                                      |
| GOOGLE_APPLICATION_CREDENTIALS | Path to the Service Account .json file to allow Bucket write access                                                                    |

## Build and run
To build the microservice : 
```
go build main.go
```

To run the microservice, you must configure the required environment variables first. Then run :
```
go run main.go
```

## Build the Docker image
The provided Dockerfile allow the microservice to be packaged into a container. To build the Docker image:
```
docker build .
```

## List of possible improvements 
 * Automated tests
 * Use OAuth2 
 * Granular permission and public image access
 * Ability to share images with other users
 * CICD Pipeline that builds and upload to Docker Hub a new Docker image at each merge to the master branch