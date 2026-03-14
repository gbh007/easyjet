package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gbh007/easyjet/internal/adapter/repository/gorm"
	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/golang-cz/devslog"
	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.Default()
	lt := "dev"
	llv := slog.LevelDebug

	switch lt {
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: llv,
		}))
	case "dev":
		logger = slog.New(devslog.NewHandler(os.Stdout, &devslog.Options{
			HandlerOptions: &slog.HandlerOptions{
				Level: llv,
			},
		}))
	case "tint":
		logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level: llv,
		}))
	}

	err := os.Remove("main.db")
	if err != nil {
		panic(err)
	}

	r, err := gorm.NewRepo(logger, "sqlite", "main.db")
	if err != nil {
		panic(err)
	}

	_, err = r.SetProject(context.TODO(), entity.Project{
		Stages: []entity.ProjectStage{
			{
				Script: "1",
			},
			{
				Script: "2",
			},
			{
				Script: "3",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	id, err := r.SetProject(context.TODO(), entity.Project{
		Dir:    "hello",
		GitURL: "world",
		Name:   "123",
		Stages: []entity.ProjectStage{
			{
				Script: "1",
			},
			{
				Script: "2",
			},
			{
				Script: "3",
			},
		},
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
		Stages: []entity.ProjectStageRun{
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
