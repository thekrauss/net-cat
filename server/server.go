package server

import (
	"fmt"
	"net"
	"os"
)

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

	listener, err := net.Listen("tpc", fmt.Sprintf("%s:%s", server.IP, server.PORT))
	if err != nil {
		return err
	}
	server.Listener = listener
	fmt.Println(ColorRed, "[SERVER] : Serveur lancé ...", ColorReset)
	return nil

}

func (server *Server) Run() {
	if err := server.startConn(); err != nil {
		server.gestionErreur(err)
	}
	for {
		client, err := server.Listener.Accept()
		if err != nil {
			server.gestionErreur(err)
		}
		go server.addClient(client)
	}

}

func (server *Server) addClient(client net.Conn) {
	if len(server.clients) >= MaxConnections {
		server.sendToClient(client, "ERREUR : Nombre maximal de connexions atteint. Veuillez réessayer plus tard.\n")
		client.Close()
	} else {
		if !server.checkUsernameClient(client) {
			return
		}

		message := fmt.Sprintf("INFO : %s a rejoint le serveur\n", server.clients[client])
		server.sendToClient(client, "Vous pouvez commencer la discussion avec les autres invités...\n\n")

		// Envoyer l'historique du chat au nouveau client

		if (server.clients) > 1 {
			server.sendChatHistory(client)
		}

		server.msgToAll(client, message)
		server.addLog(fmt.Sprintf("%s connecté depuis %s\n", server.clients[client], client.RemoteAddr()))
		server.receive(client)
	}
}

func (server *Server) gestionErreur(err error) {
	if err != nil {
		server.Listener.Close()
	}
	fmt.Println(ColorRed, "Server closes", ColorReset)
	fmt.Println(err)
	os.Exit(2)
}

func (server *Server) sendToClient(client net.Conn, message string) {
	_, err := client.Write([]byte(message))
	if err != nil {
		fmt.Println(ColorRed, "ERREUR lors de l'envoi au client :", err, ColorReset)
	}
}
func (server *Server) checkUsernameClient(client net.Conn) bool {
	username, err := server.catchClientUsername(client)
	if err != nil {
		return false
	}

	for server.isUsernameExist(username, client) {
		server.sendToClient(client, "ERREUR : Votre nom d'utilisateur existe déjà sur le serveur \n")
		username, err := server.catchClientUsername(client)
		if err != nil {
			return false
		}
	}

	server.clients[client] = username
	server.sendToClient(client, "[SUCCÈS] Vous êtes connecté avec succès !\n")
	return true
}
