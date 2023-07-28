package k8s

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var clientSet = Getk8sClient()

func GetJob(ctx context.Context, jobName string) (*batchv1.Job, error) {

	// TODO pass namespace
	jobsClient := clientSet.BatchV1().Jobs(metav1.NamespaceDefault)

	return jobsClient.Get(ctx, jobName, metav1.GetOptions{})

}

func WatchJob(ctx context.Context, name string) {

	job, err := GetJob(ctx, name)

	if err != nil {
		log.Fatal(err)
	}

	pod, err := GetPodForJob(ctx, job)

	if err != nil {
		log.Fatal(err)
	}

	log.Info(pod)

}

func GetPodForJob(ctx context.Context, job *batchv1.Job) (*corev1.Pod, error) {

	name := job.ObjectMeta.Name
	nameSpace := job.ObjectMeta.Namespace

	podList, err := clientSet.CoreV1().Pods(nameSpace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", name),
	})

	if err != nil {
		return nil, err
	}

	if len(podList.Items) == 0 {
		return nil, fmt.Errorf("No pod has been found for job %s yet...", name)
	}

	if len(podList.Items) != 1 {
		return nil, fmt.Errorf("Expected job %s to have one pod, but instead got %v.", name, len(podList.Items))
	}
	return &podList.Items[0], nil
}
