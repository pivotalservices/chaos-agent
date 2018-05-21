package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	doSpike = kingpin.Flag("do-spike", "open an irresponsible number of files").Bool()
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.LUTC)
	kingpin.Parse()

	log.Println(os.Getpid())

	if *doSpike {
		log.Println("Time to do some pretty intensive computing!")
		log.Printf("Logs will be located in %s/security-report*\n", os.TempDir())

		quit := make(chan os.Signal, 1)
		done := make(chan bool, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			for {
				select {
				case <-quit:
					log.Println("Finishing up work...")
					done <- true
					return
				default:
					for j := 0; j < 10; j++ {
						f, _ := ioutil.TempFile("", "security-report")
						f.Write([]byte("Couldn't connect to host!"))
						defer func() {
							f.Close()
							log.Println("File closed!")
						}()
					}

					time.Sleep(1 * time.Second)
					log.Printf("There are %d open files\n", countOpenFiles())
				}
			}
		}()

		<-done
	}
}

func countOpenFiles() int {
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -p %v", os.Getpid())).Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	return len(lines) - 1
}
