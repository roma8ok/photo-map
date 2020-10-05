package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/roma8ok/goexif/exif"
)

type Image struct {
	Path      string
	DateTime  time.Time
	Latitude  float64
	Longitude float64
}

type HTMLMarker struct {
	Image
	FillColor         string
	FormattedDateTime string
}

func (image *Image) Scan(file *os.File) error {
	exifData, err := exif.Decode(file)
	if err != nil {
		return err
	}
	image.Path = file.Name()

	lat, long, err := exifData.LatLong()
	if err != nil {
		return err
	}
	image.Latitude = lat
	image.Longitude = long

	dateTime, err := exifData.DateTime()
	if err != nil {
		return err
	}
	image.DateTime = dateTime

	return nil
}

func ScanPathForImages(path string) (image Image, err error) {
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

// Convert year to pantone color of the year.
func ConvertYearToColor(year int) (color string) {
	switch year {
	// TODO: change color in 2021 year.
	case 2021:
		color = "#788995"
	case 2020:
		color = "#34558B"
	case 2019:
		color = "#FA7268"
	case 2018:
		color = "#5f4b8b"
	case 2017:
		color = "#91b54d"
	case 2016:
		color = "#93a9d1"
	case 2015:
		color = "#964f4c"
	case 2014:
		color = "#b163a3"
	case 2013:
		color = "#009874"
	case 2012:
		color = "#e2492f"
	case 2011:
		color = "#d94f70"
	case 2010:
		color = "#45b8ac"
	case 2009:
		color = "#efc050"
	case 2008:
		color = "#5b5ea6"
	case 2007:
		color = "#9b2335"
	case 2006:
		color = "#decdbe"
	case 2005:
		color = "#55b4b0"
	case 2004:
		color = "#e15d44"
	case 2003:
		color = "#7fcdcd"
	case 2002:
		color = "#bc243c"
	case 2001:
		color = "#c34e7c"
	case 2000:
		color = "#98b4d4"
	default:
		color = "#FCEA76"
	}

	return
}

func Walk(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return
	}

	isDir := info.IsDir()

	if isDir {
		fmt.Println("start check directory:", file.Name())

		dirNames, err := file.Readdirnames(0)
		if err != nil {
			return err
		}

		for _, dirName := range dirNames {
			err := Walk(filepath.Join(file.Name(), dirName))
			if err != nil {
				return err
			}
		}
	} else {
		extension := strings.ToLower(filepath.Ext(file.Name()))
		if extension == ".jpg" || extension == ".jpeg" || extension == ".png" {
			paths = append(paths, file.Name())
		}
	}

	return
}

var (
	paths                   = make([]string, 0)
	pathsWithoutCoordinates = make([]string, 0)
	pathsWithinCoordinates  = make([]Image, 0)
)

func main() {
	startTime := time.Now()

	if len(os.Args) < 2 {
		fmt.Println("enter the path to the directory with the photos")
		return
	}
	pathToPhotos := os.Args[1]

	walkError := Walk(pathToPhotos)
	if walkError != nil {
		fmt.Println(walkError)
		return
	}

	for _, path := range paths {
		image, err := ScanPathForImages(path)
		if err != nil {
			pathsWithoutCoordinates = append(pathsWithoutCoordinates, path)
		} else {
			pathsWithinCoordinates = append(pathsWithinCoordinates, image)
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Println("elapsedTime:", elapsedTime)
	fmt.Println("len(paths):", len(paths))
	fmt.Println("len(pathsWithinCoordinates):", len(pathsWithinCoordinates))
	fmt.Println("len(pathsWithoutCoordinates):", len(pathsWithoutCoordinates))

	mapPath := filepath.Join(pathToPhotos, "photo-map.html")

	err := os.Remove(mapPath)
	if err != nil {
		fmt.Println(err)
	}
	htmlFile, err := os.OpenFile(
		mapPath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer htmlFile.Close()

	htmlMarkers := make([]HTMLMarker, 0, len(pathsWithinCoordinates))
	for _, image := range pathsWithinCoordinates {
		htmlMarker := HTMLMarker{
			Image:             image,
			FillColor:         ConvertYearToColor(image.DateTime.Year()),
			FormattedDateTime: image.DateTime.Format("2006-01-02 15:04:05"),
		}
		htmlMarkers = append(htmlMarkers, htmlMarker)
	}

	t, err := template.ParseFiles("html/map.html")
	if err != nil {
		fmt.Println("can't parse template")
		return
	}
	templateData := struct {
		Markers []HTMLMarker
	}{
		Markers: htmlMarkers,
	}

	err = t.Execute(htmlFile, templateData)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("show photo map here: %s\n", mapPath)
}
