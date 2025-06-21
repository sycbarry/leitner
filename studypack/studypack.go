package studypack

import (
	"fmt"
	"os"
	"path/filepath"
)

type StudyPackInterface interface {
	InitializePackage()
	Name() string
}

type StudyPackage struct {
	PackageName string
}

func (sp *StudyPackage) InitializePackage() {
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting user home directory:", err)
		os.Exit(1)
	}

	rootPath := filepath.Join(homeDirPath, ".leitner")

	// Create root folder if it doesn't exist
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		fmt.Println("Error: .leitner directory not found. Please run 'leitner init' to initialize.")
		os.Exit(1)
	}

	// Create package folder
	packagePath := filepath.Join(rootPath, sp.PackageName)
	if err := os.Mkdir(packagePath, 0755); err != nil && !os.IsExist(err) {
		fmt.Println("error creating package directory:", err)
		os.Exit(1)
	}
}

func (sp *StudyPackage) Name() string {
	return sp.PackageName
}
