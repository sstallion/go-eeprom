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

import "io/ioutil"

var writeStart, writeCount, writePagesize int

func init() {
	cmd := &command{
		name: "write",
		exec: write,
		help: `usage: eeprom write [-start addr] [-count n] [-pagesize n] file

The write command writes the contents of the specified file to the device.

The flags are:

    -start addr
		starting address; by default this is 0.
    -count n
		number of bytes to write; by default this is the length of
		the file.
    -pagesize n
		page size to use when writing; by default page writes are
		disabled for compatibility.
`,
	}
	cmd.flag.IntVar(&writeStart, "start", 0, "")
	cmd.flag.IntVar(&writeCount, "count", 0, "")
	cmd.flag.IntVar(&writePagesize, "pagesize", 0, "")
	addCommand(cmd)
}

func write(args ...string) error {
	if len(args) < 1 {
		return errUsage
	}
	data, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	d, err := openDevice()
	if err != nil {
		return err
	}
	defer d.Close()

	if writeCount == 0 || writeCount > len(data) {
		writeCount = len(data)
	}
	if writePagesize > 0 {
		d.SetPageSize(writePagesize)
		err = d.WritePages(uint16(writeStart), data[:writeCount])
	} else {
		err = d.WriteBytes(uint16(writeStart), data[:writeCount])
	}
	if err != nil {
		d.Reset()
	}
	return err
}
