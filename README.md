# YT Toolbox

A command line program to extract youtube subscriptions, playlist and history. Using [go-is](https://github.com/J-Siu/go-is) and [rod](https://github.com/go-rod/rod).

- [Install](#install)
- [Usage](#usage)
- [Limitation](#limitation)
- [License](#license)
<!--more-->

### Install

Go install

```sh
go install github.com/J-Siu/yt-toolbox@latest
```

Download

- https://github.com/J-Siu/yt-toolbox/releases

### Usage

```sh
YouTube toolbox

Usage:
  go-yt-toolbox [command]

Aliases:
  go-yt-toolbox, gyt

Available Commands:
  completion   Generate the autocompletion script for the specified shell
  config       Print configurations
  help         Help about any command
  history      Get Youtube History
  playlist     Get Youtube Playlist
  subscription Youtube Subscriptions

Flags:
      --DevToolsHost string   Devtools Host
      --DevToolsPort int      Devtools Port (default 9222)
      --DevToolsUrl string    Devtools Url (override Devtools host, ver and port)
      --DevToolsVer string    Devtools Version
  -c, --config string         Config file (default "$HOME/.config/yt-toolbox.json")
      --debug                 Enable debug
  -h, --help                  help for go-yt-toolbox
  -s, --scroll-max int        Unlimited -1
      --trace                 Enable trace (include debug)
  -v, --verbose               Verbose
      --version               version for go-yt-toolbox

Use "yt-toolbox [command] --help" for more information about a command.
```

### Limitation

> Must use remote browser as function require youtube login.

1. Close running Chrome. Start Chrome with following option:

    ```sh
    chrome --remote-debugging-port=9222
    ```

    Or Chromium with the same option.

2. Login youtube.com
3. Run yt-toolbox

### License

The MIT License (MIT)

Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
