package k8s

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

	labelSelector := fmt.Sprintf("job-name=%s", job.GetName())

	watchPodForSelector(ctx, labelSelector, job.GetNamespace())

}
