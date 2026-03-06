package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"github.com/fsnotify/fsnotify"
	"github.com/go-enry/go-enry/v2"
	"io/ioutil"
)

func runTests(dir string) {
	log.Println(">>>>>>>>>>>>>> RESTARTING LOOP <<<<<<<<<<<<<<<")

	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("go test failed in %s: %v\n%s", dir, err, string(output))
		return
	}
	log.Printf("go test succeeded in %s\n%s", dir, string(output))
}
func getUserArgs() []string {
	userArgs := os.Args[1:]
	// log.Printf("first userArg:%s\n", userArgs[0])
	return userArgs
}
func getFileLang(filename string) string {
	// Example: Detect the language of a file
	filePath := "example.go"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	language := enry.GetLanguage(filePath, content)

	//fmt.Printf("File: %s\n", filePath)
	log.Printf("File:%s\n", string(filePath))
	log.Printf("Detected Language:%s\n", language)
	// log.Printf("Reliable:%s\n", string(reliable))

	return language
}

func main() {
	// userArgs := os.Args[1:]
	//  if len(userArgs) > 0 {
	// 	log.Printf("User arguments: \n%s", string(userArgs[0]))
    // }
	getUserArgs()

	// log.Fatal(string("Force stop"))

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Recursively add directories
	err = filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Ignore directories starting with '.'
			if strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
				return filepath.SkipDir
			}
			log.Println("Watching:", path)
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Process events
	go func() {

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("Event:", event)
				dir := filepath.Dir(event.Name)
				runTests(dir)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	<-make(chan struct{}) // Block forever
}
