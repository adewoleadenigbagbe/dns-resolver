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
	var (
		err      error
		query    string
		response []byte
	)
	buf := new(bytes.Buffer)
	query, err = buildQuery(buf)
	if err != nil {
		log.Fatal(err)
	}

	response, err = sendQuery(query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

func buildQuery(buf *bytes.Buffer) (string, error) {
	var (
		err error
	)
	h, err := encodeHeader(buf)

	if err != nil {
		return "", nil
	}

	q, err := encodeQuestion(buf)
	if err != nil {
		return "", nil
	}

	return h + q, nil
}

func sendQuery(query string) ([]byte, error) {
	p := make([]byte, 1024)
	conn, err := net.Dial("udp", "8.8.8.8:53")
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(conn, query)
	_, err = bufio.NewReader(conn).Read(p)

	if err != nil {
		return nil, err
	}
	return p, nil
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
