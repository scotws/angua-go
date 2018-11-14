CPU emulated mode files for the Angua 65816 Emulator in Go
Scot W. Stevenson <scot.stevenson@gmail.com>
First version: 06. Nov 2018
This version: 14. Nov 2018

The CPU of the Angua 65816 works like every other part of the system: With
brute force.  There are two separately defined CPUs, one for native mode and
the other for emulated mode. The native CPU does not know the emulated CPU
exists and vice versa. This avoids having to test for various modes, makes
debugging easier.
