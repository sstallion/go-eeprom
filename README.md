# USB EEPROM Programming

[![](https://travis-ci.org/sstallion/go-eeprom.svg?branch=master)][1]
[![](https://godoc.org/github.com/sstallion/go-eeprom?status.svg)][2]
[![](https://goreportcard.com/badge/github.com/sstallion/go-eeprom)][3]
[![](https://img.shields.io/github/license/sstallion/go-eeprom.svg)][LICENSE]

Package `eeprom` provides an idiomatic interface to USB EEPROM programmers that
conform to http://github.com/sstallion/usb-eeprom/wiki/Protocol. Due to the
chip-agnostic nature of the protocol, constraints such as capacity and alignment
must be enforced by the caller.

## Documentation

Up-to-date documentation can be found on [GoDoc][2], or by issuing the `go doc`
command after installing the package:

    $ go doc -all github.com/sstallion/go-eeprom

## Installation

> **Note**: [libusb][4], is required to use this package and should be installed
> using the system package manager (eg. `libusb-dev` on Debian-based
> distributions via `apt-get`).

Package `eeprom` may be installed via the `go get` command:

    $ go get github.com/sstallion/go-eeprom

### eeprom

A command named `eeprom` is provided, which manages USB EEPROM programmers.
`eeprom` may be installed by issuing:

    $ go get github.com/sstallion/go-eeprom/cmd/eeprom

Once installed, issue `eeprom help` to display usage.

## Contributing

Pull requests are welcome! If a problem is encountered using this package,
please file an issue on [GitHub][5].

## License

Source code in this repository is licensed under a Simplified BSD License. See
[LICENSE] for more details.

[1]: https://travis-ci.org/sstallion/go-eeprom
[2]: https://godoc.org/github.com/sstallion/go-eeprom
[3]: https://goreportcard.com/report/github.com/sstallion/go-eeprom
[4]: https://libusb.info/
[5]: https://github.com/sstallion/go-eeprom/issues/new

[LICENSE]: LICENSE
