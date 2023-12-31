/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"fmt"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
)

func CreateJobFromUnstructuredObj(obj interface{}) (job *batchv1.Job) {

	job = new(batchv1.Job)

	createdUnstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		log.Fatal(err)
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(createdUnstructuredObj, &job)
	if err != nil {
		log.Fatal(err)
	}

	return

}

func CreateConfigMapFromUnstructuredObj(obj interface{}) (cm *v1.ConfigMap) {

	cm = new(v1.ConfigMap)

	createdUnstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		log.Fatal(err)
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(createdUnstructuredObj, &cm)
	if err != nil {
		log.Fatal(err)
	}

	return

}

func verifyJobCompletionStatus(prStatus string) (jobStatus string) {

	jobStatusMessage, _ := sthingsBase.GetRegexSubMatch(prStatus, `Complete\s(\w+)`)

	if jobStatusMessage != "True" {
		jobStatus = "running"
	} else {
		jobStatus = "finished"
	}

	return

}

func verifyInformerStatus(kind, function string, obj interface{}) {

	switch kind {

	case "pipelineruns":
		pipelineRu := CreateJobFromUnstructuredObj(obj)
		log.Println("pipelineRun " + function + ": " + pipelineRu.Name)
		pipelineRunStatusMessage := verifyJobCompletionStatus(fmt.Sprintln(pipelineRu.Status))
		fmt.Println(pipelineRunStatusMessage)
		// produceStatus("job-"+pipelineRu.Name, pipelineRunStatusMessage)

	case "jobs":
		job := CreateJobFromUnstructuredObj(obj)
		log.Println("job " + function + ": " + job.Name)
		jobStatusMessage := verifyJobCompletionStatus(fmt.Sprintln(job.Status))
		produceStatus("job-"+job.Name, jobStatusMessage)

	case "configmaps":
		cmStatus := "notExisting"
		cm := CreateConfigMapFromUnstructuredObj(obj)
		log.Println("configMap " + function + ": " + cm.Name)

		if function != "deleted" {
			cmStatus = "created"
		}

		produceStatus("cm-"+cm.Name, cmStatus)
	}

}
