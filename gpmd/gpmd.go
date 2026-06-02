package gpmd

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/ghostosproject/gpmd/config"
	"github.com/ghostosproject/gpmd/misc"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
)

type GPMD struct {
	Id        string
	Nodes     map[string]Node
	NodeConns map[net.Conn]string
	Networks  map[string]Network
	Peers     map[string]Peer
	Resources Resources
	P2PHost   host.Host
	Messages  map[string]string // keeps track of messages that require a response
}

type Node struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Peer    string // what peer the node is connected to
	Network string `json:"network"` // what network the node is connected to
	Active  bool
	Type    int
	GPMD    net.Conn // could be null if it is not a local node
}

type Connection struct {
	Type string `json:"type"`
	Node Node   `json:"node"`
}

type Network struct {
	Id        string
	Name      string
	Peers     []string
	Nodes     []string
	Type      int
	Private   bool // this means that a key is needed to connect to this network, will implement private options later
	Local     bool // this means that the network only exists on the current device
	Connected bool // gpmd currently supports this network
}

type Peer struct {
	Id        string
	Nodes     []string
	Addresses []multiaddr.Multiaddr
	Networks  []string // keeps the id for the network which is mapped to the Networks Section
	Stream    network.Stream
}

type Resources struct {
	Modules map[string]Module
}

type Module struct {
	Name     string                   `json:"name"`
	Versions map[string]ModuleVersion `json:"versions"`
}

type ModuleVersion struct {
	Version string `json:"version"`
	File    string `json:"file"`
	Hash    string `json:"hash"`
}

const discoverProtocolID = "/ghost/gpmd/0.1.0"

func (md GPMD) Run(cfg config.Config) {
	md.P2PHost.SetStreamHandler(discoverProtocolID, func(s network.Stream) {
		go handleGPMDStream(s, &md)
	})

	go runBackgroundServer(&md, cfg.Ports.Background)
	runGVMServer(&md, cfg.Ports.GVM)
}

func CreateGPMD(p2p bool, p2pPort, path string) GPMD {
	var node host.Host
	id := ""
	if p2p {
		node = setupP2P(p2pPort, path)
		id = node.ID().String()
	}

	netID := uuid.New().String()
	netw := Network{
		Id:        netID,
		Name:      "default",
		Peers:     []string{},
		Nodes:     []string{},
		Private:   true,
		Local:     true,
		Connected: true,
	}

	pmd := GPMD{
		Id:        id,
		Nodes:     map[string]Node{},
		Networks:  map[string]Network{netw.Id: netw},
		Peers:     map[string]Peer{},
		NodeConns: map[net.Conn]string{},
		P2PHost:   node,
	}

	return pmd
}

func setupP2P(port, path string) host.Host {
	var sk crypto.PrivKey
	if misc.FileExists(path) {
		bz, err := os.ReadFile(path)
		if err != nil {
			// Handle error
			panic(err)
		}

		sk, err = crypto.UnmarshalPrivateKey(bz)
		if err != nil {
			// Handle error
			panic(err)
		}
	} else {
		var r io.Reader
		var err error
		sk, _, err = crypto.GenerateEd25519Key(r)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}

		bz, err := crypto.MarshalPrivateKey(sk)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		err = os.WriteFile("./private.key", bz, 0600)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
	node, err := libp2p.New(
		// libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/3002"),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/udp/"+port+"/quic-v1/"),
		libp2p.Ping(false),
		libp2p.Identity(sk),
	)
	if err != nil {
		panic(err)
	}
	return node
}
