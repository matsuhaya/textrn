# textrn

rename multiple files by editor

## Usage

Change directory containing files you want to rename.

```
.
├── aaa.txt
├── bbb.txt
└── zzzzzz
```

```
textrn
```

Rename multiple files by your favorite editor.

```
aaa_20220101.txt
bbb_20220102.txt
zzzzzz
```

```
.
├── aaa_20220101.txt
├── bbb_20220102.txt
└── zzzzzz
```

## Installation

```
$ go install github.com/matsuhaya/textrn/cmd/textrn@latest
```

## License
MIT

## Author
matsuhaya