package main

import (
	"flag"
	"log"

	"github.com/moemoe89/file-storage/pkg/rpcseqdiagram"
)

func main() {
	var rpc string

	flag.StringVar(&rpc, "RPC", "", "RPC, e.g. GetData,GetList")
	flag.Parse()

	rpcSeqDiagram, err := rpcseqdiagram.New()
	if err != nil {
		log.Fatal(err)
	}

	err = rpcSeqDiagram.Generate(rpc)
	if err != nil {
		log.Fatal(err)
	}
}
