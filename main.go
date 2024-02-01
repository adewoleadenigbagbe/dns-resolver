package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
)

const (
	nameserver = "dns.google.com"
)

type header struct {
	id      uint16
	flags   int16
	qdcount int16
	ancount int16
	nscount int16
	arcount int16
}

type question struct {
	qtype  int16
	qclass int16
}

type answer struct {
}

func main() {
	buf := new(bytes.Buffer)
	query := buildQuery(buf)
	sendQuery(query)
}

func buildQuery(buf *bytes.Buffer) string {
	var (
		err error
	)
	h, err := encodeHeader(buf)

	if err != nil {
		log.Fatal(err)
	}

	q, err := encodeQuestion(buf)
	if err != nil {
		log.Fatal(err)
	}

	return h + q
}

func sendQuery(query string) {
	p := make([]byte, 1024)
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	fmt.Fprintf(conn, query)
	_, err = bufio.NewReader(conn).Read(p)

	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()
}

func encodeHeader(buf *bytes.Buffer) (string, error) {
	h := &header{
		id:      uint16(rand.Intn(65536)),
		flags:   0x100, // standard query with recursion desired
		qdcount: 1,
		ancount: 0,
		nscount: 0,
		arcount: 0,
	}

	err := binary.Write(buf, binary.BigEndian, h)
	if err != nil {
		return "", err
	}

	s := buf.String()
	buf.Reset()
	return s, nil
}

func encodeQuestion(buf *bytes.Buffer) (string, error) {
	var encoded string
	for _, part := range strings.Split(nameserver, ".") {
		b := byte(uint8(len(part)))
		encoded += string(b) + part
	}
	encoded += string(byte(0))

	q := &question{
		qtype:  1,
		qclass: 1,
	}

	err := binary.Write(buf, binary.BigEndian, q)
	if err != nil {
		return "", err
	}
	encoded += buf.String()

	buf.Reset()
	return encoded, nil
}
