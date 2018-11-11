Switcher Daemon for Angua 65816 
Scot W. Stevenson <scot.stevenson@gmail.com>
First version: 11. Nov 2018
This version: 11. Nov 2018

The Switcher Daemon is started when the user calls init from the CLI. It is
run as a goroutine and spawns its own goroutines: The two main loops for
the 8-bit and 16-bit variants (emulated and native) of the 65816 CPUs; for
testing (currently) a "heartbeat" that prints an "I'm alive!" message every
two minutes; and then goes into an infinite loop where it waits for a request
from the CPUs to switch to the other mode. This is so the CLI doesn't have to
be aware of any of this.

