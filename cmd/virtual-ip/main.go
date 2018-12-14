package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/darxkies/virtual-ip/pkg"
	"github.com/darxkies/virtual-ip/version"
	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Printf("Version: %s\n", version.Version)

	id := flag.String("id", "vip", "ID of this node")
	bind := flag.String("bind", "0.0.0.0", "RAFT bind addreess")
	virtualIP := flag.String("virtual-ip", "192.168.0.25", "Virtual/Floating IP")
	timeout := flag.Int("timeout", 10, "Shell command timeout")
	peersList := flag.String("peers", "", "Peers as a comma separated list of peer-id=peer-address:peer-port including the id and the bind of this instance")
	_interface := flag.String("interface", "enp3s0:1", "Network interface")
	flag.Parse()

	commandRunner := pkg.NewShellCommandRunner(time.Duration(*timeout))

	peers := pkg.Peers{}

	if len(*peersList) > 0 {
		for _, peer := range strings.Split(*peersList, ",") {
			peerTokens := strings.Split(peer, "=")

			if len(peerTokens) != 2 {
				log.WithFields(log.Fields{"peer": peer}).Error("Peers malformated")

				os.Exit(-1)
			}

			peers[peerTokens[0]] = peerTokens[1]
		}
	}

	logger := pkg.Logger{}

	vipManager := pkg.NewVIPManager(*id, *bind, *virtualIP, peers, logger, *_interface, commandRunner)
	if error := vipManager.Start(); error != nil {
		log.WithFields(log.Fields{"error": error}).Error("Start failed")

		os.Exit(-1)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan

	vipManager.Stop()
}
