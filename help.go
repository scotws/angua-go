// Long Help messages for the Angua Command Line Interface (CLI)
// Scot W. Stevenson <scot.stevenson@gmail.com>
// First version: 11. Nov 2018
// This version: 01. Jan 2019

package main

const (
	longHelpAbort string = `Trigger the Abort Vector

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpBeep string = `Produce a beeping noise

Seriously, that is all this command does. It goes BEEP.

Example:
		beep		; duh, right?`

	longHelpBoot string = `Boot the machine.

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpDestroy string = `Destory the current machine.

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpDisasm string = `Disassemble a memory range.

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpDump string = `Create a hex dump of the address range

Supply an address range in the usual format "<ADDRESS> to <ADDRESS>" or as a
bank with "bank <BANK>"; or dump the stack or direct page. Output is formatted
in hex with the printable ASCII characters in a separate column. Unicode is
currently not supported.

Examples:
                dump $00:1000 to $00:1FFF
                dump 0x1000 to 0xFFFF     ; defaults to bank 0
                dump bank 2
                dump bank 0x1F
		dump stack
		dump dp		          ; synonyms "direct" and "directpage"
		
Because of the size of the stack, only the first elements are shown.`

	longHelpEcho string = `Print a character string.

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpFill string = `Fill a memory range with a byte

Fill can be used to set a memory range to a byte, for instance to store 0xEA
(the nop instruction). The format takes the usual address range or bank number.
Used with a zero byte, this can be used to erase large address ranges.

Examples:
		fill $00:2000 to $00:2FFF with 0xea
		fill bank 2 with 00`

	longHelpHalt string = `Halt a running machine.

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpInit string = `Initialize a new machine.

Init loads a configuration file from the configs folder, sets up some background
stuff and starts the Switcher Daemon which is responsible for switching the
65816 from emulated to native modes and vice versa. If no configuration file is
provided, it uses the "default.cfg" file. It is an error to try to initialize an
already initialized machine.

Examples:
		init			; uses configs/default.cfg
		init my65816.cfg	; runs file configs/myconfig65816.cfg`

	longHelpIRQ string = `Trigger a maskable Interrupt (IRQ).

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpLoad string = `Load binary file to memory location

Given a filename, load it as binary data to the specified memory range.

Examples:
		load myOS.bin to $00:FFFF
		load myStuff.bin to bank 2`

	longHelpInfo string = `Query built-in manual pages

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

Examples:
		info jmp.xi		; information on the opcode
		info 0x00		; information on brk instruction`

	longHelpMemory string = `Define a memory chunk.

Memory in Angua is organized as "chunks", which are continuous regions that can
be either read-only ("ROM") or read and write ("RAM"). They are defined either
by passing a memory range or a block number.

Example:
	memory 00:000 $00:ffff is ram 
	memory bank 0 is ram

Calling "memory" by itself will print the current memory configuration and is the
same as "show memory".`

	longHelpNMI string = `Trigger a non-maskable interrupt (NMI).

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpReading string = `Set a special address for reading.

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpResume string = `Resume execution of a halted machine.

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpReset string = `Trigger the RESET signal.

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpRun string = `Run an initialized machine.

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpSave string = `Save an address range of memory to a file.

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpSet string = `Set various parameters.

Use this to set verbose and trace modes. 

	set step [on|off]	- Single step mode
	set trace [on|off]	- Print trace (lots and lots of output)
	set verbose [on|off]    - Give more information (more output)

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpShow string = `Show information on the system

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

	longHelpStatus string = `Display high-level status of machine or host.

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpStore string = `Store a byte at a given address in memory.

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`

	longHelpWriting string = `Define special address for writing.

(THE REST OF THIS ENTRY IS MISSING)

Example:
		(THE EXAMPLE IS MISSING)`
)
