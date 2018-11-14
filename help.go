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

	longHelpLoad = `Load binary file to memory location

Given a filename, load it as binary data to the specified memory range.

Examples:
		load myOS.bin to $00:FFFF
		load myStuff.bin to bank 2`

	longHelpInfo = `Query built-in manual pages

The manual pages include information on SAN mnemonics and opcodes as well as other
topics. A better name for this command would have been "help", but it was already
taken by the shell software.

Information on the opcodes includes their SAN and traditional mnemonics, size in
bytes, execution time in cycles formal name, MPU models they are available on,
and notes.  

Options:
		<SAN MNEMONIC>  - Returns information on the instruction
		<OPCODE>        - Returns information on the instruction
		all             - Gives a list of all instructions
		json <FILENAME> - Saves a JSON file with the information to file

Why the "json" option? The information about the opcodes can be used for other
projects, and sooner or later you'll want to have it.

Examples:
		info jmp.xi		; information on the opcode
		info 0x00		; information on brk instruction
		info json opcodes.json   ; saves databank to file`

	longHelpShow = `Show information on the system

Show produces information on larger elements of the system.

Options:

		breakpoints	- Lists defined breakpoints
		config		- Current configuration of machine
		memory		- List of chunks in memory
		specials	- Special addresses
		vectors		- Boot and interrupt vectors

For information on the computer angua is running on, use "status system". For
information on the CPU, use "status". To see information on the stack, use
"dump stack". For the contents of the Direct Page", use "dump direct". For
a hex dump of a given range, use "dump <ADDRESS RANGE>".`
)
