# The Angua emulator for the 65816 (Go version)
Scot W. Stevenson <scot.stevenson@gmail.com>
First version: 26. Sep 2017
This version: 11. Nov 2018

Angua is an emulator for the 65816 CPU, a 8/16-bit hybrid processor that is the
sibling of the famous 6502 8-bit processor of the Apple II, VIC-20, C64, and
other machines. 

Angua is written in Go (golang), and in fact a major part of the project was a
opportunity to learn Go better. Therefore, early parts of the system are rather
crude. 

The 65816 is a complex processor which switch between 8 and 16 bit modes and
mutate the size of the A, X, and Y registers. Therefore, every attempt was made
to keep the code as straightforward and easy to understand as possible. Where
possible, "brute force" coding was used to keep every routine easy to understand
at the cost of massive repetitive code in some places. It is hoped that readers
not that familiar with Go will be able to understand the code that way. 

No serious attempt has been made at this point to make the system more efficient
or fast -- in the compiler we trust. 

The emulator comes with online help. Type "help" to get a list of shell
commands, "<COMMAND> help" for more information on individual shell commands,
and "info <MNEMONIC>" or "info <OPCODE>" for information on mnemonics and opcodes.

