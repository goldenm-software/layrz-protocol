package clients

import (
	"errors"
	"log"
	"net"
	"regexp"
	"strconv"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/server"
)

type TcpComm struct {
	Host   string
	Port   int
	Ident  string
	Passwd string

	initialized   bool
	callback      *func(*any)
	conn          *net.Conn
	authenticated bool
	splitRegExp   *regexp.Regexp
	endRegExp     *regexp.Regexp
}

// New creates a new intance of LayrzProtocol using TCP communication
func (p *TcpComm) New(host string, port int, ident, password string) {
	p.Host = host
	p.Port = port
	p.Ident = ident
	p.Passwd = password

	p.initialized = true
	p.callback = nil
	p.authenticated = false

	p.splitRegExp = regexp.MustCompile(`(<[A-Z][a-z]>[\s\S]*?</[A-Z][a-z]>)`)
	p.endRegExp = regexp.MustCompile(`</[A-Z][a-z]>`)
}

// SetCallback sets the callback function to be called when a packet is received
func (p *TcpComm) SetCallback(callback func(*any)) error {
	if !p.initialized {
		return errors.New("tcp comm not initialized")
	}

	p.callback = &callback
	return nil
}

// Send sends a packet to the server
func (p *TcpComm) Send(packet any) error {
	if !p.initialized {
		return errors.New("tcp comm not initialized")
	}

	data, err := EncodeClientPacket(packet)
	if err != nil {
		return err
	}

	*data += "\r\n"

	log.Printf("Sending %s\n", *data)
	if _, err := (*p.conn).Write([]byte(*data)); err != nil {
		return err
	}

	return nil
}

// Connect connects to the server
func (p *TcpComm) Connect() error {
	if !p.initialized {
		return errors.New("tcp comm not initialized")
	}
	p.authenticated = false

	conn, err := net.Dial("tcp", net.JoinHostPort(p.Host, strconv.Itoa(p.Port)))
	if err != nil {
		return err
	}

	p.conn = &conn
	go p.listen()

	if err := p.Send(&client.PaPacket{
		Ident:    &p.Ident,
		Password: &p.Passwd,
	}); err != nil {
		return err
	}

	for attempts := 0; !p.authenticated; attempts++ {
		if attempts > 60 {
			return errors.New("authentication timeout")
		}
		log.Println("Waiting for authentication...")
		time.Sleep(time.Second)
	}
	return nil
}

// Listen for incoming data
func (p *TcpComm) listen() {
	buffer := make([]byte, 4096)
	messages := []byte{}

	for {
		n, err := (*p.conn).Read(buffer)
		if err != nil {
			log.Println("Connection closed:", err)
			panic(err)
		}

		messages = append(messages, buffer[:n]...)
		if p.endRegExp.Match(messages) {
			// Split messages
			messagesSplit := p.splitRegExp.FindAll(messages, -1)
			for _, message := range messagesSplit {
				log.Printf("Received message %s\n", message)
				packet, err := DecodeServerOutput(string(message))
				if err != nil {
					log.Printf("Error handling server output %s\n", err)
					break
				}

				switch (*packet).(type) {
				case *server.AsPacket:
					p.authenticated = true

				case *server.AuPacket: //nolint:staticcheck
					log.Println("Deprecated AuPacket...")

				default:
					if p.callback != nil {
						log.Println("Calling callback function...")
						(*p.callback)(packet)
					} else {
						log.Println("No callback function set...")
					}
				}
			}

			messages = []byte{}
		}
	}

}

// Close closes the connection
func (p *TcpComm) Close() error {
	if !p.initialized {
		return errors.New("tcp comm not initialized")
	}

	return (*p.conn).Close()
}
