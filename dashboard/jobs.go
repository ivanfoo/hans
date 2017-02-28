package dashboard

import (
	"encoding/json"
	"net/http"

	"github.com/ivanfoo/hans/worker/job"

	"github.com/gin-gonic/gin"
)

func GetBasicStats(c *gin.Context) {
	var lastJobs []job.Job
	url := "http://localhost:8080/api/jobs"
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(&lastJobs)

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "failed",
			"message": "api is unavailable",
		})
		return
	}

	c.HTML(200, "last_jobs.tmpl", gin.H{
		"lastJobs": lastJobs,
	})
}
