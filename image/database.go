package image

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// Return a connection Pool that can query the DB parameterized by environment variables
func NewConnectionPool() (*sql.DB, error) {
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USERNAME")
	dbIP, ok := os.LookupEnv("DB_IP")
	if !ok {
		dbIP = "127.0.0.1"
	}

	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true", dbUser, dbPassword, dbIP, dbName)

	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, err
	}

	err = createTableIfNotExist(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Create the images table if it doesn't exist
func createTableIfNotExist(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS images (UUID binary(16) not null primary key, " +
		"name varchar(64) not null, owner varchar(32) not null, extension varchar(12) not null, height int null, " +
		"length int null, bucket varchar(64) not null, bucketPath varchar(128) not null, status varchar(32) null)")
	if err != nil {
		return err
	}
	return nil
}

// Create a DB record for the specified image
func CreateImage(db *sql.DB, image Image) error {
	uuidToCreate, err := image.UUID.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO images (UUID, name, owner, extension, height, length, bucket, bucketPath, "+
		"status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", uuidToCreate, image.Name, image.Owner, image.Extension, image.Height,
		image.Length, image.Bucket, image.BucketPath, image.Status)
	if err != nil {
		return err
	}
	return nil
}

// Update the image record with the same uuid as the one that is passed as parameter
func UpdateImage(db *sql.DB, image Image) error {
	uuidToUpdate, err := image.UUID.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE images SET name = ?, owner = ?, extension = ?, height = ?, length = ?, bucket = ?, "+
		"bucketPath = ?, status = ? WHERE uuid = ?", image.Name, image.Owner, image.Extension, image.Height, image.Length,
		image.Bucket, image.BucketPath, image.Status, uuidToUpdate)

	if err != nil {
		return err
	}
	return nil
}

// Delete an image record with the provided uuid
func DeleteImage(db *sql.DB, id uuid.UUID) (*sql.Tx, error) {
	uuidToDelete, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec("DELETE FROM `images` WHERE `uuid` = ?", uuidToDelete)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// Return the record associated to the image uuid
func GetImage(db *sql.DB, id uuid.UUID) (*Image, error) {
	uuidToGet, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow("SELECT * FROM images WHERE UUID = ?", uuidToGet)
	image := &Image{}
	var uuidToParse []byte

	err = row.Scan(&uuidToParse, &image.Name, &image.Owner, &image.Extension, &image.Height, &image.Length, &image.Bucket,
		&image.BucketPath, &image.Status)
	if err != nil {
		return nil, err
	}

	err = image.UUID.UnmarshalBinary(uuidToParse)
	if err != nil {
		return nil, err
	}
	return image, nil
}

// Return the record(s) of the image owned by the user passed as parameter
func GetImages(db *sql.DB, username string) ([]Image, error) {
	rows, err := db.Query("SELECT * FROM images WHERE owner = ? LIMIT 500", username)
	images := make([]Image, 0)

	for rows.Next() {
		var image Image
		var uuidToParse []byte
		err = rows.Scan(&uuidToParse, &image.Name, &image.Owner, &image.Extension, &image.Height, &image.Length,
			&image.Bucket, &image.BucketPath, &image.Status)
		if err != nil {
			return nil, err
		}
		err = image.UUID.UnmarshalBinary(uuidToParse)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}
