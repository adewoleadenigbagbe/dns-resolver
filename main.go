package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
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
	Name     string
	Type     uint16
	Class    uint16
	Ttl      uint32
	RdLength uint16
	RData    []byte
}

func main() {
	buf := new(bytes.Buffer)
	var records []string
	records, err := resolveDns(buf, "198.41.0.4", ":53")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(records)
}

func resolveDns(buf *bytes.Buffer, nameServer string, port string) ([]string, error) {
	query, err := buildQuery(buf)
	if err != nil {
		return nil, err
	}

	data, err := sendQuery(query, nameServer, port)
	if err != nil {
		return nil, err
	}

	offset := 12
	_, err = buf.Write(data[:offset])
	if err != nil {
		return nil, err
	}

	h := &header{}
	binary.Read(buf, binary.BigEndian, h)

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

	var answerRecords []string
	for i := 0; i < int(h.Ancount); i++ {
		r, newOffset := parseRecordSection(buf, data, offset, h.Qdcount)
		if r.Type == 1 {
			var ip string
			for _, r := range r.RData {
				ip += strconv.Itoa(int(r)) + "."
			}
			ip = ip[:len(ip)-1]
			answerRecords = append(answerRecords, ip)
		}
		offset = newOffset
	}

	if len(answerRecords) > 0 {
		return answerRecords, nil
	}

	var nsRecords []string
	for i := 0; i < int(h.Nscount); i++ {
		r, newOffset := parseRecordSection(buf, data, offset, h.Qdcount)
		if r.Type == 2 {
			var ip string
			for _, r := range r.RData {
				ip += strconv.Itoa(int(r)) + "."
			}
			ip = ip[:len(ip)-1]
			nsRecords = append(nsRecords, ip)
		}
		offset = newOffset
	}

	for i := 0; i < int(h.Arcount); i++ {
		r, newOffset := parseRecordSection(buf, data, offset, h.Qdcount)
		if r.Type == 1 {
			var ip string
			for _, r := range r.RData {
				ip += strconv.Itoa(int(r)) + "."
			}
			ip = ip[:len(ip)-1]
			return resolveDns(buf, ip, ":53")
		}
		offset = newOffset
	}

	buf.Reset()

	return nil, nil
}

func parseRecordSection(buf *bytes.Buffer, data []byte, offset int, questionCount int16) (record, int) {
	var r record
	//check if there is a message compression which is 2 bytes data by checking whether it has a pointer
	//Ex: if the binary format is 110001111, the first two bit start with 11 makes it a pointer , hence a message compression
	b := fmt.Sprintf("%b", data[offset])
	pointer := b[:2]
	if pointer == "11" {
		b = b[2:] + fmt.Sprintf("%b", data[offset+1])
		qNameOffset, _ := strconv.ParseInt(b, 2, 64)

		var name string
		for i := 0; i < int(questionCount); i++ {
			for data[qNameOffset] != 0 {
				if data[qNameOffset] > 47 && data[qNameOffset] < 123 {
					name += string(data[qNameOffset])
				} else {
					name += "."
				}
				qNameOffset += 1
			}
			r.Name = name[1:]
			qNameOffset += 5
		}

	} else {
	}

	offset += 2

	r.Type = binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2

	r.Class = binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2

	r.Ttl = binary.BigEndian.Uint32(data[offset : offset+4])
	offset += 4

	r.RdLength = binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2

	r.RData = data[offset : offset+int(r.RdLength)]
	offset += int(r.RdLength)

	buf.Reset()
	return r, offset
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

func sendQuery(query string, ipAddress string, port string) ([]byte, error) {
	p := make([]byte, 1024)
	conn, err := net.Dial("udp", ipAddress+port)
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
