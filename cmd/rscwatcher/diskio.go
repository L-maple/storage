package main

import (
	"fmt"
	"github.com/tommenx/storage/pkg/config"
	"os/exec"
	"strconv"
	"time"
)

type Diskio struct {
	pause int32   // watcher internal-times, second
	diskRead float32  // disk-read speed, K/s
	diskWrite float32  // disk-write speed, K/s
	tid	uint16
}

func (d *Diskio) Run() {
	pidStr :=  strconv.FormatInt(int64(d.tid), 10)

	cmd := exec.Command("/bin/sh",
						 "-c",
						 "sudo iotop -bkt -n 1 -p " + pidStr)

	// get the iotop's output
	if output, err := cmd.Output(); err != nil {
		fmt.Println(err.Error())
		return
	}else {
		// parse the diskRead and diskWrite
		fmt.Println(time.Now().Second())
		fmt.Println(string(output))
	}
}

/*
	function: return a Diskio object;
	parameters:
		pid: process id;
		pauseTime: pause time, second, for exmaple: 1s
 */
func NewDiskIO (pid uint16, pauseTime int32) *Diskio {
	return &Diskio{
		pause: pauseTime,
		tid: pid,
	}
}

func getPauseTime(path string) int32 {
	config.Init(path)
	return config.GetPauseTime()
}

func RunDiskIO(pid uint16) {
	// read pauseTime from config file
	pauseTime := getPauseTime("../../config.toml")

	// get the result
	diskio := NewDiskIO(uint16(pid), pauseTime)
	for ;; {
		diskio.Run()
	}
}

func main() {
	RunDiskIO(21344)
}
