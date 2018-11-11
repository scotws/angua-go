// Long Help messages for the Angua Command Line Interface (CLI)
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 11. Nov 2018
// This version: 11. Nov 2018

package main

const (
	longHelpBeep = `Produce a beeping noise

Seriously, that is all this command does. It goes BEEP.

Example:
		beep		; duh, right?`

	longHelpDump = `Create a hex dump of the address range

Supply an address range in the usual format "<ADDRESS> to <ADDRESS>" or as a
bank with "bank <BANK>". Output is formatted in hex with the printable ASCII
characters in a separate column. Unicode is currently not supported.

Examples:
                dump $00:1000 to $00:1FFF
                dump 0x1000 to 0xFFFF     ; defaults to bank 0
                dump bank 2
                dump bank 0x1F`

	longHelpFill = `Fill a memory range with a byte

Fill can be used to set a memory range to a byte, for instance to store 0xEA
(the nop instruction). The format takes the usual address range or bank number.
Used with a zero byte, this can be used to erase large address ranges.

Examples:
		fill $00:2000 to $00:2FFF with 0xea
		fill bank 2 with 00`

	longHelpInit = `Initialize a new machine

Init loads a configuration file from the configs folder, sets up some background
stuff and starts the Switcher Daemon which is responsible for switching the
65816 from emulated to native modes and vice versa. If no configuration file is
provided, it uses the "default.cfg" file. It is an error to try to initialize an
already initialized machine.

Examples:
		init			; uses configs/default.cfg
		init my65816.cfg	; runs file configs/myconfig65816.cfg`
)
