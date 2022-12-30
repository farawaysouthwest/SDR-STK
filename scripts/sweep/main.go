package main

import (
	"github.com/alexflint/go-arg"
)

type Args struct {
	Options     string  `help:"Input rtl_power options as a string." arg:"-o"`
	Pipe        bool    `help:"No rtl_power options, run in pipe mode." arg:"-i"`
	Server_Ip   string  `help:"Server IP address (default: localhost)" arg:"-s" default:"127.0.0.1"`
	Server_Port string  `help:"Server Port (default: 1236)" arg:"-p" default:"1236"`
	Client_Ip   string  `help:"Client IP address (default: localhost)" arg:"-c" default:"127.0.0.1"`
	Client_Port string  `help:"Client Port (default: 1234)" arg:"-a" default:"1234"`
	Db_Limit    float64 `help:"Set dBm peak detect limit (default: 0dBm)" arg:"-d" default:"0"`
	Logging     bool    `help:"Enable frequency and dBm logging in the console" arg:"-v"`
}

func main() {

	var args Args
	arg.MustParse(&args)

	forward, err := NewForward(&args)
	if err != nil {
		panic(err.Error())
	}

	listener, err := NewListener(&args, forward)
	if err != nil {
		panic(err.Error())
	}

	server := NewServer(&args, listener)

	err = server.Start()
	if err != nil {
		panic("failed to start server")
	}

}
