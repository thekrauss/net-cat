package server

import "net"

const (
	ColorRed       = "\033[31m"
	ColorReset     = "\033[0m"
	MaxConnections = 10
)

type Server struct {
	IP                  string
	PORT                string
	Listener            net.Listener
	clients             map[net.Conn]string
	chatHistory         []string
	connectionSemaphore chan struct{}
}

func State(IP, PORT string) *Server {
	return &Server{
		IP:                  IP,
		PORT:                PORT,
		clients:             make(map[net.Conn]string),
		connectionSemaphore: make(chan struct{}, MaxConnections),
	}
}

func (server *Server) startConn() error {
	return nil
}

func (server *Server) Run() {

}

func (server *Server) addClient(clients net.Conn) {

}
