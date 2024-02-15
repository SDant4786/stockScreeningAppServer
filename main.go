package main

import (
	a "./algorithm"
	"./api"
	"log"
	"os"
)
func main() {
	log.Println("Application started")
	go api.StartServer()

	for _, algorithm := range a.GetAllAlgorithms() {
		a.AddToAlgorithmMap(algorithm)
		if algorithm.RunOnStart == true {
			a.StartAlgorithm(algorithm.UserName, algorithm.UniqueID)
		}
	}

	for {
		switch {
		case <- api.ProgramShutDown:
			log.Println("Shutting down program")
			os.Exit(1)
		}
	}

}