
utilities for parsing OpenStreetMap PBF files and extracting geographic data

### installation

> tested on go version go1.9, recommended go1.10, cross-compilation possibly broken in recent versions of go

```bash
$ go get github.com/missinglink/pbf
```

### disclaimer

this repo contains a bunch of experiments including a monkey-patched parsing library, it should be considered as permanently in development.

please feel free to fork and contribute, at some stage it might reach a level of maturity and stability where it could be published as a more general-purpose tool.

### running the cli

```bash
$ pbf --help

NAME:
   pbf - utilities for parsing OpenStreetMap PBF files and extracting geographic data

USAGE:
   pbf [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     stats                    pbf statistics
     json                     convert to overpass json format, optionally using bitmask to filter elements
     json-flat                convert to a json format, compulsorily using bitmask to filter elements and leveldb to denormalize where possible
     xml                      convert to osm xml format, optionally using bitmask to filter elements
     opl                      convert to opl, optionally using bitmask to filter elements
     nquad                    convert to nquad, optionally using bitmask to filter elements
     cypher                   convert to cypher format used by the neo4j graph database, optionally using bitmask to filter elements
     sqlite3                  import elements in to sqlite3 database, optionally using bitmask to filter elements
     leveldb                  import elements in to leveldb database, optionally using bitmask to filter elements
     genmask                  generate a bitmask file by specifying feature tags to match
     genmask-boundaries       generate a bitmask file containing only elements referenced by a boundary:administrative relation
     genmask-super-relations  generate a bitmask file containing only relations which have at least one another relation as a member
     bitmask-stats            output statistics for a bitmask file
     store-noderefs           store all node refs in leveldb for records matching bitmask
     boundaries               write geojson osm boundary files using a leveldb database as source
     xroads                   compute street intersections
     streets                  export street segments as merged linestrings, encoded in various formats
     noderefs                 count the number of times a nodeid is referenced in file
     index                    index a pbf file and write index to disk
     index-info               display a visual representation of the index file
     find                     random access to pbf
     help, h                  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### get more detailed information on a specific command

```bash
$ pbf help stats

NAME:
   pbf stats - pbf statistics

USAGE:
   pbf stats [command options] [arguments...]

OPTIONS:
   --interval value, -i value  write stats every i milliseconds (default: 0)
```

### running the tests

```bash
$ go test $(go list ./... | grep -v /vendor/)
```

### issues / bugs

please open a github issue / open a pull request.

if you are planning a non-trivial feature, please open an issue to discuss it first.

### license

```bash
The MIT License (MIT)

Copyright (c) 2017, Peter Johnson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
