# eeprom

    import "github.com/sstallion/go-eeprom"

Package eeprom provides a low-level interface to USB EEPROM programmers that
conform to http://github.com/sstallion/usb-eeprom/wiki/Protocol. Due to the
chip-agnostic nature of the protocol, constraints such as capacity and alignment
must be enforced by the caller.

## Usage

```go
const (
        // MaxBytes is the maximum amount of addressable data.
        MaxBytes = 1 << 16
)
```

#### func  Walk

```go
func Walk(fn func(*Device) error) error
```
Walk calls the specified function for each supported device attached to the
host. To ensure proper reference counting, Open must be called within the
context of a Walk.

#### type Device

```go
type Device struct {
}
```

Device represents an attached USB EEPROM programmer.

#### func  First

```go
func First() (*Device, error)
```
First returns the first supported device attached to the host. Unlike Walk, the
returned Device is opened automatically. This function exists primarily for
testing.

#### func (*Device) Close

```go
func (d *Device) Close() error
```
Close releases the interface and closes the device. A device may not be opened
again after calling this method. Returned errors may be safely ignored.

#### func (*Device) Erase

```go
func (d *Device) Erase() error
```
Erase issues a Chip Erase.

#### func (*Device) ID

```go
func (d *Device) ID() string
```
ID returns a string suitable for uniquely identifying the device.

#### func (*Device) Open

```go
func (d *Device) Open() error
```
Open opens an attached device and claims the interface. To ensure proper
reference counting, Open must be called within the context of a Walk.

#### func (*Device) Read

```go
func (d *Device) Read(start uint16, data []byte) error
```
Read reads into the given slice at the supplied starting address.

#### func (*Device) Reset

```go
func (d *Device) Reset() error
```
Reset issues a device reset. This method may be called after a failed transfer
to reset the interface. Returned errors may be safely ignored.

#### func (*Device) SetPageSize

```go
func (d *Device) SetPageSize(pagesize int)
```
SetPageSize sets the number of bytes written per page by WritePages. By default,
the maximum packet size supported by the endpoint is used.

#### func (*Device) WriteBytes

```go
func (d *Device) WriteBytes(start uint16, data []byte) error
```
WriteBytes writes the given slice starting at the supplied starting address.

#### func (*Device) WritePages

```go
func (d *Device) WritePages(start uint16, data []byte) error
```
WritePages writes the given slice starting at the supplied starting address.
