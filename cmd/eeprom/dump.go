// Copyright (C) 2014 Steven Stallion <sstallion@gmail.com>
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

package main

import (
	"encoding/hex"
	"io"
	"os"

	"github.com/sstallion/go-eeprom"
)

var dumpStart, dumpCount int

func init() {
	cmd := &command{
		name: "dump",
		exec: dump,
		help: `usage: eeprom dump [-start addr] [-count n] [file]

The dump command reads data from the device and emits a hexdump to standard
output. If specified, dump will write the contents of the device to the given
file, creating it if necessary.

The flags are:

    -start addr
		starting address; by default this is 0.
    -count n
		number of bytes to read; by default this is the maximum
		number of bytes supported by the device.
`,
	}
	cmd.flag.IntVar(&dumpStart, "start", 0, "")
	cmd.flag.IntVar(&dumpCount, "count", 0, "")
	addCommand(cmd)
}

func dump(args ...string) error {
	var w io.WriteCloser

	if len(args) == 0 {
		w = hex.Dumper(os.Stdout)
	} else {
		var err error

		w, err = os.Create(args[0])
		if err != nil {
			return err
		}
	}
	defer w.Close()

	d, err := openDevice()
	if err != nil {
		return err
	}
	defer d.Close()

	if dumpCount == 0 {
		dumpCount = eeprom.MaxBytes - dumpStart
	}
	data := make([]byte, dumpCount)
	if err := d.Read(uint16(dumpStart), data); err != nil {
		d.Reset()
		return err
	}
	_, err = w.Write(data)
	return err
}
