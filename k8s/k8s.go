package k8s

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

// TODO pass this to context and define in main?
var clientSet = Getk8sClient()

func Getk8sClient() *kubernetes.Clientset {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/root"
	}
	configPath := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Panic("failed to create K8s config")
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic("Failed to create K8s clientset")
	}

	return clientSet
}
