package server

import (
	"bufio"
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

// Server représente le serveur de chat.

type Server struct {
	IP                  string
	PORT                string
	Listener            net.Listener
	clients             map[net.Conn]string
	chatHistory         []string
	connectionSemaphore chan struct{}
}

// State initialise un nouvel objet Server.

func State(IP, PORT string) *Server {

	return &Server{
		IP:                  IP,
		PORT:                PORT,
		clients:             make(map[net.Conn]string),
		connectionSemaphore: make(chan struct{}, MaxConnections),
	}
}

// Run lance la connexion du serveur et attend indéfiniment de nouveaux clients.
func (server *Server) Run() {
	if err := server.startCon(); err != nil {
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

// startCon initialise la connexion du serveur.
func (server *Server) startCon() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.IP, server.PORT))
	if err != nil {
		return err
	}
	server.Listener = listener
	fmt.Println(ColorRed, "[SERVER] : Serveur lancé ...", ColorReset)
	return nil
}

// addClient gère l'ajout d'un nouveau client au serveur.
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
		if len(server.clients) > 1 {
			server.sendChatHistory(client)
		}

		server.msgToAll(client, message)
		server.addLog(fmt.Sprintf("%s connecté depuis %s\n", server.clients[client], client.RemoteAddr()))
		server.receive(client)
	}

}

// sendChatHistory envoie l'historique du chat à un nouveau client
func (server *Server) sendChatHistory(client net.Conn) {
	for _, message := range server.chatHistory {
		server.sendToClient(client, (message))
	}
}

// gestionErreur gère les erreurs du serveur.
func (server *Server) gestionErreur(err error) {

	if err != nil {
		if server.Listener != nil {
			server.Listener.Close()
		}
		fmt.Println(ColorRed, "Serveur fermé", ColorReset)
		fmt.Println(err)
		os.Exit(2)
	}
}

// removeClient retire un client du serveur et informe les autres clients de sa déconnexion.
func (server *Server) removeClient(client net.Conn) {
	message := fmt.Sprintf("INFO : %s s'est déconnecté\n", server.clients[client])
	server.msgToAll(client, message)
	server.addLog(fmt.Sprintf("%s est déconnecté [total clients %d]\n", server.clients[client], len(server.clients)-1))
	delete(server.clients, client)
	client.Close()
}

// receive gère la réception de messages d'un client.
func (server *Server) receive(client net.Conn) {
	defer server.removeClient(client) // Appel à removeClient lorsque la boucle se termine

	buf := bufio.NewReader(client)

	for {
		message, err := buf.ReadString('\n')
		if err != nil {
			break // Sortir de la boucle si une erreur se produit
		}
		message = fmt.Sprintf("[%s]: %s", server.clients[client], message)
		server.msgToAll(client, message)
	}
}

// msgToAll envoie un message à tous les clients, sauf à l'émetteur.
func (server *Server) msgToAll(sender net.Conn, message string) {
	formattedMessage := fmt.Sprintf(" %s", (message))
	server.chatHistory = append(server.chatHistory, datetimeLine(formattedMessage))

	for client := range server.clients {
		if sender != client {
			server.sendToClient(client, message)
		}
	}
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

// isUsernameExists vérifie si un nom d'utilisateur existe déjà.
func (server *Server) isUsernameExists(username string, client net.Conn) bool {
	for _, existingUsername := range server.clients {
		if existingUsername == username && server.clients[client] != existingUsername {
			return true
		}
	}
	return false
}

// checkUsernameClient vérifie le nom d'utilisateur d'un client et l'ajoute au serveur.
func (server *Server) checkUsernameClient(client net.Conn) bool {
	username, err := server.catchClientUsername(client)
	if err != nil {
		return false
	}

	for server.isUsernameExists(username, client) {
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

func datetimeLine(text string) string {
	datetimeNow := time.Now().Format("02/01/2006 15:04:05")
	return fmt.Sprintf("[%s] %s", datetimeNow, text)
}

// addLog ajoute une ligne de journal avec horodatage.
func (server *Server) addLog(line string) {

	file, err := os.OpenFile("souche.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		fmt.Println(ColorRed, "ERREUR :", err, ColorReset)
		return
	}
	line = datetimeLine(line)
	_, err = file.WriteString(line)
	if err != nil {
		fmt.Println(ColorRed, "ERREUR :", err, ColorReset)
		return
	}
	fmt.Println(line)

	defer file.Close()
}

// sendToClient envoie un message à un client.
func (server *Server) sendToClient(client net.Conn, message string) {
	_, err := client.Write([]byte(message))
	if err != nil {
		fmt.Println(ColorRed, "ERREUR lors de l'envoi au client :", err, ColorReset)
	}
}
