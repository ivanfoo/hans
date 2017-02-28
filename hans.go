package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ivanfoo/hans/api"
	"github.com/ivanfoo/hans/dashboard"
	"github.com/ivanfoo/hans/utils"
	"github.com/ivanfoo/hans/webhook"
	"github.com/ivanfoo/hans/worker"

	"github.com/gin-gonic/gin"
	"github.com/kr/beanstalk"
	"gopkg.in/mgo.v2"
)

var (
	listenPort      = flag.String("port", "8080", "listen port")
	beanstalkServer = flag.String("beanstalk", "127.0.0.1:11300", "beanstalkd server")
	mongoServer     = flag.String("mongo", "127.0.0.1:27017", "mongodb server")
	numWorkers      = flag.Int("workers", 3, "number of parellel ansible jobs")
)

func main() {
	flag.Parse()

	if os.Getenv("HANS_GIT_CREDENTIALS") == "" {
		log.Fatal("You must define HANS_GIT_CREDENTIALS")
	}

	if os.Getenv("HANS_GIT_REPOSITORY") == "" {
		log.Fatal("You must define HANS_GIT_REPOSITORY")
	}

	if os.Getenv("HANS_LISTEN_ADDRESS") != "" {
		*numWorkers, _ = strconv.Atoi(os.Getenv("HANS_LISTEN_ADDRESS"))
	}

	if os.Getenv("HANS_MONGO_SERVER") != "" {
		*numWorkers, _ = strconv.Atoi(os.Getenv("HANS_MONGO_SERVER"))
	}

	if os.Getenv("HANS_BEANSTALK_SERVER") != "" {
		*numWorkers, _ = strconv.Atoi(os.Getenv("HANS_BEANSTALK_SERVER"))
	}

	if os.Getenv("HANS_WORKERS") != "" {
		*numWorkers, _ = strconv.Atoi(os.Getenv("HANS_WORKERS"))
	}

	utils.SetGitCredentials(os.Getenv("HANS_GIT_REPOSITORY"), os.Getenv("HANS_GIT_CREDENTIALS"))
	utils.CloneGitRepository(os.Getenv("HANS_GIT_REPOSITORY"), "origin")

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(MongoDB(*mongoServer))
	router.Use(Beanstalkd(*beanstalkServer))

	pwd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	router.Static("/public", filepath.Join(pwd, "/dashboard/public"))
	router.LoadHTMLGlob(filepath.Join(pwd, "/dashboard/templates/*"))

	dbConn, err := mgo.Dial(*mongoServer)

	if err != nil {
		log.Fatal(err)
	}

	queueConn, err := beanstalk.Dial("tcp", *beanstalkServer)

	if err != nil {
		log.Fatal(err)
	}

	var workers = make([]*worker.Worker, *numWorkers)
	for key, _ := range workers {
		workers[key] = worker.New(queueConn, dbConn)
	}

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/dashboard/jobs")
	})

	apiEndpoint := router.Group("/api")
	{
		apiEndpoint.GET("/queue/:id", api.GetQueuedJob)
		apiEndpoint.GET("/jobs", api.GetLastJobs)
		apiEndpoint.GET("/jobs/:id", api.GetJob)
		apiEndpoint.POST("/jobs", api.QueueJob)
	}

	statsEndpoint := router.Group("/dashboard")
	{
		statsEndpoint.GET("/jobs", dashboard.GetBasicStats)
	}

	router.POST("/webhook", webhook.RefreshGitRepo)

	for key, _ := range workers {
		go workers[key].Run()
	}

	router.Run(":" + *listenPort)
}
