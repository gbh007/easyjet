package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gbh007/easyjet/internal/adapter/repository/gorm"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func main() {
	err := os.Remove("main.db")
	if err != nil {
		panic(err)
	}

	r, err := gorm.NewRepo("sqlite", "main.db")
	if err != nil {
		panic(err)
	}

	id, err := r.SetProject(context.TODO(), entity.Project{
		Dir:    "hello",
		GitURL: "world",
		Name:   "123",
		Stages: []string{"1", "2", "3"},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(id)

	p, err := r.Project(context.TODO(), id)
	if err != nil {
		panic(err)
	}

	fmt.Println(p)

	p.Name = "vasya"

	id, err = r.SetProject(context.TODO(), p)
	if err != nil {
		panic(err)
	}

	fmt.Println(id)

	p, err = r.Project(context.TODO(), id)
	if err != nil {
		panic(err)
	}

	fmt.Println(p)

	runID, err := r.SetProjectRun(context.TODO(), entity.ProjectRun{
		ProjectID: id,
		Success:   true,
		Stages: []entity.ProjectRunStage{
			{
				StageNumber: 2,
				Success:     true,
				Log:         "aaaa\ndasdas",
			},
			{
				StageNumber: 5,
				Success:     false,
				Log:         "aa231aa\nda35sdas",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(runID)

	runs, err := r.ProjectRuns(context.TODO(), id)
	if err != nil {
		panic(err)
	}

	fmt.Println(runs)
}
