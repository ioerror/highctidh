GHA Workflows
=============

Written in YS

## Usage

```
$ make update
```

## Overview

Each `.name.yaml` workflow file is generated from a `name.ys` file.

GHA uses the "hidden" `.yaml` file, but you can easily refactor them by using
YS as the inputs.


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
