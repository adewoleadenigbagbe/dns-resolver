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
	Id      uint16
	Flags   uint16
	Qdcount int16
	Ancount int16
	Nscount int16
	Arcount int16
}

type question struct {
	qtype  int16
	qclass int16
}

type record struct {
	Name     []byte
	Type     uint16
	Class    uint16
	Ttl      uint32
	RdLength uint16
	RData    []byte
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

	//fmt.Println(response)
	resolveDns(buf, response)
}

func resolveDns(buf *bytes.Buffer, data []byte) ([]record, error) {
	offset := 12
	_, err := buf.Write(data[:offset])
	if err != nil {
		return nil, err
	}

	h := &header{}
	binary.Read(buf, binary.BigEndian, h)
	fmt.Println("nscount : ", h.Nscount)
	var qr uint16 = 1 << 15

	//qr is set to 0 which indicate (query:0) when sending, name server sets it to 1 which indicate (response:1)
	if !(h.Flags&qr != 0) {
		return nil, fmt.Errorf("expected response is 1 but got %d", h.Flags&qr)
	}

	//with the number of question sent to the name server
	for i := 0; i < int(h.Qdcount); i++ {
		for data[offset] != 0 {
			offset += 1
		}
		//qclass + qtype is 32 bits (4 bytes), increment the offset by 4 plus 1 = 5 which is the dermacated byte when found
		offset += 5
	}

	buf.Reset()

	answerRecords := make([]record, h.Ancount)
	for i := 0; i < int(h.Ancount); i++ {
		//skip for the name
		offset += 2

		answerRecords[i].Type = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		answerRecords[i].Class = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		answerRecords[i].Ttl = binary.BigEndian.Uint32(data[offset : offset+4])
		offset += 4

		answerRecords[i].RdLength = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		answerRecords[i].RData = data[offset : offset+int(answerRecords[i].RdLength)]
		offset += int(answerRecords[i].RdLength)
	}

	if len(answerRecords) > 0 {
		buf.Reset()
		return answerRecords, nil
	}

	for i := 0; i < int(h.Nscount); i++ {

	}

	return answerRecords, nil
}

func parseRecordSection() {

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
		Id:      uint16(rand.Intn(65536)),
		Flags:   0x100, // standard query with recursion desired
		Qdcount: 1,
		Ancount: 0,
		Nscount: 0,
		Arcount: 0,
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
	//append with zero, this dermacate name to be resolved , in case you have multiple names
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
