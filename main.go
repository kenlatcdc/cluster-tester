package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// testDeps is a required implementation for testing.MainStart
type testDeps struct{}

func (d testDeps) MatchString(pat, str string) (bool, error) { return true, nil }
func (d testDeps) StartCPUProfile(w interface{}) error       { return nil }
func (d testDeps) StopCPUProfile()                           {}
func (d testDeps) WriteHeapProfile(w interface{}) error      { return nil }
func (d testDeps) WriteProfileTo(string, interface{}, int) error {
	return nil
}
func (d testDeps) ImportPath() string         { return "" }
func (d testDeps) StartTestLog(w interface{}) {}
func (d testDeps) StopTestLog() error         { return nil }

func main() {
	var testFiles []string
	var wg sync.WaitGroup

	// Root directory where your script will run
	rootDir := "."

	// Parse all folders in the root directory and add those with a "tests" subfolder to testFiles
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the directory contains a "tests" subfolder
		if info.IsDir() && filepath.Base(path) == "tests" {
			testFiles = append(testFiles, path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error while walking directories: %v\n", err)
		os.Exit(1)
	}

	if len(testFiles) == 0 {
		fmt.Println("No test folders found.")
		return
	}

	fmt.Printf("Found %d test folders: %v\n", len(testFiles), testFiles)

	// Run tests for each folder concurrently
	for _, testFile := range testFiles {
		wg.Add(1)
		go func(testFolder string) {
			defer wg.Done()
			fmt.Printf("Running tests for folder: %s\n", testFolder)

			for i := 0; i < 1000; i++ {
				fmt.Printf("Run %d for folder: %s\n", i+1, testFolder)

				// You could modify this to run specific tests for the test folder.
				result := testing.MainStart(testDeps{}, []testing.InternalTest{
					{"TestMyFunction", TestMyFunction}, // Replace with actual tests for each folder
				}, nil, nil).Run()

				if result != 0 {
					fmt.Printf("Test failed on run %d for folder: %s\n", i+1, testFolder)
					os.Exit(1)
				}
			}
		}(testFile)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All tests completed.")
}
