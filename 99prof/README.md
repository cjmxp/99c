# Table of Contents

1. Usage
1. Installation
1. Changelog

# 99prof

Command 99prof profiles programs produced by the 99c compiler.

The profile is written to stderr.

### Usage

    99prof [-functions] [-lines] [-instructions] [-rate] a.out [arguments]

### Installation

To install or update 99prof

     $ go get [-u] -tags virtual.profile github.com/cznic/99c/99prof

Online documentation: [godoc.org/github.com/cznic/99c/99prof](http://godoc.org/github.com/cznic/99c/99prof)

### Changelog

2017-10-09: Initial public release.
