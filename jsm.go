package jsm

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

const version = "1.0.0"

type Jsm struct {
	AppName  string
	Debug    bool
	Version  string
	RootPath string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	config   config
}

type config struct {
	port     string
	renderer string
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

	//create loggers
	infoLog, errorLog := j.startLoggers()
	j.InfoLog = infoLog
	j.ErrorLog = errorLog

	j.Version = version
	j.RootPath = rootPath
	j.Debug, err = strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		return err
	}

	j.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
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

func (j *Jsm) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}
