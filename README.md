Just Pikd WMS
=========

Development Setup:
-------------
1. Currently just boot the vagrant image from the just-pikd repo
2. Clone this repo into ~/Documents/just-pikd-wms
3. `vagrant ssh` to the dev box
4. Need to fix permissions via chef, but do this after provisioning: `sudo chown -R vagrant:vagrant /opt/go`
5. Add to ~/.bashrc (and source it)
    `export GOPATH='/opt/go'
    `export PATH=$GOPATH:$PATH
6. I use gin for live code reloading: `go get github.com/codegangsta/gin`
7. Dependency management: `go get github.com/tools/godep`
8. Set cwd to the WMS directory: `cd /opt/go/src/just-pikd-wms`
9. Install deps: `godep restore`
10. Run `gin` to start the app, and it should reload when you change code (it can be somewhat slow).

I'm also using GoSublime. I had to set GOPATH and install some packages on my mac to get formatting/syntax highlighting to work. Standardizing that is another TODO item. In GoSublime package settings, I did this:

```
{
        "shell": ["/usr/bin/bash"],
        "env": {"GOPATH": "/Users/Scott/go/",
        "GOROOT": "/opt/local/go"},
}
```

Running Tests:
-------------
Use `./test` to run all tests, or the `go test` command to get finer grained control over running individual tests

Test Data Generation:
-------------
In the test_data folder, generate_sql.xlsx was used to generate some static records.
generate_test_data.sql contains statements based off of that to generate additional data, as well as examples of how to dump it to json files
