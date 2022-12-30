package main

import (
	"fmt"
	"math"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	DEFAULT_BUFFER_SIZE      = 4096
	DEFAULT_DELAY            = 0.0001
	SET_FREQUENCY       byte = 0x01
)

type listener struct {
	args       *Args
	bufferSize int
	delay      float64
	optionArgs []string
	topFreq    topFreq
	server     net.Listener
	forward    Forward
}

type Listener interface {
	// test
	Listen() error

	//test
	Close() error
}

func NewListener(args *Args, forward Forward) (Listener, error) {

	server, err := net.Listen("tcp", args.Server_Ip+":"+args.Server_Port)
	if err != nil {
		return nil, err
	}

	return &listener{
		args:       args,
		server:     server,
		bufferSize: DEFAULT_BUFFER_SIZE,
		delay:      DEFAULT_DELAY,
		optionArgs: strings.Split(args.Options, " "),
		topFreq:    topFreq{freq: 0, dBm: 0},
		forward:    forward,
	}, nil
}

func (r *listener) Close() error {
	err := r.server.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *listener) Listen() error {

	for {
		// debounce
		time.Sleep(time.Duration(r.delay))

		connection, err := r.server.Accept()
		if err != nil {
			return err
		}

		fmt.Println("connection request accepted")
		go r.processClient(connection)
	}
}

func (r *listener) processClient(connection net.Conn) {

	out, err := exec.Command(strings.Join(r.optionArgs, " ")).Output()
	if err != nil {
		fmt.Println("Error starting power instance", err.Error())
	}

	defer connection.Close()

	lowFreq := float64(0)

	for {
		if out != nil {

			input := strings.Split(string(out), ", ")

			step := input[4]
			db := input[6:]

			maxIndex := r.findMaxIndex(db)
			maxValue := math.Max(float64(maxIndex), float64(len(db)))

			freq, startFreq := r.calcFreq(input[2], maxIndex, step)

			if startFreq >= lowFreq {
				lowFreq = startFreq

				cmd := make([]byte, 0)
				cmd = append(cmd, []byte(">BI")...)
				cmd = append(cmd, SET_FREQUENCY)
				cmd = append(cmd, byte(r.topFreq.freq))

				r.forward.Send(&cmd)

				r.topFreq.dBm = r.args.Db_Limit
			}

			if maxValue > r.topFreq.dBm {
				r.topFreq.freq = freq
				r.topFreq.dBm = maxValue
			}
		}
	}
}

func (r *listener) calcFreq(startFreq string, maxIndex int64, step string) (float64, float64) {

	sFreq, err := strconv.ParseFloat(startFreq, 64)
	if err != nil {
		fmt.Println("power input invalid")
	}

	mFreq, err := strconv.ParseFloat(step, 64)
	if err != nil {
		fmt.Println("power input invalid")
	}

	return sFreq + float64(maxIndex)*mFreq, sFreq
}

func (r *listener) findMaxIndex(array []string) int64 {

	var highest int64

	for i, v := range array {
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			fmt.Println("error finding max index", err.Error())
		}

		if parsed > highest {
			highest = int64(i)
		}
	}

	return highest
}
