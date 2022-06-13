# 😐 MEH: Medium Export Helper

This tool transforms Medium's default data export format (HTML) into JSON. This is work in progress and, so far, being tested only on my own export data.

#### Usage
```
$ go run . -in=/path/to/archive -out=/path/to/out
```

#### Flags

```
-in string
    path to the (uncompressed) medium archive
-out string
    output directory
-verbose
    whether to print logs to stdout
-withImages
    whether to download images from medium cdn
```