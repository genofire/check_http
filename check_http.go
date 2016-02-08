package main

import (
		"io/ioutil"
		"net"
		"strings"
		"os"
		"sync"
		"log"
		"regexp"
		"crypto/tls"
//		"encoding/json"
		"gopkg.in/yaml.v2"
		"github.com/olekukonko/tablewriter"
)

type Domain struct{
	Domain string	`yaml:"domain"`
	Find string	`yaml:"regex"`
}

type Settings struct {
		Ipv4 bool	`yaml:"ipv4"`
		Ipv6 bool	`yaml:"ipv6"`
		Domains	[]Domain	`yaml:"domains"`
}

type StatusDomain struct{
	Domain string
	Ipv4_http int
	Ipv6_http int
	Ipv4_https int
	Ipv6_https int
}

var (
	statusoutput = make(chan StatusDomain,5)
	outputwg = sync.WaitGroup{}
	domainwg = sync.WaitGroup{}
   config = &Settings{}
)
func main() {
	readConfig()
	outputwg.Add(1)
	go printTable()

   for _,domain := range config.Domains{
		domainwg.Add(1)
		fetchDomain(domain)
	}
	domainwg.Wait()
	close(statusoutput)
	outputwg.Wait()

}
func fetchDomain(domain Domain){
	n := StatusDomain{Domain:domain.Domain}
	find,err := regexp.Compile(domain.Find)
	if err !=nil{
		log.Println(err)
	}

	if config.Ipv4 {
		n.Ipv4_http = requestHttp("tcp4",domain.Domain,find)
		n.Ipv4_https = requestHttps("tcp4",domain.Domain,find)
	}
	if config.Ipv6 {
		n.Ipv6_http = requestHttp("tcp6",domain.Domain,find)
		n.Ipv6_https = requestHttps("tcp6",domain.Domain,find)
	}
	statusoutput <- n
	domainwg.Done()
}
func printTable(){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Domain", "IPv4", "IPv6", "SSL-IPv4", "SSL-IPv6"})
	for v := range statusoutput {
		table.Append([]string{v.Domain,printStatus(v.Ipv4_http),printStatus(v.Ipv6_http),printStatus(v.Ipv4_https),printStatus(v.Ipv6_https)})
	}
	table.Render()
	outputwg.Done()
}
func printStatus(s int) string{
	switch s {
	case 1:
		return "0";
	case 200:
		return "W";
	case 300:
		return "30";

	}
	return "E";
}

func requestHttp(tcp, domain string,find *regexp.Regexp) int {
	conn, err := net.Dial(tcp, domain+":http")
	handleError(err)
	_, err = conn.Write([]byte("GET / HTTP/1.0\nHOST: "+domain+"\r\n\r\n"))
	handleError(err)
	result, _ := ioutil.ReadAll(conn)
	return requestOutputInterpreter(result,find,false,false)
}
func requestHttps(tcp, domain string,find *regexp.Regexp) int {
	conf := tls.Config{}
	conn, err := tls.Dial(tcp, domain+":https",&conf)
	handleError(err)
	_, err = conn.Write([]byte("GET / HTTP/1.0\nHOST: "+domain+"\r\n\r\n"))
	handleError(err)
	result, _ := ioutil.ReadAll(conn)
	ssl := false;
	out := requestOutputInterpreter(result,find,true,ssl)
	return out;
}

func requestOutputInterpreter(result []byte,find *regexp.Regexp,ssl bool, sslwork bool) int{
	result_s := string(result)
	if strings.Contains(result_s,"HTTP/1.1 200") || strings.Contains(result_s,"HTTP/1.1 30"){
		if find.Match(result){
			return 1;
		}else{
			if strings.Contains(result_s,"HTTP/1.1 200"){
				return 200
			}else{
				return 300
			}
		}
	}
	return -1;
}
func readConfig() {
	file, _ := ioutil.ReadFile("/home/geno/.checkhttprc")
	err := yaml.Unmarshal(file,&config)
	if err != nil {
	  log.Panic("error:", err)
	}
}

func handleError(err error){
	if err != nil {
			log.Println(err)
			os.Exit(1)
	}
}
