# redgreen
Red-green testing with (some) batteries included.

## quickstart

Download/build/install:

    go get -u github.com/jncornett/redgreen/cmd/redgreen
    
Run the server:

    redgreen serve
    
Profit!

You can use `redgreen getall` to examine keys.
The web client is located at `http://localhost:8080` by default.

## todo
- add motivation
- add quickstart
- add overlayfs for customization
- add error parsing plugins
- add (better) server logging and verbosity switches
- add systemd service file
- add filtering search bar for static web client
- add sorting for static web client
