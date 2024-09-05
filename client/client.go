package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	ColorRed   = "\033[31m"
	ColorReset = "\033[0m"
	ColorBlue  = "\033[34m"
)

// WaitGroup pour synchroniser les goroutines
var wg sync.WaitGroup

const LengthUsername int = 20

type Client struct {
	IP          string
	PORT        string
	Username    string
	Conn        net.Conn
	IsConnected bool
}

// State initialise et renvoie un Client avec l'IP et le PORT spécifiés
func State(IP string, PORT string) Client {
	var client Client
	client.IP, client.PORT = IP, PORT
	return client
}

// usernameHandle gère la saisie et la validation du nom d'utilisateur
func (client *Client) usernameHandle() {
	for {
		Username := client.getUsernameInput()
		lengthUsername := len(strings.TrimSuffix(Username, "\n"))

		if lengthUsername > LengthUsername || lengthUsername == 0 {
			fmt.Println("ERREUR : Votre nom d'utilisateur ne doit pas être vide et ne doit pas dépasser", LengthUsername, "caractères")
			continue
		}

		client.Username = strings.TrimSuffix(Username, "\n") // Stocke le nom d'utilisateur dans la structure Client
		client.Conn.Write([]byte(Username))
		message, _ := client.read()
		fmt.Print((message))

		if strings.Contains(message, "SUCCÈS") {
			break
		}
	}
}

// getUsernameInput obtient la saisie du nom d'utilisateur de l'utilisateur
func (client *Client) getUsernameInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(ColorBlue, "Entrez votre nom d'utilisateur : ", ColorReset)
	username, err := reader.ReadString('\n')

	client.check(err)
	return username
}

// connect établit une connexion avec le serveur
func (client *Client) connect() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", client.IP, client.PORT))
	client.check(err)
	fmt.Println("Connexion à", conn.RemoteAddr(), "SERVEUR en cours...")
	client.IsConnected = true
	client.Conn = conn
}

// check vérifie les erreurs et gère la déconnexion
func (client *Client) check(err error) {
	if err != nil {
		client.IsConnected = false
		if client.Conn != nil {
			client.Conn.Close()
		}
		fmt.Println(err)
		fmt.Println(ColorRed, "Vous êtes maintenant déconnecté !", ColorReset)
		os.Exit(2)
	}
}

// send gère l'envoi de messages au serveur
func (client *Client) send() {
	defer wg.Done()
	for {
		reader := bufio.NewReader(os.Stdin)
		//fmt.Print(ColorBlue, "Message: ", ColorReset)
		message, err := reader.ReadString('\n')
		if !client.IsConnected {
			break
		}
		client.check(err)
		client.Conn.Write([]byte(message))

		// Efface la ligne de saisie
		fmt.Print("\033[1A") // Déplacer le curseur d'une ligne vers le haut
		fmt.Print("\033[K")  // Effacer la ligne

		fmt.Print(datetimeSendLine(message))
	}
}

// read lit les messages du serveur
func (client *Client) read() (string, error) {
	messageBuffer := make([]byte, 4096)
	length, err := client.Conn.Read(messageBuffer)
	if err != nil {
		fmt.Println(ColorRed, "[INFO] Le serveur est hors ligne, appuyez sur Entrée pour fermer la session", ColorReset)
		client.IsConnected = false
	}
	message := string(messageBuffer[:length])

	return message, err
}

// receive gère la réception de messages du serveur
func (client *Client) receive() {
	defer wg.Done()
	for {
		message, err := client.read()
		if !client.IsConnected {
			break
		}
		if err != nil {
			fmt.Println(ColorRed, "[INFO] Le serveur est hors ligne, appuyez sur Entrée pour fermer la session", ColorReset)
			client.IsConnected = false
			break
		}
		fmt.Print(datetimeLine(message))
	}
}

// datetimeLine retourne une chaîne avec le texte et l'horodatage
func datetimeLine(text string) string {
	datetimeNow := time.Now().Format("02/01/2006 15:04:05")
	return fmt.Sprintf("[%s] %s", datetimeNow, text)
}

func datetimeSendLine(text string) string {
	datetimeNow := time.Now().Format("02/01/2006 15:04:05")
	return fmt.Sprintf("[%s] [Vous]: %s", datetimeNow, text)
}

// Run démarre le client en établissant la connexion, gérant le nom d'utilisateur et les goroutines d'envoi/réception
func (client *Client) Run() {
	client.connect()
	client.usernameHandle()
	wg.Add(2)
	go client.send()
	go client.receive()
	wg.Wait()
	fmt.Println("Vous êtes maintenant déconnecté !")
}
