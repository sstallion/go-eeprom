# eeprom

eeprom is a tool that manages USB EEPROM programmers that conform to
http://github.com/sstallion/usb-eeprom/wiki/Protocol.

Usage:

    eeprom [-id device] command [arguments]

The flags are:

        -id device
                identifies device to use; by default the first supported
                device is selected.

The commands are:

    dump        dump contents of device
    erase       erase contents of device
    reset       hard reset device
    verify      verify contents of device
    write       write file to device

Use "eeprom help [command]" for more information about a command.
