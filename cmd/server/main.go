package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gbh007/easyjet/internal/adapters/database"
	"github.com/gbh007/easyjet/internal/entities"
)

func main() {
	err := os.Remove("main.db")
	if err != nil {
		panic(err)
	}

	r, err := database.NewRepo("sqlite", "main.db")
	if err != nil {
		panic(err)
	}

	id, err := r.SetProject(context.TODO(), entities.Project{
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
}
