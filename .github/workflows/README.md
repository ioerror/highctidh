GHA Workflows
=============

Written in YS

## Usage

```
$ make update
$ make update force=1
$ make update force=1 devel=1
```

## Overview

Each `.name.yaml` workflow file is generated from a `name.ys` file.

GHA uses the "hidden" `.yaml` file, but you can easily refactor them by using
YS as the inputs.


## Makefile Options

The `make update force=1` command will rebuild all `.yaml` files.

By default `ys` is used from a Docker container which is more portable but
slow.

With `make update devel=1` a local copy of `ys` is installed under the `.git/`
directory the first time and then used.
This runs much faster.


## How to Refactor the Workflows

Simply make changes to any of the input files and then run `make update`.
If no errors occur, run `git status`.
If no `.*.yaml` files have changed, your refactoring worked.

If something did change you made a mistake.
You can undo the changes and try again.

To see if a file was well formed, run:

```
(source <(make env) && ys -Uc name.ys)
```

To see the output a file produces, run:

```
(source <(make env) && ys -Y name.ys)
```

See the [YS Documentation](https://yamlscript.org/doc/) for more debugging
techniques.
