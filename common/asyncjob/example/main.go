package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"social_todo_list/common/asyncjob"
)

func main() {
	job1 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 2)
		fmt.Println("I am job 1")

		return nil
	}, asyncjob.WithName("Job 1"))

	job2 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 3)
		fmt.Println("I am job 2")

		return nil
	}, asyncjob.WithName("Job 2"))

	if err := asyncjob.NewGroup(true, job1, job2).Run(context.Background()); err != nil {
		log.Println(err)
	}

	//if err := job1.Execute(context.Background()); err != nil {
	//	log.Println(err)
	//
	//	for {
	//		err := job1.Retry(context.Background())
	//
	//		if err == nil || job1.State() == asyncjob.StateRetryFailed {
	//			break
	//		}
	//	}
	//}
}
