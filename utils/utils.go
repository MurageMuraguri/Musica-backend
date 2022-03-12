package utils

import (
	"log"
	"os"

	"github.com/kardianos/osext"

	"github.com/joho/godotenv"
)

var dev bool = IsDev()

func LoadEnv() {

	// load .env file
	errx := godotenv.Load(Path(".env"))
	if errx != nil {
		log.Fatalf("Error loading .env file")
	}

}

func Env(key string) string {
	LoadEnv()
	return os.Getenv(key)
}

func Path(p string) string {
	if dev {
		folderPath := "./"
		return folderPath + p
	} else {
		folderPath, _ := osext.ExecutableFolder()
		return folderPath + "/" + p
	}

}

func PathA(p string) string {
	if dev {
		folderPath := "./"
		return folderPath + p
	} else {
		folderPath, _ := osext.ExecutableFolder()
		return folderPath + "/" + p
	}

}
func IsDev() bool {
	return os.Getenv("TRILL_MODE")=="DEV"
}
