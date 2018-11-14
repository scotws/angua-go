Memory subsystem for Angua 
Scot W. Stevenson <scot.stevenson@gmail.com>
First version: 14. Nov 2018
This version: 14. Nov 2018

Memory is divided into "chunks" of continuous memory with a max range of 24 bit.
Chunks are made to be threadsafe for reading and writing. The program will
throw an error if the address is not in a legal range. 

