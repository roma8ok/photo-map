package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/roma8ok/photo-map/additional"
)

var (
	paths                   = make([]string, 0)
	pathsWithoutCoordinates = make([]string, 0)
	pathsWithinCoordinates  = make([]additional.Image, 0)
)

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
		image, scanImageError := additional.ScanImage(path)
		if scanImageError != nil {
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

	mapRemoveErr := os.Remove(mapPath)
	if mapRemoveErr != nil {
		fmt.Println(mapRemoveErr)
	}
	htmlFile, htmlOpenFileErr := os.OpenFile(
		mapPath, os.O_WRONLY|os.O_CREATE, 0755)
	if htmlOpenFileErr != nil {
		fmt.Println(htmlOpenFileErr)
		return
	}
	defer htmlFile.Close()

	markers := make([]string, 0, len(pathsWithinCoordinates))
	for _, image := range pathsWithinCoordinates {
		markerString := fmt.Sprintf(
			`    L.circleMarker([%f, %f], { color: "#343E40", weight: 1, fillColor: "%s", fillOpacity: 0.5 }).addTo(mymap).bindPopup("<a href='%s' target='_blank'>%s</a>");`,
			image.Lat, image.Long, additional.ConvertYearToColor(image.DateTime.Year()),
			image.Path, image.DateTime.Format("2006-01-02 15:04:05"))
		markers = append(markers, markerString)
	}

	_, htmlFileWriteErr := htmlFile.WriteString(
		additional.StartHTML + strings.Join(markers, "\n") + additional.EndHTML)
	if htmlFileWriteErr != nil {
		fmt.Println(htmlFileWriteErr)
		return
	}

	fmt.Printf("show photo map here: %s\n", mapPath)
}
