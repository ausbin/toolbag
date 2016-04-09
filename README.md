toolbag
-------

This repository holds the source for some Go utilities I use on my
website. You can find a listing [here][1]. Just for kicks, I've licensed
them under the AGPLv3.

`notfound` and `figlet` are probably the most fun. The former requires
texlive, and the latter requires figlet and my
[`code.austinjadams.com/execd`][2] package.

### how the figlet tool works

where `<--->` is a socket and `==>` indicates executing another program:

#### in production

    (internet)
        ^
        | tcp :80, :443
        v
      nginx
        ^
        | fastcgi (unix socket)
        |            ___________________
        v           |                   |
     toolbag<-------|-->execd => figlet |
               tcp  |___________________|
                   systemd-nspawn container

For more information on the systemd-nspawn container, see the [figlet
README][3] in my [execd repository][2].

#### in development

    ./run => execd => figlet
      ||       ^
      ||       | tcp
      ||       v
      |====> toolbag<----------->your browser
                      tcp :8030

### trying it out

to try toolbag:

    # apt-get install figlet texlive
    $ go get code.austinjadams.com/toolbag
    $ go get code.austinjadams.com/execd
    $ cd $GOPATH/src/code.austinjadams.com/execd/execd
    $ go build
    $ cd $GOPATH/src/code.austinjadams.com/toolbag/toolbag
    $ go build
    $ pushd share
    $ make
    $ popd
    $ ./run
    $ firefox localhost:8030/tools

or, in words:

 0. install figlet, which the figlet tool requires, and a compatible
    latex distribution (probably just texlive), which you'll need for
    building the notfound images
 1. use `go get` (or plain `git` or whatever) to download
    `code.austinjadams.com/toolbag` and `/execd`.  (execd is a
    dependency of the figlet tool. unfortunately, `go get` doesn't seem
    to understand the dependencies of subpackages, so you'll have to
    retrieve it by hand.)
 2. `go build` `code.austinjadams.com/execd/execd` and
    `code.austinjadams.com/toolbag/toolbag` (these are not typos -- the
    binaries are sub-packages)
 3. run the Makefile in the `tools/share` directory (`toolbag/share`, a
    symlink, points there)
 4. from `/toolbag`, call the `run` shellscript to start a minimal http
    server and an execd server.
 5. make http requests to `localhost:8030` to enjoy life

[1]: https://austinjadams.com/tools
[2]: https://code.austinjadams.com/execd
[3]: https://code.austinjadams.com/execd/plain/execd/figlet/README
