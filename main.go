package main

import "k8s-go-prototype/kubectl"
import "k8s-go-prototype/k8s"
import "context"

func main() {
	// Apply job manifest
	kubectl.ApplyYamlFile("hello.yaml")

	// Watch logs
	k8s.WatchJob(context.TODO(), "hello-world")

	//k8s.ConnectToK8s()
}
