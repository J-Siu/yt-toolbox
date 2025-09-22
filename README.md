# YT Toolbox

A command line program to extract youtube subscriptions, playlist and history. Using [go-is](https://github.com/J-Siu/go-is) and [rod](https://github.com/go-rod/rod).

### Limitation

> Must use remote browser as function require youtube login.

1. Close running Chrome. Start Chrome with following option:

    ```sh
    chrome --remote-debugging-port=9222
    ```

    Or Chromium with the same option.

2. Login youtube.com
3. Run yt-toolbox

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
