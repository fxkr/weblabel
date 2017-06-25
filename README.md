# weblabel - A web interface for label printers

![image](https://cloud.githubusercontent.com/assets/389616/25075287/7a33bc96-2311-11e7-91b2-8dedc52a7c8b.png)


## Installation

1. Add the [packagecloud.io/fxkr/weblabel](https://packagecloud.io/fxkr/weblabel/install) repository and install `weblabel` via your package manager. Currently, Debian 8 (x86_64, armhf) and Fedora 25 (x86_64). Contact me if you need more.

2. Install a tool like [ptouch-print](https://github.com/dradermacher/ptouch-print) that takes a text to print on the command line. At the moment, you'll likely need to compile this yourself. This won't be necessary in the future.

3. Copy `/usr/share/weblabel/config.yml` to `/etc/weblabel/config.yml` and edit appropriately. For `ptouch-print`, it could look like this:

    ```yaml
    ---
    Address: "0.0.0.0:80"
    PrintCommand: "/usr/bin/ptouch-print --text {}"
    StaticPath: "/usr/share/weblabel/static"
    ```

4. (Re-)start and enable (if needed) weblabel:

    ```sh
    systemctl restart weblabel
    systemctl enable weblabel
    ```

5. If you use a Brother 2430PC, make sure the switch is in the "E" (not "EL") position.


## API

Note: the API is not stable yet.

Printing text:

```sh
curl 'http://localhost:8081/api/v1/printer/print' -H 'Accept: application/json' --data '{"document": {"text":"bar"}}'
```

Rendering a preview:

```sh
curl 'http://localhost:8081/api/v1/renderer/render' -H 'Accept: application/json' --data '{"document": {"text":"bar"}}'
```

Printing a PNG file:

```sh
curl 'http://localhost:8081/api/v1/printer/image' -H 'Accept: application/json' -F 'data={}' -F 'image=@label.png'
```


## Development

You'll need [go](https://golang.org/) for the backend and  [yarn](https://yarnpkg.com/lang/en/) for the frontend.

It's recommended to use a simple image viewer like [feh](https://feh.finalrewind.org/) to simulate printing.

The following sections show some useful commands.

### Configuration

Using a real labelprinter would get expensive fast. `notify-send`, which shows an OSD, provides a good alternative. Use port 8081 for the backend - the frontend dev server at port 8080 will proxy API requests to there.

Put this in `config.yml`:


```yaml
---
Address: "127.0.0.1:8081"
PrintCommand: "feh {}"
StaticPath: ./static/dist/
```

### Backend

```
export GOPATH="$(pwd)"
```

Fetching source code and dependencies:

```sh
go get github.com/fxkr/weblabel
```

Running directly:
```sh
go run cmd/weblabel/weblabel.go
```

Running tests:
```sh
go test ./...
```

Compiling a binary:
```sh
go build github.com/fxkr/weblabel/cmd/weblabel
```

Cross compiling for Raspberry Pi:
```sh
GOARM=6 GOOS=linux GOARCH=arm go build -v github.com/fxkr/weblabel/cmd/weblabel
```

### Frontend

```sh
cd static
```

Fetching dependencies:
```sh
yarn install
```

Running a dev webserver:
```sh
export PATH="$PATH:$(yarn bin)"
webpack-dev-server
```
