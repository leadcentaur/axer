package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/common-nighthawk/go-figure"
	"github.com/gookit/color"
	"github.com/miekg/dns"
)

var wg sync.WaitGroup
var rg sync.WaitGroup
var rd sync.WaitGroup

func readDn(ch chan string, fname string) {
	defer rd.Done()
	file, err := os.Open(fname)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dn := scanner.Text()
		fmt.Println(dn)
		ch <- dn
	}
	close(ch)
	if err := scanner.Err(); err != nil {
		color.Red.Print(err)
	}
}

func nsLookup(ch chan string, domain string, chB chan bool) {
	defer wg.Done()
	r, err := net.LookupNS(domain)
	if err != nil {
		fmt.Println(err)
		chB <- false
	} else {
		chB <- true
		for _, ns := range r {
			l := len(ns.Host)
			//annoying check we have to do :(
			if ns.Host[l-1] == '.' {
				ch <- strings.TrimSuffix(ns.Host, ".")
			} else {
				ch <- ns.Host
			}
		}
	}
}

func doAXFR(ch chan string, respch chan string, domain string) {
	defer rg.Done()
	for ns := range ch {
		fmt.Printf("Attempting AXFR, NS Len: %d @%s %s\n", len(ch), ns, domain)

		tran := new(dns.Transfer)
		msg := new(dns.Msg)
		msg.SetAxfr(dns.Fqdn(domain))
		reqNs := net.JoinHostPort(ns, "53")

		respChan, err := tran.In(msg, reqNs)
		if err != nil {
			color.Red.Print(os.Stderr, "An Error has occured: ", err)
			recover()
		}
		//dump everything
		for envelope := range respChan {
			fmt.Println("Attempting AXFR on Domain:", domain)
			if envelope.Error != nil {
				color.Red.Print("Transfer Failed - ")
				log.Printf("%s %s\n", envelope.Error, domain)

			} else {
				color.Green.Print("\t\t Transfer sucessful \n")
				for _, rr := range envelope.RR {
					fmt.Println(rr.Header().String())
					switch v := rr.(type) {
					case *dns.A:
						fmt.Println(v.Header().Name)
					case *dns.AAAA:
						fmt.Println(v.Header().Name)
					default:
					}
				}
			}
		}
	}
}

var nsParsed string
func main() {

	var inputStr, fpath string

	fig := figure.NewFigure("Axer v1.2", "", true)
	fig.Print()
	color.Green.Print("\t\t# Coded by Brenden #\n")

	flag.StringVar(&fpath, "f", "deg", "Read a list of domain name from a file.")
	flag.Parse()

	tail := flag.Args()

	boolCh := make(chan bool)
	nservers := make(chan string)
	respch := make(chan string)
	domainsd := make(chan string)

	if len(tail) == 0 || fpath =="deg" {
		if len(tail) == 1 {
			fpath = tail[0]
		}
		rd.Add(1)
		go readDn(domainsd, fpath)
		for dn := range domainsd {
			if string(dn) != "" {
				wg.Add(1)
				go nsLookup(nservers, dn, boolCh)
				readb := <-boolCh
				if readb == true {
					rg.Add(1)
					go doAXFR(nservers, respch, dn)
				}
			}
		}
		wg.Wait()
		close(nservers)
		rg.Wait()
		close(respch)
		rd.Wait()
	}
}
