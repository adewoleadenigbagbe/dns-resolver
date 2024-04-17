# Dns-Resolver

A simple lightweight command line tool to resolve dns that can resolve the IP address for a host

DNS which means Domain Name system

When you type a url in your browser, a some action take place behind the scene

## STEPS TO RESOLVING
The broswer check if it has the url in its DNS cache if it doesnt
The operating system (OS) provides a function to check it own DNS cache if it doesnt
The Os request for a server called a DNS Resolver (tool built) also has its cache if not present it cache , knows how to find the authoritative nameservers by sending a DNS query
Tha authoritative nameservers are the servers where DNS records are actually stored.


For more information on dns resolver, visit this links below

[text](https://wizardzines.com/comics/cast-of-characters/) - A visual comic of how DNS resolver work
[text](https://en.wikipedia.org/wiki/Domain_Name_System)
[text](https://datatracker.ietf.org/doc/html/rfc1035)

## Usage
In this example we use google IP Address as the name server
Clone the Repository, check in on the root folder and run the main.go file, you see a list address printed on the cmd line

