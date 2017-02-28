package job

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/ivanfoo/hans/utils"

	"gopkg.in/mgo.v2/bson"
)

var ansibleCommand = "ansible-playbook"

type Job struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"id" binding:"required"`
	QueueId  uint64        `bson:"queueId" json:"queueId"`
	Playbook string        `json:"playbook" binding:"required"`
	Hosts    []string      `json:"hosts" binding:"required"`
	Aws      *AwsParams    `bson:"aws,omitempty" json:"aws,omitempty"`
	Status   string        `json:"status" binding:"required"`
}

type AwsParams struct {
	Session string `bson:"session,omitempty" json:"session,omitempty"`
	RoleArn string `bson:"roleArn,omitempty" json:"roleArn,omitempty"`
}

func New(queueId uint64, body []byte) *Job {
	var j Job
	json.Unmarshal(body, &j)
	j.QueueId = queueId

	return &j
}

func (j *Job) Run() error {
	if j.Aws != nil {
		os.Setenv("AWS_PROFILE", j.Aws.Session)

		clearCache := exec.Command("rm", "-rf", "~/.ansible/tmp/*")
		err := clearCache.Run()
		if err != nil {
			return err
		}

		err = utils.AssumeAwsRole(j.Aws.RoleArn, j.Aws.Session)
		if err != nil {
			return err
		}
	}

	return j.runPlaybook()
}

func (j *Job) runPlaybook() error {
	stringifiedHosts := strings.Join(j.Hosts[:], ",")
	os.Chdir(os.Getenv("HANS_ANSIBLE_STUFF"))

	cmd := exec.Command(ansibleCommand, j.Playbook, "-l", stringifiedHosts)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
