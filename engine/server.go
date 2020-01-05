package engine

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type clientHandlerFunc func(rw *bufio.ReadWriter)

type boyomaServer struct {
	listener net.Listener
	handlers map[string]clientHandlerFunc
	lock     sync.RWMutex
}

func newServer(port int) (*boyomaServer, error) {
	return &boyomaServer{
		handlers: map[string]clientHandlerFunc{},
	}, nil
}

func (e *boyomaServer) addHandlerFunc(name string, f clientHandlerFunc) {
	e.lock.Lock()
	e.handlers[name] = f
	e.lock.Unlock()
}

func (e *boyomaServer) listen(port int) error {
	var err error
	e.listener, err = net.Listen("tcp", string(port))
	if err != nil {
		fmt.Println(err)
		return errors.Wrapf(err, "Cannot listen on port %d\n", port)
	}

	fmt.Println("StartServer listening on", e.listener.Addr().String())
	for {
		conn, err := e.listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection request:", err)
			continue
		}

		go e.handleConnection(conn)
	}
}

func (e *boyomaServer) handleConnection(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	for {
		cmd, err := rw.ReadString('\n')
		switch {
		case err == io.EOF:
			fmt.Println("Closing this connection")
			return
		case err != nil:
			fmt.Println("errorLogger reading command, got:"+cmd+"\n", err)
			return
		}

		cmd = strings.Trim(cmd, "\n ")
		log.Println(cmd)

		e.lock.Lock()
		handler, ok := e.handlers[cmd]
		e.lock.Unlock()
		if !ok {
			log.Println("Command " + cmd + " is not registered")
			return
		}

		handler(rw)
	}
}

func StartServer(port int) error {
	ep, err := newServer(port)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ep.addHandlerFunc("STRING", handleStrings)
	ep.addHandlerFunc("GOB", handleGOB)

	return ep.listen(port)
}

func handleGOB(rw *bufio.ReadWriter) {
	//var data complexData
	//
	//
	//dec := gob.NewDecoder(rw)
	//err := dec.Decode(&data)
	//if err != nil {
	//	log.Println("Error decoding GOB data:", err)
	//	return
	//}
	//
	//log.Printf("Outer complexData struct: \n%#v\n", data)
	//log.Printf("Inner complexData struct: \n%#v\n", data.C)
}

func handleStrings(rw *bufio.ReadWriter) {
	//s, err := rw.ReadString('\n')
	//if err != nil {
	//	log.Println("Cannot read from connection.\n", err)
	//}
	//s = strings.Trim(s, "\n ")
	//log.Println(s)
	//_, err = rw.WriteString("Thank you.\n")
	//if err != nil {
	//	log.Println("Cannot write to connection.\n", err)
	//}
	//err = rw.Flush()
	//if err != nil {
	//	log.Println("Flush failed.", err)
	//}
}
