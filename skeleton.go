// =============================================
// file: server/main.go
// =============================================
package main

import "tcp-chat/server" // TODO: replace with your actual module path

func main() {
    // entry point – запускаем TCP‑сервер на порту 8080
    server.StartServer(":8080")
}

// =============================================
// file: server/server.go
// =============================================
package server

import "net"

// Server инкапсулирует состояние и логику чата на стороне сервера.
//   - listener   – TCP‑слушатель
//   - clients    – мапа активных клиентов (соединение → имя)
//   - broadcast  – общий канал для рассылаемых сообщений
//
// StartServer(addr) – публичная точка входа; настраивает listener,
//                    инициализирует структуры данных и запускает цикл Accept.
// broadcastMessages() – горутина, читающая из broadcast и отправляющая
//                       строку всем клиентам из clients.
// HandleClient() – горутина на каждого клиента: читает имя, регистрирует,
//                  затем в цикле читает сообщения и кладёт их в broadcast.

type Server struct {
    listener  net.Listener
    clients   map[net.Conn]User
    broadcast chan string
}

// StartServer запускает TCP‑сервер.
func StartServer(addr string) {
    // TODO: listen on addr
    // TODO: create Server instance & launch broadcastMessages goroutine
    // TODO: accept loop – для каждого conn go HandleClient(conn, srv)
}

// broadcastMessages рассылает входящие строки всем подключённым клиентам.
func (s *Server) broadcastMessages() {
    // TODO: range over s.broadcast и делать conn.Write()
}

// HandleClient обслуживает одного клиента.
func (s *Server) HandleClient(conn net.Conn) {
    // TODO: запросить имя, зарегистрировать conn в s.clients
    // TODO: в цикле читать строки и отправлять в s.broadcast
    // TODO: обработать обрыв соединения, удалить клиента из s.clients
}

// =============================================
// file: server/client_handler.go
// =============================================
package server

import "net"

// (вариант) Выделенный обработчик клиента, если не хочется держать метод у Server.
// Здесь можно вынести часть логики из Server.HandleClient.
func HandleClient(conn net.Conn, srv *Server) {
    // TODO: read name → srv.clients
    // TODO: loop read → srv.broadcast <- fmt.Sprintf("[%s]: %s", name, msg)
}

// =============================================
// file: client/main.go
// =============================================
package main

import (
    "fmt"
    "net"
    "tcp-chat/client" // TODO: заменить модулем
)

func main() {
    // TODO: подключиться к серверу
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("не удалось подключиться:", err)
        return
    }
    defer conn.Close()

    client.StartChat(conn) // запускаем клиентскую логику
}

// =============================================
// file: client/chat.go
// =============================================
package client

import "net"

// StartChat запускает две горутины:
//   1. readMessages – читает данные с сервера и печатает в терминал
//   2. writeMessages – читает пользовательский ввод и пишет на сервер
func StartChat(conn net.Conn) {
    // TODO: go readMessages(conn)
    // TODO: writeMessages(conn)
}

// readMessages читает строки из conn и выводит их.
func readMessages(conn net.Conn) {
    // TODO
}

// writeMessages читает строки из stdin и отправляет их серверу.
func writeMessages(conn net.Conn) {
    // TODO
}

// =============================================
// file: common/types.go
// =============================================
package common

// Message – опциональный формат для обмена (если решишь кодировать JSON).
// Пока можно передавать голые строки; структура пригодится для расширений.

type Message struct {
    Sender  string `json:"sender"`
    Content string `json:"content"`
}
