# One-Shot to print opcode tables 
# Scot W. Stevenson <scot.stevenson@gmail.com>
# First version: 05. Jan 2019
# This version: 05. Jan 2019

for a in range(0xF+1):

    for b in range(0xF+1):
        print(f'{a:X}{b:X} ', end=' ')

    print()


