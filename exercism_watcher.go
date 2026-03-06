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

type codeSpec struct {
    language string
}

func runTests(dir string, language string) {
	log.Println(">>>>>>>>>>>>>> RESTARTING LOOP <<<<<<<<<<<<<<<")
	log.Printf("Detected Language:%s\n", language)

	// TODO: quick and simple for now.  We can abstract this into func later (i.e. )
	cmd := exec.Command("go", "test", "./...")

	// cmd := getTestCommand(language)

	switch {
		case language == "Python" :
			cmd = exec.Command("pytest")
		default : 
		// defaults to go
			cmd = exec.Command("go", "test", "./...")
	}
	
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s test failed in %s: %v\n%s", language, dir, err, string(output))
		return
	}
	log.Printf("%s test succeeded in %s\n%s", language, dir, string(output))
}

func getTestCommand(language string) {

	log.Println(">getTestCommand<")
	// switch {
	// cmd := exec.Command("go", "test", "./...")
	// }
}

func getUserArgs() []string {
	var userArgs []string
	userArgs = os.Args[1:]
	return userArgs
	
}
func getFileLang(filePath string) string {
	// Example: Detect the language of a file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	language := enry.GetLanguage(filePath, content)

	//fmt.Printf("File: %s\n", filePath)
	// log.Printf("File:%s\n", string(filePath))
	// log.Printf("Detected Language:%s\n", language)

	return language
}

func main() {

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
			// Ignore Python specific directories
			if d.Name() == "__pycache__" {
				return filepath.SkipDir
			}
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

		userArgs := getUserArgs()
		codeLanguage := getFileLang(userArgs[0])

		if codeLanguage == "" {
			// default to go when 0 args
			codeLanguage = "go"
		}

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("Event:", event)
				dir := filepath.Dir(event.Name)
				runTests(dir, codeLanguage)
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
