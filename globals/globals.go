package globals

import (
	"log"
	"os"
)

func RootDir() string {
	path, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	return path
}

func ArtifactsDir() string {
	return RootDir() + "/artifacts"
}
