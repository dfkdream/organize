# Organize
Template-based file organizer
## Build
```shell
$ git clone github.com/dfkdream/organize && cd organize
$ go build
```
## Usage

For template syntax, see `text/template` [documentation](https://pkg.go.dev/text/template).

### Flags

| Flag | Usage                 | Default Value      |
| ---- | --------------------- | ------------------ |
| `-i` | Input directory       | `.`                |
| `-o` | Output directory      | `.`                |
| `-p` | File name pattern     | `*`                |
| `-r` | Enable recursive walk | `false`            |
| `-t` | Template              | `{{ .Info.Name }}` |
| `-dry-run` | Do not apply result to filesystem | `false` |

### Template Variables
| Name | Description |
| ---- | ----------- |
| `.From` | Original file path |
| `.Info` | Same with [`fs.FileInfo`](https://pkg.go.dev/io/fs#FileInfo) |

### Template functions

| Name | Input | Output | Usage |
| ---- | ----- | ------ | ----- |
| `count` | `format string` | `string` | Return incremental count with provided `fmt.Printf` format | 
| `ext` | `file name or path` | `string` | Return file extension with `.` prefix (e.g. `.jpg`) |
| `chTimes` | `file path, time.Time` | `none` | Set file modification timestamp |
| `parseTime` | `format string, string` | `time.Time` | Parse time string (same with `time.Parse`) |
| `skip` | `none` | `none` | Skip file |

### Example
Add `-dry-run` argument to prevent unwanted file loss
* Parse filename as timestamp, change ModTime of every `.jpg` files to parsed timestamp, and move to `../result/YYYY/MM/YYYY-MM-DD (4 Digit Hex Count).jpg`
```shell
$ organize -r -t '{{ $t := (parseTime "20060102.jpg" .Info.Name) }}{{ chTimes .From $t }}{{ $t.Format "2006/01/2006-01-02" }} {{ count "%04X" }}{{ ext .Info.Name }}' -p "*.jpg" -i . -o ../result
```
