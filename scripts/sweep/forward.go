package main

import (
	"C"
	"fmt"
	"net"
)
import (
	"sync"
)

type forward struct {
	connection net.Conn
}

type Forward interface {
	Send(buffer *[]byte) error
	Close()
}

func NewForward(args *Args) (Forward, error) {

	var wg sync.WaitGroup
	wg.Add(1)

	go CreateRTL_TCP(args, &wg)

	wg.Wait()
	socketConn, err := net.Dial("tcp", args.Client_Ip+":"+args.Client_Port)
	if err != nil {
		return nil, err
	}

	fmt.Println("connected to client device on", socketConn.LocalAddr().String())

	return forward{
		connection: socketConn,
	}, nil

}

func (r forward) Send(buffer *[]byte) error {

	dataBytes, err := r.connection.Write(*buffer)
	if err != nil {
		return err
	}

	fmt.Println("sent ", dataBytes, " bytes")
	return nil
}

func (r forward) Close() {
	err := r.connection.Close()
	if err != nil {
		panic(err)
	}
}
