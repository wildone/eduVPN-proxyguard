package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"codeberg.org/eduVPN/proxyguard"
)

func main() {
	fwmark := flag.Int("fwmark", -1, "[Client + Linux only] The fwmark/SO_MARK to use on the TCP client socket. -1 is disable.")
	listen := flag.String("listen", "", "The IP:PORT to listen for UDP (client) or TCP (server) traffic.")
	to := flag.String("to", "", "The IP:PORT to which to send the converted traffic to. For the client you would specify a TCP server, for the server the WireGuard endpoint. The WireGuard endpoint for the client is automatically inferred from the first UDP packet")
	ic := flag.Bool("client", false, "Indicates that this should function as a client, proxying UDP packets to a TCP server")
	is := flag.Bool("server", false, "Indicates that this should function as a server, proxying TCP packets to UDP")
	flag.Parse()
	// Both a client and server flag are not supplied
	if !*ic && !*is {
		fmt.Fprintln(os.Stderr, "Invalid invocation error: Please supply the --client or --server flag")
		flag.PrintDefaults()
		os.Exit(1)
	}
	// Both a client and server flag are supplied
	if *ic && *is {
		fmt.Fprintln(os.Stderr, "Invalid invocation error: Please supply only one --client or --server flag")
		flag.PrintDefaults()
		os.Exit(1)
	}
	// listen and to flags are also mandatory
	if *listen == "" {
		fmt.Fprintln(os.Stderr, "Invalid invaction error: Please supply the --listen flag")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *to == "" {
		fmt.Fprintln(os.Stderr, "Invalid invocation error: Please supply the --to flag")
		flag.PrintDefaults()
		os.Exit(1)
	}
	// Warn that the server does not use fwmark
	if *is && *fwmark != -1 {
		fmt.Fprintln(os.Stderr, "Invalid invocation warning: The --fwmark flag is a NO-OP for the server")
		*fwmark = -1
	}
	// fwmark flag is given for the client but we are not linux
	if *ic && *fwmark != -1 && runtime.GOOS != "linux" {
		fmt.Fprintln(os.Stderr, "Invalid invocation warning: The --fwmark flag is a NO-OP when you're not using Linux. We will ignore it...")
		*fwmark = -1
	}
	// We are a client
	if *ic {
		err := proxyguard.Client(context.Background(), *listen, *to, *fwmark)
		if err != nil {
			log.Fatalf("error occurred when setting up a client: %v", err)
		}
		return
	}
	// We are a server
	err := proxyguard.Server(context.Background(), *listen, *to)
	if err != nil {
		log.Fatalf("error occurred when setting up server side: %v", err)
	}
}
