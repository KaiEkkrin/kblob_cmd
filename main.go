/* A simple command line interface to komblobulate. */

package main

import (
	"flag"
	"fmt"
	"github.com/kaiekkrin/komblobulate"
	"io"
	"os"
)

type Params struct {
	Resist *string
	Cipher *string

	DataPieceSize    *int
	DataPieceCount   *int
	ParityPieceCount *int

	CipherChunkSize *int
	Password        *string
}

func (p *Params) GetRsParams() (int, int, int) {
	return *p.DataPieceSize, *p.DataPieceCount, *p.ParityPieceCount
}

func (p *Params) GetAeadChunkSize() int {
	return *p.CipherChunkSize
}

func (p *Params) GetAeadPassword() string {
	return *p.Password
}

func doEncode(inFile, outFile string, resist, cipher byte, params *Params) (err error) {
	var inf, outf *os.File
	inf, err = os.Open(inFile)
	if err != nil {
		return
	}
	defer func() {
		inf.Close()
	}()

	outf, err = os.Create(outFile)
	if err != nil {
		return
	}
	defer func() {
		outf.Close()
	}()

	var writer io.WriteCloser
	writer, err = komblobulate.NewWriter(outf, resist, cipher, params)
	if err != nil {
		return
	}
	defer func() {
		writer.Close()
	}()

	_, err = io.Copy(writer, inf)
	return
}

func doDecode(inFile, outFile string, params *Params) (err error) {
	var inf, outf *os.File
	inf, err = os.Open(inFile)
	if err != nil {
		return
	}
	defer func() {
		inf.Close()
	}()

	outf, err = os.Create(outFile)
	if err != nil {
		return
	}
	defer func() {
		outf.Close()
	}()

	var reader io.Reader
	reader, err = komblobulate.NewReader(inf, params)
	if err != nil {
		return
	}

	_, err = io.Copy(outf, reader)
	return
}

func main() {
	params := new(Params)

	inFile := flag.String("in", "", "Input file")
	outFile := flag.String("out", "", "Output file")

	// TODO More modes.  Verify, scrub?
	encode := flag.Bool("encode", false, "Encode a file")
	decode := flag.Bool("decode", false, "Decode a file")

	resist := flag.String("resist", "rs", "Resist type (rs or none)")
	cipher := flag.String("cipher", "aead", "Cipher type (aead or none)")

	params.DataPieceSize = flag.Int("dps", 508, "Size of each encoded data piece")
	params.DataPieceCount = flag.Int("dpc", 8, "Number of data pieces per chunk")
	params.ParityPieceCount = flag.Int("ppc", 1, "Number of parity pieces per chunk")

	params.CipherChunkSize = flag.Int("ccs", 256*1024, "Chunk size of enciphered data")
	// TODO Make it do a password input at the command
	// line, optionally/instead
	params.Password = flag.String("password", "", "Password for enciphered data")

	flag.Parse()

	if len(*inFile) == 0 || len(*outFile) == 0 || *encode == *decode || (*resist != "rs" && *resist != "none") && (*cipher != "aead" && *cipher != "none") {
		flag.PrintDefaults()
		os.Exit(2)
	}

	var err error
	if *encode {
		r := komblobulate.ResistType_None
		if *resist == "rs" {
			r = komblobulate.ResistType_Rs
		}

		c := komblobulate.CipherType_None
		if *cipher == "aead" {
			c = komblobulate.CipherType_Aead
		}

		err = doEncode(*inFile, *outFile, r, c, params)
	} else {
		err = doDecode(*inFile, *outFile, params)
	}

	exitCode := 0
	if err != nil {
		fmt.Printf(err.Error())
		exitCode = 1
	}

	os.Exit(exitCode)
}
