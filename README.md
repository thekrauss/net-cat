
# TCP-Chat (NetCat Recreation)

This project is a recreation of the `NetCat` command-line utility, built using Go. The main objective of this project was to create a server-client architecture for a TCP group chat, capable of handling multiple clients simultaneously.

## Objectives

The project implements the following key features:

- **TCP connection** between a server and multiple clients.
- Each client must **provide a name** before joining the chat.
- Limit of **10 connections** at a time.
- Clients can **send messages** to the group chat.
- Messages are **timestamped** and tagged with the client's name (e.g., `[2024-09-15 14:00:00][client_name]: message`).
- New clients receive **all previous chat history** upon joining.
- Clients are notified when a **new client joins** or when a client **leaves**.
- **No empty messages** are broadcasted.
- If no port is specified, the default port used is `8989`.

## How It Works

1. Clone the repository:
    ```bash
        https://github.com/thekrauss/net-cat.git
    ```

### Server

To start the server, use the following commands:

```bash
$ go run main.go server
Listening on port: 8989
```

Or specify a custom port:

```bash
$ go run . 2525
Listening on port: 2525
```

### Client

To connect to the server, use NetCat:

```bash
$ go run main.go client
```

When a client connects, they are greeted with a welcome message and asked for their name:

```
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]:
```

Once the client provides their name, they can start chatting with others. For example:

```
[2024-09-15 14:00:00][Kevin]: Hello, everyone!
[2024-09-15 14:01:00][Alice]: Hi, Kevin!
```

### Handling Clients

- When a **new client** joins, the server broadcasts the event to all clients.
- When a client **leaves**, the server informs the remaining clients.
- **Errors** are handled both on the server and client side to ensure the system is stable.


## Project Structure

- **Server:** The server is responsible for listening on a specific port and accepting incoming connections.
- **Client:** Each client can connect to the server, send messages, and receive messages from other clients.
- **Concurrency:** The project uses **Goroutines** to manage multiple client connections simultaneously.
- **Synchronization:** **Channels** or **Mutexes** are used to manage communication between the Goroutines.

## Future Improvements

- Implement a **Terminal UI** using the `gocui` package.
- Save chat logs to a file.
- Support for multiple group chats.

## Learning Outcomes

This project helped in understanding:

- **Go concurrency** using Goroutines.
- Managing **TCP connections** and communication.
- Using **channels** and **mutexes** for thread safety.
- Socket programming and network protocols.

## License

This project is licensed under the MIT License.
