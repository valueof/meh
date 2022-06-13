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
-version
    print version and exit
-withImages
    whether to download images from medium cdn
```

#### How To Help

So far this has been tested only on my own data. If you want to help, head over to [medium.com/me/export](https://medium.com/me/export), export your data, run it through `meh`, and let me know if anything is broken. Thank you!
