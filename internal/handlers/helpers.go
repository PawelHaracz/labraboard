package handlers

import (
	"fmt"
	"os"
)

func createLabraboardBackendFile(path string, endpoint string, projectId string) error {
	content := `terraform {
  backend "http" {
    address = "%[1]s/api/v1/state/terraform/%[2]s"
    lock_address = "%[1]s/api/v1/state/terraform/%[2]s/lock"
    unlock_address = "%[1]s/api/v1/state/terraform/%[2]s/lock"
  }
}`
	file, err := os.Create(fmt.Sprintf("%s/backend_override.tf", path))
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, err = fmt.Fprintf(file, content, endpoint, projectId)
	if err != nil {
		return err
	}

	return nil
}

func createBackendFile(path string, statePath string) error {
	content := `terraform {
  backend "local" {
    path = "%s"
  }
}`

	file, err := os.Create(fmt.Sprintf("%s/backend_override.tf", path))
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, err = fmt.Fprintf(file, content, statePath)
	if err != nil {
		return err
	}

	return nil
}
