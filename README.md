# gophernes

Inspired to try writing an emulator, I put together this Go-based NES emulator over the span of a couple of weeks. This work is largely based on poring over the following resources:

* [Nesdev wiki](http://wiki.nesdev.com/w/index.php/Nesdev_Wiki)
* [obelisk.me 6502 documentation](http://www.obelisk.me.uk/6502/index.html)
* [oxyron.de opcode matrix](http://www.oxyron.de/html/opcodes02.html)
* More that I am likely forgetting now

I was able to come up with a functioning emulator that can read a couple of simple cartridge formats, emulate the full 6502 CPU instruction set, and emulate the NES PPU (graphics processor). Notably missing are the APU (audio), and controller inputs. I tested using a Pacman ROM only, so there are likely incompatibilities with other ROMs, but it is able to successfully boot and display the demo.

I'm using [ebiten](https://github.com/hajimehoshi/ebiten) to display the basic graphics onscreen.

There are some "quality of life" command line options for development, such as running headlessly, tracing CPU/PPU instructions, and the ability to run the emulation at any rate. I was working on the beginnings of the APU, but this was never completed.
