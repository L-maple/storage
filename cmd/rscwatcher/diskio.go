package main

import (
	"fmt"
	"os/exec"
	//"time"
)

type Diskio struct {
	pause float32   // watcher internal-times, second
	diskRead float32  // disk-read speed, K/s
	diskWrite float64  // disk-write speed, K/s
	tid	uint16
}

func (*Diskio) Run() {
	cmd := exec.Command("/bin/sh", "-c", "sudo iotop -bk -n 2 -d 5 -p 15841")  // which iotop -> file location
	var output []byte
	var err error
	if output, err = cmd.Output(); err != nil {
		fmt.Println(err.Error())
		return
	}
	// parse the diskRead and diskWrite
	fmt.Println(string(output))
}

/*
	pid: process id;
	pauseTime: pause time, second, for exmaple: 1s
 */
func NewDiskIO (pid uint16, pauseTime float32) *Diskio {
	return &Diskio{
		pause: pauseTime,
		tid: pid,
	}
}

func main() {
	diskio := NewDiskIO(15841, 2)
	diskio.Run()
}
