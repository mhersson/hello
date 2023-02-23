package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func checkError(w *http.ResponseWriter, msg string, err error) {
	if err != nil {
		fmt.Fprintf(*w, "%s: %s\n", msg, err)
	}
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

	fmt.Fprintln(w, "Hello from devspaces")
	fmt.Fprintf(w, "My workspace is %s!\n", hostname)
	fmt.Fprintln(w, "You can find me here:")

	for _, iface := range interfaces {
		addr, _ := iface.Addrs()
		if len(addr) > 0 && iface.Name != "lo" {
			fmt.Fprintf(w, "%s - %s\n", iface.Name, addr[0])
		}
	}
}

func main() {
	http.HandleFunc("/", simpleHandler)
	panic(http.ListenAndServe(":8080", nil))
}
