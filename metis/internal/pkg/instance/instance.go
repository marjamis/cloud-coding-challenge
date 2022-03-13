package instance

import (
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// InstanceID is used to store the instance id for use throughout metis
	InstanceID string
)

func init() {
	InstanceID = GetInstanceID()
}

// GetInstanceID curls the metadata service to get the instances id or if not found report this.
func GetInstanceID() (instanceID string) {
	instanceID = "NOT_FOUND"

	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "GetInstanceId#1",
		}).Error(err)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"Function": "GetInstanceId#2",
		}).Error(err)
	}
	instanceID = string(data)

	return
}
