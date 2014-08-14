package main

import (
	"flag"
	"fmt"
	"github.com/vbatts/go-bt/bencode"
	"github.com/vbatts/go-bt/torrent"
	"os"
)

var (
	flOutput = flag.String("o", "", "output the re-encoded torrent to file at this path")
)

func main() {
	flag.Parse()

	if len(*flOutput) > 0 && flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "-o and multiple input files can not be used together")
		os.Exit(1)
	}

	for _, arg := range flag.Args() {
		fh, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			continue
		}

		data, err := bencode.Decode(fh)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			fh.Close()
			continue
		}
		fh.Close()

		tf, err := torrent.DecocdeTorrentData(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			continue
		}
		fmt.Printf("Loaded: %s (%d files)\n", tf.Info.Name, len(tf.Info.Files))

		if len(*flOutput) > 0 {
			fhOutput, err := os.Create(*flOutput)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				continue
			}
			err = bencode.Marshal(fhOutput, *tf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			fmt.Printf("wrote: %s\n", fhOutput.Name())
			fhOutput.Close()
		}
	}
}
