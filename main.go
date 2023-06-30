package main

import "k8s-go-prototype/kubectl"

func main() {
    // Apply job manifest
	kubectl.ApplyYamlFile("hello.yaml")
}