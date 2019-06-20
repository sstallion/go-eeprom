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

package eeprom_test

import (
	"testing"

	"github.com/sstallion/go-eeprom"
)

type dataCommand func(*eeprom.Device, uint16, []byte) error

func TestReset(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	d, err := eeprom.First()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	d.Reset()
}

func TestVerify(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	d, err := eeprom.First()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	tests := []struct {
		name string
		cmd  dataCommand
	}{
		{"WriteBytes", (*eeprom.Device).WriteBytes},
		{"WritePages", (*eeprom.Device).WritePages},
	}
	for _, test := range tests {
		t.Log("verifying", test.name)
		if err := d.Erase(); err != nil {
			t.Fatal(err)
		}
		wbuf := make([]byte, eeprom.MaxBytes)
		for i := range wbuf {
			wbuf[i] = byte(i)
		}
		if err := test.cmd(d, 0, wbuf); err != nil {
			t.Fatal(err)
		}
		rbuf := make([]byte, eeprom.MaxBytes)
		if err := d.Read(0, rbuf); err != nil {
			t.Fatal(err)
		}
		for i, b := range wbuf {
			if rbuf[i] != b {
				t.Fatalf("rbuf[%d]: expected %#x; got %#x", i, b, rbuf[i])
			}
		}
	}
}

func benchmarkDataCommand(b *testing.B, cmd dataCommand, n int) {
	d, err := eeprom.First()
	if err != nil {
		b.Fatal(err)
	}
	defer d.Close()

	data := make([]byte, n)
	if err := d.Erase(); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := cmd(d, 0, data); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRead8K(b *testing.B) {
	benchmarkDataCommand(b, (*eeprom.Device).Read, 8*1024)
}

func BenchmarkWriteBytes8K(b *testing.B) {
	benchmarkDataCommand(b, (*eeprom.Device).WriteBytes, 8*1024)
}

func BenchmarkWritePages8K(b *testing.B) {
	benchmarkDataCommand(b, (*eeprom.Device).WritePages, 8*1024)
}

func BenchmarkErase(b *testing.B) {
	d, err := eeprom.First()
	if err != nil {
		b.Fatal(err)
	}
	defer d.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := d.Erase(); err != nil {
			b.Fatal(err)
		}
	}
}
