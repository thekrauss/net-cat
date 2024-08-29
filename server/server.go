package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
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
		username, err = server.catchClientUsername(client)
		if err != nil {
			return false
		}
	}

	server.clients[client] = username
	server.sendToClient(client, "[SUCCÈS] Vous êtes connecté avec succès !\n")
	return true
}

// catchClientUsername récupère le nom d'utilisateur d'un client.
func (server *Server) catchClientUsername(client net.Conn) (string, error) {

	var usernameBuffer [4096]byte
	length, err := client.Read(usernameBuffer[:])
	if err != nil {
		server.addLog(fmt.Sprintf("Le client depuis %s a interrompu la saisie du nom d'utilisateur\n", client.RemoteAddr()))
		return "", errors.New("saisie du client interrompue")
	}
	return strings.TrimSuffix(string(usernameBuffer[:length]), "\n"), nil
}

func (server *Server) isUsernameExist(username string, client net.Conn) bool {
	for _, existUsername := range server.clients {
		if existUsername == username && server.clients[client] != existUsername {
			return true
		}
	}
	return false
}

func (server *Server) addLog(line string) {
	file, err := os.OpenFile("souche.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		fmt.Println(ColorRed, "ERREUR :", err, ColorReset)
		return
	}
	line = dateTimeLine(line)
	_, err = file.WriteString(line)
	if err != nil {
		fmt.Println(ColorRed, "ERREUR :", err, ColorReset)
		return
	}
	fmt.Println(line)
	defer file.Close()
}

func (server *Server) sendChatHistory(client net.Conn) {
	for _, message := range server.chatHistory {
		server.sendToClient(client, (message))
	}
}

func dateTimeLine(text string) string {
	datetimeNow := time.Now().Format("02/01/2006 15:04:05")
	return fmt.Sprintf("[%s] %s", datetimeNow, text)
}
