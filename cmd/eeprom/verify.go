// Copyright (c) 2014, Steven Stallion
// All rights reserved.
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
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
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
	"fmt"
	"io/ioutil"
)

var verifyCount, verifyStart int

func init() {
	cmd := &command{
		name: "verify",
		exec: verify,
		help: `usage: eeprom verify [-start addr] [-count n] file

The verify command reads data from the device and performs a bytewise
comparison against the specified file.

The flags are:

    -start addr
		starting address; by default this is 0.
    -count n
		number of bytes to verify; by default this is the length of
		the file.
`,
	}
	cmd.flag.IntVar(&verifyStart, "start", 0, "")
	cmd.flag.IntVar(&verifyCount, "count", 0, "")
	addCommand(cmd)
}

func verify(args ...string) error {
	if len(args) < 1 {
		return errUsage
	}
	file, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	d, err := openDevice()
	if err != nil {
		return err
	}
	defer d.Close()

	if verifyCount == 0 || verifyCount > len(file) {
		verifyCount = len(file)
	}
	data := make([]byte, verifyCount)
	if err := d.Read(uint16(verifyStart), data); err != nil {
		d.Reset()
		return err
	}
	for i, b := range file {
		if data[i] != b {
			return fmt.Errorf("%s:%d: expected %#x; got %#x", args[0], i, b, data[i])
		}
	}
	return nil
}
