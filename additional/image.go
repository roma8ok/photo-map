package additional

import (
	"os"
	"time"

	"github.com/roma8ok/goexif/exif"
)

type Image struct {
	Path     string
	DateTime time.Time
	Lat      float64
	Long     float64
}

func (image *Image) Scan(file *os.File) error {
	exifData, exifDataError := exif.Decode(file)
	if exifDataError != nil {
		return exifDataError
	}
	image.Path = file.Name()

	lat, long, latLongError := exifData.LatLong()
	if latLongError != nil {
		return latLongError
	}
	image.Lat = lat
	image.Long = long

	dateTime, dateTimeError := exifData.DateTime()
	if dateTimeError != nil {
		return dateTimeError
	}
	image.DateTime = dateTime

	return nil
}

func ScanImage(path string) (image Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	err = image.Scan(file)
	if err != nil {
		return
	}

	return
}
