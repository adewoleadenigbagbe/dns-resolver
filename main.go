package main

import (
	"math/rand"
)

type dnsMessage struct {
}

type header struct {
	ID      int16
	Flags   int16
	QDCOUNT int16
	ANCOUNT int16
	NSCOUNT int16
	ARCOUNT int16
}

type question struct {
	QNAME  string
	QTYPE  int16
	QCLASS int16
}

type answer struct {
}

func main() {
	h := header{
		ID:      int16(rand.Intn(65536)),
		Flags:   0x100, // standard query with recursion desired
		QDCOUNT: 1,
		ANCOUNT: 0,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}

	q := question{
		QNAME:  "",
		QTYPE:  1,
		QCLASS: 1,
	}
}
