package k8s

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	toolsWatch "k8s.io/client-go/tools/watch"
)

func watchPodForSelector(ctx context.Context, labelSelector string, namespace string) (*corev1.Pod, error) {

	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		// TODO pass this as a parameter
		timeOut := int64(60)
		return clientSet.CoreV1().Pods(namespace).Watch(context.Background(), metav1.ListOptions{
			TimeoutSeconds: &timeOut,
			TypeMeta:       metav1.TypeMeta{},
			LabelSelector:  labelSelector,
			FieldSelector:  "",
		})
	}

	watcher, _ := toolsWatch.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFunc})

	logsWatched := false
	for event := range watcher.ResultChan() {

		pod := event.Object.(*corev1.Pod)

		log.WithFields(log.Fields{
			"podName":  pod.GetName(),
			"event":    event.Type,
			"podPhase": pod.Status.Phase,
		}).Info("Event for pod watched")

		if pod.Status.Phase == corev1.PodSucceeded {
			log.Info("Pod Succeeded")
			return pod, nil
		}

		if pod.Status.Phase == corev1.PodFailed {
			log.Info("Pod Failed")
			return pod, nil
		}

		if pod.Status.Phase == corev1.PodUnknown {
			log.Info("Pod status Unknown")
			return pod, nil
		}

		// Running - show logs
		if pod.Status.Phase == corev1.PodRunning && logsWatched == false {
			log.Info("Pod is running, show logs below:")
			GetPodLogs(pod, true)
			logsWatched = true
		}

	}

	return nil, fmt.Errorf("Running pod has not been found for labelSelector %s", labelSelector)
}

func GetPodLogs(pod *corev1.Pod, follow bool) error {

	namespace := pod.ObjectMeta.Namespace
	podName := pod.ObjectMeta.Name

	count := int64(100)
	podLogOptions := corev1.PodLogOptions{
		Follow:    follow,
		TailLines: &count,
	}

	podLogRequest := clientSet.CoreV1().
		Pods(namespace).
		GetLogs(podName, &podLogOptions)
	stream, err := podLogRequest.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		buf := make([]byte, 2000)
		numBytes, err := stream.Read(buf)

		if err == io.EOF {
			break
		}
		if numBytes == 0 {
			continue
		}

		if err != nil {
			return err
		}
		message := string(buf[:numBytes])
		fmt.Print(message)
	}
	return nil
}
