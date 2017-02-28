package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kr/beanstalk"
)

func GetQueuedJob(c *gin.Context) {
	queueConn := c.MustGet("queueConn").(*beanstalk.Conn)
	jobId, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	statsJob, err := queueConn.StatsJob(jobId)

	if err != nil {
		if err.Error() == "stats-job: not found" {
			c.JSON(200, gin.H{
				"id":    jobId,
				"state": "not found",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"id":    statsJob["id"],
		"state": statsJob["state"],
		"age":   statsJob["age"],
	})
}
