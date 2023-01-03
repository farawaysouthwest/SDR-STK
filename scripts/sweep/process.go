package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func CreateRTLProcess(optionArgs []string, wg *sync.WaitGroup) *[]byte {

	options := strings.Join(optionArgs, " ")

	out, err := exec.Command("rtl_power", options, options).Output()
	if err != nil {
		fmt.Println("Error starting rtl_power instance", err.Error())
	}

	wg.Done()

	return &out
}

func CreateRTL_TCP(args *Args, wg *sync.WaitGroup) {
	cmd := exec.Command("rtl_tcp", " -a "+args.Client_Ip, " -p "+args.Client_Port)
	fmt.Println("running commend: ", strings.Join(cmd.Args, ""))

	out, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error starting rtl_tcp instance", err.Error())
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("started rtl_tcp process")

	wg.Done()

	outputBuff := make([]byte, 10)
	_, err = out.Read(outputBuff)
			if err != nil {
				fmt.Println(err)
			}

	go func() {
		for {
			time.Sleep(time.Duration(1) * time.Second)
			fmt.Println(outputBuff)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error starting rtl_tcp instance", err.Error())
	}

}
