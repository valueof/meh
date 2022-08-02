# 😐 MEH: Medium Export Helper

This tool transforms Medium's default data export format (HTML) into JSON. This is work in progress and, so far, being tested only on my own export data.

#### Usage

This command will convert original Medium archive located in `/path/to/archive` into JSON with results being written into `/path/to/out`:
```
$ meh -dir=/path/to/archive -out=/path/to/out
```


#### All Flags

```
-dir string
    path to the uncompressed medium archive
-out string
    output directory
-server string
    run web version of meh on provided address
-verbose
    whether to print logs to stdout
-version
    print version and exit
-withImages
    whether to download images from medium cdn
-zip string
    path to the compressed medium archive
```

#### How To Help

So far this has been tested only on my own data. If you want to help, head over to [medium.com/me/export](https://medium.com/me/export), export your data, run it through `meh`, and let me know if anything is broken by creating an issue.

You can also try Medium Export Helper on the web: [meh.antonkovalyov.com](https://meh.antonkovalyov.com)

Thank you!
