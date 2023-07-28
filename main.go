package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"k8s-go-prototype/k8s"
	"k8s-go-prototype/kubectl"
)

func main() {

	jobManifest := "hello.yaml"

	jobName := "hello-world"

	// Check if the job exists already
	_, err := k8s.GetJob(context.TODO(), jobName)
	if err == nil {
		log.WithFields(log.Fields{
			"jobName": jobName,
		}).Fatal("That job exists already")
	}

	// Apply job manifest
	kubectl.ApplyYamlFile(jobManifest)

	// Watch logs
	k8s.WatchJob(context.TODO(), jobName)
}
