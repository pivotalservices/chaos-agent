package main

import (
	"log"
	"math/rand"
	"time"

	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	probabilityOfFork  = kingpin.Flag("fork-probability", "The probability (0 - 1) that this app will spawn child processes").Short('f').Default("0.5").Float()
	probabilityOfSpike = kingpin.Flag("spike-probability", "The probability (0 - 1) that this app will open a LOT of file handles").Short('s').Default("0.5").Float()
)

func main() {
	workerProcess := filepath.Dir(os.Args[0]) + "/worker"

	kingpin.Parse()
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.LUTC)

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// for {

	forkValue := rng.Float64()
	spikeValue := rng.Float64()

	doFork := (forkValue <= *probabilityOfFork)
	doSpike := (spikeValue <= *probabilityOfSpike)

	if doFork {
		log.Print("Spawning worker processes ...")
		args := make([]string, 0, 1)

		if doSpike {
			args = append(args, "--do-spike")
		}

		for i := 0; i < 10; i++ {
			cmd := exec.Command(workerProcess, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		}
	}

	// 	time.Sleep(time.Duration(30+rng.Intn(30)) * time.Second)
	// }
}
