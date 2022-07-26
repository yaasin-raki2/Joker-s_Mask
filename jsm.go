package jsm

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/yaasin-raki2/Joker-s_Mask/render"
	"github.com/yaasin-raki2/Joker-s_Mask/session"
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
	JetViews *jet.Set
	Session  *scs.SessionManager
	DB       Database
	config   config
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
	database    databaseConfig
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

	//connect to database
	if os.Getenv("DATABASE_TYPE") != "" {
		db, err := j.OpenDB(os.Getenv("DATABASE_TYPE"), j.BuildDsn())
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}
		j.DB = Database{
			DataType: os.Getenv("DATABASE_TYPE"),
			Pool:     db,
		}
	}

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
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifeTime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSIST"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database: databaseConfig{
			database: os.Getenv("DATABASE_TYPE"),
			dsn:      j.BuildDsn(),
		},
	}

	// create a session
	sess := session.Session{
		CookieName:     j.config.cookie.name,
		CookieLifeTime: j.config.cookie.lifeTime,
		CookiePersist:  j.config.cookie.persist,
		SessionType:    j.config.sessionType,
		CookieDomain:   j.config.cookie.domain,
	}

	j.Session = sess.InitSession()

	j.JetViews = jet.NewSet(jet.NewOSFileSystemLoader(
		fmt.Sprintf("%s/views", rootPath)), jet.InDevelopmentMode())

	j.createRenderer()

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

	defer func() {
		err := j.DB.Pool.Close()
		if err != nil {
			os.Exit(1)
		}
	}()

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

func (j *Jsm) createRenderer() {
	j.Render = &render.Render{
		Renderer: j.config.renderer,
		RootPath: j.RootPath,
		Port:     j.config.port,
		JetViews: j.JetViews,
	}
}

func (j *Jsm) BuildDsn() string {
	var dsn string

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"))
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}
	default:

	}

	return dsn
}
