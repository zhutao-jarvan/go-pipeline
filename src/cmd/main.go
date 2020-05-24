package main

import (
	"fmt"
	"pipeline"
)

func main()  {
	fmt.Println("Entre main.")
	pkts := pipeline.NewPkts()

	pl := pipeline.NewPipeline(2, pkts)

	pl.Run()
	fmt.Println("Exit main.")
}

