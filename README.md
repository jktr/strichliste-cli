# go-strichliste: Go Bindings for the Strichliste API

[![GoDoc](https://godoc.org/github.com/jktr/strichliste-cli?status.svg)](https://godoc.org/github.com/jktr/strichliste-cli)

This program implements a command line client for
[hackerspace bootstrap's strichliste](https://github.com/strichliste/strichliste-backend),
which is a pretty neat tally sheet server.

Please note that this tool is of beta-level quality, but should
suffice for basic use.

There's a sister project to this library,
[go-strichliste](https://github.com/jktr/go-strichliste), which
implements the REST client and API bindings.

## Example

```
$ cat ~/.strichliste-cli/config.json
{
  "api-url": "https://demo.strichliste.org/api",
  "user": "jktr"
}

$ ./strichliste-cli user create --name jktr --balance 5.00
created user #1 (jktr)
created transaction #1
new balance for user #1 (jktr): 5.00€

$ ./strichliste-cli article create --name 'Mate Mate' --value 1.00
created article #1 (Mate Mate)

$ ./strichliste-cli buy --article 1 --count 3
created transaction #2
new balance for user #1 (jktr): 2.00€
```

## License

    Copyright (C) 2019 Konrad Tegtmeier

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    Lesser GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

## v1 API

There's a legacy `strichlist-cli` for the v1 API, which can be found
[here](https://git.cs.uni-paderborn.de/jktr/strichliste-cli). It's
deprecated; you should really migrate to Strichliste v2.

## Acknowledgments

The structure of this program is heavily based on the one of
[Hetzner's hcloud](https://github.com/hetznercloud/cli) tool.
Thanks for open sourcing it!
