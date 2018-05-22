package main

import (
	"log"
	"math/rand"
	"path/filepath"
	"sync"
	"time"

	"os"
	"os/exec"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	probabilityOfFork  = kingpin.Flag("fork-probability", "The probability (0 - 1) that this app will spawn child processes").Short('f').Default("0.5").Float()
	probabilityOfSpike = kingpin.Flag("spike-probability", "The probability (0 - 1) that this app will open a LOT of file handles").Short('s').Default("0.5").Float()
	workerProcessName  = kingpin.Flag("worker-process", "The name of the worker process to be spawned").Short('w').String()
	logFile            = kingpin.Flag("log-file", "The file to write log messages to").Short('l').OpenFile(os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
)

func main() {
	kingpin.Parse()

	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.LUTC)

	logOutput := os.Stderr

	if *logFile != nil {
		logOutput = *logFile
	}

	workerProcess := filepath.Dir(os.Args[0]) + "/worker"
	if *workerProcessName != "" {
		workerProcess = *workerProcessName
	}

	log.SetOutput(logOutput)

	log.Printf("fork probability = %v%%\n", 100*(*probabilityOfFork))
	log.Printf("spike probability = %v%%\n", 100*(*probabilityOfSpike))
	log.Printf("worker process = '%v'\n", workerProcess)
	log.Println("Logging ready")

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	forkValue := rng.Float64()
	spikeValue := rng.Float64()

	doFork := (forkValue <= *probabilityOfFork)
	doSpike := (spikeValue <= *probabilityOfSpike)

	if doFork {
		log.Println("Spawning worker processes ...")
		args := make([]string, 0, 1)

		if doSpike {
			args = append(args, "--do-spike")
		}

		var waitGroup sync.WaitGroup
		for i := 0; i < 10; i++ {
			cmd := exec.Command(workerProcess, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = logOutput

			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				if err := cmd.Start(); err != nil {
					log.Fatal(err)
				}
				cmd.Wait()
			}()
		}

		waitGroup.Wait()
	} else {
		for {
			log.Println("I didn't have to spawn any workers, so I'm going to harmlessly do some work every few seconds...")
			time.Sleep(time.Duration(10+rng.Intn(10)) * time.Second)
		}
	}
}
