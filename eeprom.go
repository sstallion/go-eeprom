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

// Package eeprom provides an idiomatic interface to USB EEPROM programmers
// that conform to http://github.com/sstallion/usb-eeprom/wiki/Protocol. Due
// to the chip-agnostic nature of the protocol, constraints such as capacity
// and alignment must be enforced by the caller.
package eeprom

/*
#cgo LDFLAGS: -lusb-1.0
#include <libusb-1.0/libusb.h>
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

const (
	// MaxBytes is the maximum amount of addressable data.
	MaxBytes = 1 << 16
)

const (
	idVendor     = 0x04d8 // Microchip Technology, Inc.
	idProduct    = 0xf4cd // 28Cxxx EEPROM Programmer
	interfaceNum = 0
	endpointNum  = 1
	endpointIN   = endpointNum | C.LIBUSB_ENDPOINT_IN
	endpointOUT  = endpointNum | C.LIBUSB_ENDPOINT_OUT
)

type libusbError struct {
	code C.int
}

func (e *libusbError) Error() string {
	return fmt.Sprintf("%s (%s)",
		C.GoString(C.libusb_strerror(C.enum_libusb_error(e.code))),
		C.GoString(C.libusb_error_name(e.code)))
}

// Device represents an attached USB EEPROM programmer.
type Device struct {
	dev      *C.libusb_device
	handle   *C.libusb_device_handle
	pagesize int
}

// ID returns a string suitable for uniquely identifying the device.
func (d *Device) ID() string {
	return fmt.Sprintf("%d:%d",
		C.libusb_get_bus_number(d.dev),
		C.libusb_get_device_address(d.dev))
}

// SetPageSize sets the number of bytes written per page by WritePages. By
// default, the maximum packet size supported by the endpoint is used.
func (d *Device) SetPageSize(pagesize int) { d.pagesize = pagesize }

// Open opens an attached device and claims the interface. To ensure proper
// reference counting, Open must be called within the context of a Walk.
func (d *Device) Open() error {
	if err := C.libusb_open(d.dev, &d.handle); err != C.LIBUSB_SUCCESS {
		return &libusbError{err}
	}
	if err := C.libusb_claim_interface(d.handle, interfaceNum); err != C.LIBUSB_SUCCESS {
		C.libusb_close(d.handle)
		return &libusbError{err}
	}
	return nil
}

// Close releases the interface and closes the device. A device may not be
// opened again after calling this method. Returned errors may be safely
// ignored.
func (d *Device) Close() error {
	defer C.libusb_close(d.handle)

	if err := C.libusb_release_interface(d.handle, interfaceNum); err != C.LIBUSB_SUCCESS {
		return &libusbError{err}
	}
	return nil
}

// Reset issues a device reset. This method may be called after a failed
// transfer to reset the interface. Returned errors may be safely ignored.
func (d *Device) Reset() error {
	defer time.Sleep(500 * time.Millisecond) // wait for device to settle

	if err := C.libusb_reset_device(d.handle); err != C.LIBUSB_SUCCESS {
		return &libusbError{err}
	}
	return nil
}

// Read reads into the given slice at the supplied starting address.
func (d *Device) Read(start uint16, data []byte) error {
	var b bytes.Buffer
	var n = uint16(len(data))

	b.WriteByte('R')
	binary.Write(&b, binary.LittleEndian, start)
	binary.Write(&b, binary.LittleEndian, n-1)

	if err := d.validate(start, data); err != nil {
		return err
	}
	if err := d.transfer(endpointOUT, b.Bytes()); err != nil {
		return err
	}
	if err := d.transfer(endpointIN, data); err != nil {
		return err
	}
	return d.verify(start + n)
}

// WriteBytes writes the given slice starting at the supplied starting address.
func (d *Device) WriteBytes(start uint16, data []byte) error {
	var b bytes.Buffer
	var n = uint16(len(data))

	b.WriteByte('W')
	binary.Write(&b, binary.LittleEndian, start)
	binary.Write(&b, binary.LittleEndian, n-1)

	if err := d.validate(start, data); err != nil {
		return err
	}
	if err := d.transfer(endpointOUT, b.Bytes()); err != nil {
		return err
	}
	if err := d.transfer(endpointOUT, data); err != nil {
		return err
	}
	return d.verify(start + n)
}

// WritePages writes the given slice starting at the supplied starting address.
func (d *Device) WritePages(start uint16, data []byte) error {
	var b bytes.Buffer
	var n = uint16(len(data))

	b.WriteByte('P')
	binary.Write(&b, binary.LittleEndian, start)
	binary.Write(&b, binary.LittleEndian, n-1)

	if err := d.validate(start, data); err != nil {
		return err
	}
	if err := d.transfer(endpointOUT, b.Bytes()); err != nil {
		return err
	}
	if err := d.transferN(endpointOUT, data, d.pagesize); err != nil {
		return err
	}
	return d.verify(start + n)
}

// Erase issues a Chip Erase.
func (d *Device) Erase() error {
	var b bytes.Buffer

	b.WriteByte('Z')

	if err := d.transfer(endpointOUT, b.Bytes()); err != nil {
		return err
	}
	return d.verify(0)
}

func (d *Device) validate(start uint16, data []byte) error {
	if len(data) == 0 {
		return errors.New("no data")
	}
	if int(start)+len(data) > MaxBytes {
		return errors.New("too much data")
	}
	return nil
}

func (d *Device) transfer(endpoint uint8, data []byte) error {
	return d.transferN(endpoint, data, 0)
}

func (d *Device) transferN(endpoint uint8, data []byte, n int) error {
	if m := int(C.libusb_get_max_packet_size(d.dev, C.uchar(endpoint))); n == 0 {
		n = m
	} else if n > m {
		return errors.New("invalid packet size")
	}

	for len, off := len(data), 0; len > 0; {
		var transferred int

		if n > len {
			n = len
		}
		if err := C.libusb_bulk_transfer(d.handle, C.uchar(endpoint), (*C.uchar)(&data[off]), C.int(n),
			(*C.int)(unsafe.Pointer(&transferred)), 2500); err != C.LIBUSB_SUCCESS {
			return &libusbError{err}
		}
		len -= transferred
		off += transferred
	}
	return nil
}

func (d *Device) verify(expected uint16) error {
	var status uint16
	var data = []byte{0xff, 0xff}

	if err := d.transfer(endpointIN, data); err != nil {
		return err
	}
	binary.Read(bytes.NewReader(data), binary.LittleEndian, &status)
	if status != expected {
		return fmt.Errorf("expected status %#x; got %#x", expected, status)
	}
	return nil
}

var context *C.libusb_context

func init() {
	if err := C.libusb_init(&context); err != C.LIBUSB_SUCCESS {
		panic(&libusbError{err})
	}
}

/*
func fini() {
	C.libusb_exit(context)
}
*/

// First returns the first supported device attached to the host. Unlike Walk,
// the returned Device is opened automatically. This function exists primarily
// for testing.
func First() (*Device, error) {
	handle := C.libusb_open_device_with_vid_pid(context, idVendor, idProduct)
	if handle == nil {
		return nil, errors.New("no devices found")
	}
	if err := C.libusb_claim_interface(handle, interfaceNum); err != C.LIBUSB_SUCCESS {
		C.libusb_close(handle)
		return nil, &libusbError{err}
	}
	return &Device{
		dev:    C.libusb_get_device(handle),
		handle: handle,
	}, nil
}

// Walk calls the specified function for each supported device attached to the
// host. To ensure proper reference counting, Open must be called within the
// context of a Walk.
func Walk(fn func(*Device) error) error {
	var list **C.libusb_device
	var found int

	n := C.libusb_get_device_list(context, &list)
	if n < C.LIBUSB_SUCCESS {
		return &libusbError{C.int(n)}
	}
	defer C.libusb_free_device_list(list, 1)

	h := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(list)),
		Len:  int(n),
		Cap:  int(n),
	}
	for _, dev := range *(*[]*C.libusb_device)(unsafe.Pointer(&h)) {
		var desc C.struct_libusb_device_descriptor

		if err := C.libusb_get_device_descriptor(dev, &desc); err != C.LIBUSB_SUCCESS {
			return &libusbError{err}
		}
		if desc.idVendor == idVendor && desc.idProduct == idProduct {
			if err := fn(&Device{dev: dev}); err != nil {
				return err
			}
			found++
		}
	}
	if found == 0 {
		return errors.New("no devices found")
	}
	return nil
}
