# ps4-20th-tool

A small tool built for educational purposes that attempts to pick apart Sony's
20 Years Of Character's website to find the secret URL. Chances are very high
this tool no longer works as details about the how the website works are
widely known.

Dean Wild has the best summary of the internals of the website that I've seen
so far in his [blog post][] on the matter.

[blog post]: https://deano2390.wordpress.com/2014/12/17/hacking-that-playstation-competition/


## Installation

### Binaries

Binaries are available for download on the
[Releases](https://github.com/jimeh/ps4-20th-tool/releases) page for:

- Mac OS X / Darwin (x86, amd64)
- FreeBSD (x86, amd64, arm)
- Linux (x86, amd64, arm)
- Windows (x86, amd64)
- NetBSD (x86, amd64, arm)
- OpenBSD (x86, amd64)
- Plan9 (x86)

### From Source

```bash
go get github.com/jimeh/ps4-20th-tool
```


## Usage

```
usage: ps4-20th-tool <command>

Commands:
   lookup  Lookup the SP (redirect code) and the secret URL.
   brute   Attempt to a brute force attack against the redirect page, trying
           every possible combination of 2 and 3 characters.
```


## License

```
        DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
                    Version 2, December 2004

 Copyright (C) 2014 Jim Myhrberg

 Everyone is permitted to copy and distribute verbatim or modified
 copies of this license document, and changing it is allowed as long
 as the name is changed.

            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION

  0. You just DO WHAT THE FUCK YOU WANT TO.
```
