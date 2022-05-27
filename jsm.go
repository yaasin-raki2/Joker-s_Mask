package jsm

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/yaasin-raki2/Joker-s_Mask/render"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const version = "1.0.0"

type Jsm struct {
	AppName  string
	Debug    bool
	Version  string
	RootPath string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Routes   *chi.Mux
	Render   *render.Render
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
	j.Routes = j.routes().(*chi.Mux)
	j.Debug, err = strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		return err
	}

	j.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	j.Render = j.createRenderer()

	return nil
}

// Init creates necessary folders for your JSM application
func (j *Jsm) Init(p initPaths) error {
	root := p.rootPath

	for _, path := range p.folderNames {
		//create folder if it doesn't exist
		err := j.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

// ListenAndServe starts the webserver
func (j *Jsm) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", j.config.port),
		ErrorLog:     j.ErrorLog,
		Handler:      j.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}
	j.InfoLog.Printf("Listening on port %s\n", j.config.port)
	err := srv.ListenAndServe()
	if err != nil {
		j.ErrorLog.Fatal(err)
	}
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

func (j *Jsm) createRenderer() *render.Render {
	myRenderer := &render.Render{
		Renderer: j.config.renderer,
		RootPath: j.RootPath,
		Port:     j.config.port,
	}
	return myRenderer
}
