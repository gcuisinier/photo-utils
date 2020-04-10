package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {

	configureLog()

	var fileToProcess []string

	rawExtenstion := []string{".RAW", ".NEF", ".ARW"}
	sort.Strings(rawExtenstion)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Scanning directory : " + dir)

	err = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if isPresentInList(rawExtenstion, filepath.Ext(path)) {
				log.Println(path)

				jpgFile := strings.TrimSuffix(path, filepath.Ext(path)) + ".jpg"
				if fileExists(jpgFile) {
					log.Println("Both file " + path + " and " + jpgFile + " exists")
					fileToProcess = append(fileToProcess, path)

				}

			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	processFile(fileToProcess)

}

func configureLog() {
	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	//log.SetOutput(f)

}

func processFile(files []string) {

	var totalSize int64 = 0

	for _, file := range files {

		fi, err := os.Stat(file)
		if err != nil {
			return
		}

		log.Println("Remove " + file + " " + ByteCountDecimal(fi.Size()))

		totalSize += fi.Size()
	}

	fmt.Println("Total size gained :" + ByteCountDecimal(totalSize))

}

func ByteCountDecimal(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func isPresentInList(extensions []string, currentExtension string) bool {

	sort.Strings(extensions)
	i := sort.SearchStrings(extensions, currentExtension)
	if i < len(extensions) && extensions[i] == currentExtension {
		return true
	}
	return false

}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
