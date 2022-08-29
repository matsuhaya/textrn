# textrn

rename multiple files with text editor

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

Open temporary text file with editor.

```
aaa.txt
bbb.txt
zzzzzz
```

Rename multiple files with your favorite editor.

```
aaa_20220101.txt
bbb_20220102.txt
zzzzzz
```

You can rename multiple files. 

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