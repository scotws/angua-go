# Angua - An Emulator for 65816 Native Mode 
Scot W. Stevenson <scot.stevenson@gmail.com>
First version: 26. Sep 2017
This version: 01. Jan 2019

Angua is an emulator for the native mode of the 65816 CPU, a 8/16-bit hybrid
processor that is the sibling of the famous 6502 8-bit processor of the Apple
II, VIC-20, C64, and other machines. 

The 65816 is a complex processor which can switch between a 6502-emulation mode
and a 65816-native mode, mutate the size of the A, X, and Y registers, and
function in binary or decimal mode for mathematical operations. However, in
practice, neither emulated nor decimal mode are used much. Angua attempts to
simplify the problem by _only_ emulating native mode and dropping decimal mode
_completely._

Angua is written in Go (golang), and in fact a major part of the project was a
opportunity to learn Go better. Therefore, early parts of the system are rather
crude. Also, every attempt was made to keep the code as straightforward and
easy to understand as possible. Where possible, "brute force" coding was used to
keep every routine easy to understand at the cost of repetitive code in some
places.  It is hoped that readers not that familiar with Go will be able to
understand the code that way. 

No serious attempt has been made at this point to make the system more efficient
or fast. In the compiler we trust. 

The emulator comes with online help. Type "help" to get a list of shell
commands, "<COMMAND> help" for more information on individual shell commands,
and "info <MNEMONIC>" or "info <OPCODE>" for information on mnemonics and
opcodes.

