package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/qiniu/log"

	"github.com/codeskyblue/readline" // a fork version in order to support windows
	"github.com/gobuild/goyaml"
	"github.com/kballard/go-shellquote"
)

var cfg struct {
	Proxies map[string]string `goyaml:"proxies"`
}

var (
	httpTimeout = 3000 * time.Millisecond
)

func init() {
	// Load config file
	cfgdata, err := ioutil.ReadFile("_config.yml")
	if err != nil {
		log.Fatal(err)
	}
	if err = goyaml.Unmarshal(cfgdata, &cfg); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Proxy engine started...")
	proxies := make([]*Proxy, 0)
	for local, remote := range cfg.Proxies {
		// eg: NewProxy("localhost:80", "10.0.0.1:4000")
		px, err := NewProxy(local, remote)
		if err != nil {
			log.Fatal(err)
		}
		proxies = append(proxies, px)
		fmt.Printf("\t%v --> %v\n", px.laddr, px.raddr)
		go px.ListenAndServe()
	}

	fmt.Printf("\nWelcome to reverse proxy console\n\n")
	prefix := "ReverseProxy $ "
	for {
		l, err := readline.String(prefix)
		if err != nil {
			if err != io.EOF {
				log.Fatal("error: ", err)
			} else {
				println("")
				break
			}
		}

		args, err := shellquote.Split(l)
		if err != nil {
			log.Println("shellquote", err)
			continue
		}
		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "h", "help":
			fmt.Println(genHelp("Usage:", map[string]string{
				"help/h":  "Show help information",
				"print/p": "Show netstat info",
				"exit":    "Exit program",
			}))
		case "p", "print":
			for _, px := range proxies {
				fmt.Printf("[%v]: %d bytes send, %d bytes received\n", px.laddr, px.sentBytes, px.receivedBytes)
			}
		case "exit":
			os.Exit(0)
		default:
			fmt.Printf("- %s: command not found, type help for more information\n", l)
			continue
		}
		readline.AddHistory(l)
	}
}
