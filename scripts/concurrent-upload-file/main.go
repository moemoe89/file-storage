package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func main() {
	var images []string

	err := filepath.Walk("scripts/concurrent-upload-file/images/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		images = append(images, path)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Non-concurrent test
	//for _, image := range images {
	//	cmd, err := exec.Command("bash", "-c", "./fs-store upload-file "+image).Output()
	//	if err != nil {
	//		log.Println(err.Error())
	//	}
	//
	//	fmt.Println(string(cmd))
	//}

	errChan := make(chan error, len(images))

	var wg sync.WaitGroup

	for _, image := range images {
		wg.Add(1)

		go func(image string) {
			defer wg.Done()

			_, err := exec.Command("bash", "-c", "./fs-store upload-file "+image).Output()
			if err != nil {
				errChan <- fmt.Errorf("failed to upload %s: %w", image, err)
			}

		}(image)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		log.Fatal(err)
	default:
		close(errChan)
		log.Printf("### Uploading finished!! ###\n")
	}
}
