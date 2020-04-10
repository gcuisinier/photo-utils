package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {

	dryRun := flag.Bool("dryRun", true, "Execute a dryRun, do not remove file")
	flag.Parse()

	logFile := configureLog()
	defer logFile.Close()

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

	processFile(fileToProcess, *dryRun)

}

func configureLog() *os.File {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	return file

}

func processFile(files []string, dryRun bool) {

	if dryRun {
		fmt.Println("*****************************")
		fmt.Println("   DRY   RUN  - NOT REMOVING")
		fmt.Println("*****************************")
	}

	var totalSize int64 = 0

	for _, file := range files {

		fi, err := os.Stat(file)
		if err != nil {
			return
		}

		log.Println("Remove " + file + " " + ByteCountDecimal(fi.Size()))
		totalSize += fi.Size()
		if !dryRun {
			var err = os.Remove(file)
			if isError(err) {
				return
			}
		}

	}

	fmt.Println("+--------------------------------------------------+")
	fmt.Println("+ Total size gained : \t\t" + ByteCountDecimal(totalSize))
	fmt.Println("+ Number of file deleted : \t", (len(files)))
	fmt.Println("+                                                  +")
	fmt.Println("+--------------------------------------------------+")

	if dryRun {
		fmt.Println("*****************************")
		fmt.Println("   DRY   RUN  - NOT REMOVING")
		fmt.Println("*****************************")
	}

}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
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
