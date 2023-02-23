package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

func checkError(w *http.ResponseWriter, msg string, err error) {
	if err != nil {
		fmt.Fprintf(*w, "%s: %s\n", msg, err)
	}
}

func getQuote() (string, error) {
	client := http.Client{
		Timeout: 250 * time.Millisecond, //nolint: go-lint
	}

	req, _ := http.NewRequest("GET", "http://hello-quotes", nil)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	hostname, err := os.Hostname()
	checkError(&w, "Failed to get hostname", err)

	interfaces, err := net.Interfaces()
	checkError(&w, "Failed to get interfaces", err)

	fmt.Fprintf(w, "Hello!\nMy name is %s!\n", hostname)
	fmt.Fprintln(w, "You can find me here:")

	for _, iface := range interfaces {
		addr, _ := iface.Addrs()
		if len(addr) > 0 && iface.Name != "lo" {
			fmt.Fprintf(w, "%s - %s\n", iface.Name, addr[0])
		}
	}

	quote, err := getQuote()
	if err == nil {
		fmt.Fprint(w, quote)
	}
}

func main() {
	http.HandleFunc("/", simpleHandler)
	panic(http.ListenAndServe(":8080", nil))
}
