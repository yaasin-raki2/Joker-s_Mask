package jsm

import (
	"fmt"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

type Jsm struct {
	AppName string
	Debug   bool
	Version string
}

func (j *Jsm) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath: rootPath,
		folderNames: []string{"handlers", "migrations", "views",
			"data", "public", "tmp", "middleware"},
	}

	err := j.Init(pathConfig)
	if err != nil {
		return err
	}

	err = j.checkDotenv(rootPath)
	if err != nil {
		return err
	}

	//read .env
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	return nil
}

func (j *Jsm) Init(p initPaths) error {
	//root := p.rootPath

	for _, path := range p.folderNames {
		//create folder if it doesn't exist
		err := j.CreateDirIfNotExist(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *Jsm) checkDotenv(path string) error {
	err := j.CreateFileIfNotExist(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}
