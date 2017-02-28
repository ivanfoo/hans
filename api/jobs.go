package api

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ivanfoo/hans/worker/job"

	"github.com/gin-gonic/gin"
	"github.com/kr/beanstalk"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func QueueJob(c *gin.Context) {
	var j job.Job
	j.Id = bson.NewObjectId()
	j.Status = "queued"

	err := c.BindJSON(&j)

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "failed",
			"message": "incorrect request",
		})
		return
	}

	b, _ := json.Marshal(&j)

	queueConn := c.MustGet("queueConn").(*beanstalk.Conn)
	tube := &beanstalk.Tube{queueConn, "jobs"}
	j.QueueId, err = tube.Put(b, 1, 0, 10*time.Minute)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "failed",
			"message": "unable to queue job",
		})
	}

	db := c.MustGet("dbConn").(*mgo.Session)
	db.DB("hans").C("jobs").Insert(j)

	c.JSON(202, j)
}

func GetJob(c *gin.Context) {
	db := c.MustGet("dbConn").(*mgo.Session)
	jobId := bson.ObjectIdHex(c.Param("id"))
	job := job.Job{}
	err := db.DB("hans").C("jobs").FindId(jobId).One(&job)
	if err != nil {
		c.JSON(404, gin.H{
			"status":  "failed",
			"message": "job not found",
		})
	}

	c.JSON(200, job)
}

func GetLastJobs(c *gin.Context) {
	db := c.MustGet("dbConn").(*mgo.Session)
	max, _ := strconv.Atoi(c.DefaultQuery("max", "100"))
	status := c.Query("status")
	var lastJobs []job.Job
	var err error
	if status == "" {
		err = db.DB("hans").C("jobs").Find(nil).Sort("-$natural").Limit(max).All(&lastJobs)
	} else {
		err = db.DB("hans").C("jobs").Find(bson.M{"status": status}).Sort("-$natural").Limit(max).All(&lastJobs)
	}
	if err != nil {
		c.JSON(404, gin.H{
			"status":  "failed",
			"message": "job not found",
		})
	}
	c.JSON(200, lastJobs)
}
