package main

const (
	DEFAULT_TOP_FREQ = 0
	DEFAULT_TOP_DBM  = "0"
)

type topFreq struct {
	freq float64
	dBm  float64
}

type server struct {
	args Args

	ingressSocket Listener
}

type Server interface {
	Start() error
}

func NewServer(args *Args, connection Listener) Server {

	return &server{
		args:          *args,
		ingressSocket: connection,
	}
}

func (r *server) Start() error {

	err := r.ingressSocket.Listen()
	if err != nil {
		return err
	}
	return nil
}
