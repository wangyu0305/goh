/* Autogenerated by Thrift Compiler (0.9.0)
 *
 * DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
 */
package main

import (
	"flag"
	"fmt"
	"librarytest"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"thrift"
)

func Usage() {
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "Functions:\n")
	fmt.Fprint(os.Stderr, "  methodThatThrowsAnException() (xwamap *ExceptionWithAMap, err error)\n")
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(0)
}

func main() {
	flag.Usage = Usage
	var host string
	var port int
	var protocol string
	var urlString string
	var framed bool
	var useHttp bool
	var help bool
	var parsedUrl url.URL
	var trans thrift.TTransport
	flag.Usage = Usage
	flag.StringVar(&host, "h", "localhost", "Specify host and port")
	flag.IntVar(&port, "p", 9090, "Specify port")
	flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
	flag.StringVar(&urlString, "u", "", "Specify the url")
	flag.BoolVar(&framed, "framed", false, "Use framed transport")
	flag.BoolVar(&useHttp, "http", false, "Use http")
	flag.BoolVar(&help, "help", false, "See usage string")
	flag.Parse()
	if help || flag.NArg() == 0 {
		flag.Usage()
	}

	if len(urlString) > 0 {
		parsedUrl, err := url.Parse(urlString)
		if err != nil {
			fmt.Fprint(os.Stderr, "Error parsing URL: ", err.Error(), "\n")
			flag.Usage()
		}
		host = parsedUrl.Host
		useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http"
	} else if useHttp {
		_, err := url.Parse(fmt.Sprint("http://", host, ":", port))
		if err != nil {
			fmt.Fprint(os.Stderr, "Error parsing URL: ", err.Error(), "\n")
			flag.Usage()
		}
	}

	cmd := flag.Arg(0)
	var err error
	if useHttp {
		trans, err = thrift.NewTHttpClient(parsedUrl.String())
	} else {
		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprint(host, ":", port))
		if err != nil {
			fmt.Fprint(os.Stderr, "Error resolving address", err.Error())
			os.Exit(1)
		}
		trans, err = thrift.NewTNonblockingSocketAddr(addr)
		if framed {
			trans = thrift.NewTFramedTransport(trans)
		}
	}
	if err != nil {
		fmt.Fprint(os.Stderr, "Error creating transport", err.Error())
		os.Exit(1)
	}
	defer trans.Close()
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
		break
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
		break
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
		break
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
		break
	default:
		fmt.Fprint(os.Stderr, "Invalid protocol specified: ", protocol, "\n")
		Usage()
		os.Exit(1)
	}
	client := librarytest.NewServiceForExceptionWithAMapClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		fmt.Fprint(os.Stderr, "Error opening socket to ", host, ":", port, " ", err.Error())
		os.Exit(1)
	}

	switch cmd {
	case "methodThatThrowsAnException":
		if flag.NArg()-1 != 0 {
			fmt.Fprint(os.Stderr, "MethodThatThrowsAnException requires 0 args\n")
			flag.Usage()
		}
		fmt.Print(client.MethodThatThrowsAnException())
		fmt.Print("\n")
		break
	case "":
		Usage()
		break
	default:
		fmt.Fprint(os.Stderr, "Invalid function ", cmd, "\n")
	}
}
