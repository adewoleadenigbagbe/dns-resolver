# Dns-Resolver

A simple lightweight command line tool to resolve dns that can resolve the IP address for a host

DNS which means Domain Name system , when you type a url in your browser, some actions take place behind the scene

## Steps to resolving
* The browser check if it has the url in its DNS cache if it doesnt
* The operating system (OS) provides a function to check it own DNS cache if it doesnt
* The Os request for a server called a DNS Resolver (tool built) knows how to find the authoritative nameservers by sending a DNS query
* The authoritative nameservers are the servers where DNS records are actually stored. Ex : Google as a nameserver


For more information on dns resolver, visit this links below

* [Wizardzines](https://wizardzines.com/comics/cast-of-characters/) - A visual comic of how DNS resolver work by Julia Evans

* [Wikipedia](https://en.wikipedia.org/wiki/Domain_Name_System) - Wiki info

* [Rfc1035](https://datatracker.ietf.org/doc/html/rfc1035) - Rfc standard on how resolver needs to be built

## Usage

In this example we use google IP Address as the name server

Clone the Repository, check in on the root folder and run the main.go file, you see a list address printed on the cmd line

