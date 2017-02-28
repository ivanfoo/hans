package worker

import (
	"time"

	"github.com/ivanfoo/hans/worker/job"

	"github.com/kr/beanstalk"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Worker struct {
	status string
	queue  *beanstalk.Conn
	conn   *mgo.Session
	jobs   chan *job.Job
}

func New(queue *beanstalk.Conn, conn *mgo.Session) *Worker {
	return &Worker{
		status: "up",
		queue:  queue,
		conn:   conn,
		jobs:   make(chan *job.Job),
	}
}

func (w *Worker) Run() {
	go w.consumeQueue()
	go w.runJobs()
}

func (w *Worker) consumeQueue() {
	tube := beanstalk.NewTubeSet(w.queue, "jobs")
	for {
		jn, body, err := tube.Reserve(5 * time.Second)
		if err != nil {
			continue
		}

		job := job.New(jn, body)
		w.jobs <- job
	}
}

func (w *Worker) runJobs() {
	for {
		job := <-w.jobs
		w.updateJobMongoStatus(job.Id, "running")
		err := job.Run()
		if err != nil {
			w.updateJobMongoStatus(job.Id, "failed")
		} else {
			w.updateJobMongoStatus(job.Id, "done")
		}
		w.queue.Delete(job.QueueId)
	}
}

func (w *Worker) updateJobMongoStatus(jobId bson.ObjectId, status string) {
	db := w.conn.DB("hans").C("jobs")
	db.UpdateId(jobId, bson.M{"$set": bson.M{"status": status}})
}
