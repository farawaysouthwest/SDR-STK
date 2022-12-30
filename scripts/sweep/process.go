package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func CreateRTLProcess(optionArgs []string) *[]byte {

	options := strings.Join(optionArgs, " ")

	out, err := exec.Command("rtl_power", options, options).Output()
	if err != nil {
		fmt.Println("Error starting rtl_power instance", err.Error())
	}

	return &out
}

func CreateRTL_TCP(args *Args, wg *sync.WaitGroup) {
	cmd := exec.Command("rtl_tcp", "-a"+args.Client_Ip, "-p"+args.Client_Port)

	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting rtl_tcp instance", err.Error())
	}
	fmt.Println("started rtl_tcp")

	wg.Done()

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error starting rtl_tcp instance", err.Error())
	}

	defer cmd.Process.Kill()
}
