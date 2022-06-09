# Simple CSV carver/grep

Works on cells and headers using regex

```
Usage:
  cmd [flags]

Flags:
  -a, --cell stringArray
  -f, --filename string
  -p, --filter stringArray
  -h, --help                 help for cmd
  -n, --noheader             CSV has no header
```

# Examples

```sh
$ cat test.csv
#domain, sending-spam, seen-before
"example.com", no, yes
example.eu, yes, yes
$ csvcarve -f test.csv -p '#domain==\.com'
#domain,sending-spam,seen-before
example.com,no,yes
$ csvcarve -f test.csv -p '#domain==\.eu' -a '1!=yes'
#domain,sending-spam,seen-before
$ csvcarve -f test.csv -p '#domain==\.eu' -a '1==yes'
#domain,sending-spam,seen-before
example.eu,yes,yes
```

# Install

```sh
$ go install github.com/desdic/csvcarve/cmd/csvcarve@latest
```
