# devtime

Track your coding time from the terminal. Local-only, no account, no cloud.

## Prerequisites

Devtime currently works with **VS Code only**. Install the extension first:

[devtime for VS Code](https://marketplace.visualstudio.com/items?itemName=arnaudhrt.devtime-local)

The extension runs in the background and tracks your coding activity locally in `~/.devtime/`.

## Installation

**Mac (Apple Silicon):**

```bash
mkdir -p ~/.local/bin && curl -sSL https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_darwin_arm64.tar.gz | tar xz -C ~/.local/bin devtime
```

**Mac (Intel):**

```bash
mkdir -p ~/.local/bin && curl -sSL https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_darwin_amd64.tar.gz | tar xz -C ~/.local/bin devtime
```

**Linux (amd64):**

```bash
mkdir -p ~/.local/bin && curl -sSL https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_linux_amd64.tar.gz | tar xz -C ~/.local/bin devtime
```

**Linux (arm64):**

```bash
mkdir -p ~/.local/bin && curl -sSL https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_linux_arm64.tar.gz | tar xz -C ~/.local/bin devtime
```

> Make sure `~/.local/bin` is in your `PATH`. Add `export PATH="$HOME/.local/bin:$PATH"` to your `~/.zshrc` or `~/.bashrc` if needed.

**Windows (PowerShell):**

```powershell
cd ~
curl https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_windows_amd64.zip -OutFile devtime.zip
Expand-Archive devtime.zip -Force -DestinationPath $env:LOCALAPPDATA\Microsoft\WindowsApps
Remove-Item devtime.zip
```

**From source:**

```bash
go install github.com/arnaudhrt/devtime@latest
```

## Usage

### `devtime all`

```
$ devtime all

  Total time:     124h 10m
  Daily average:  2h 45m
  Days tracked:   45
  Tracking since: Mar 11, 2026

  Projects:
    my-app   80h 20m  █████████████░░░░░░░   65%
    devtime  28h 40m  ████░░░░░░░░░░░░░░░░   23%
    my-proj  15h 10m  ██░░░░░░░░░░░░░░░░░░   12%

  Languages:
    TypeScript  72h 30m  ███████████░░░░░░░░░   58%
    Go          38h 25m  ██████░░░░░░░░░░░░░░   31%
    CSS         13h 15m  ██░░░░░░░░░░░░░░░░░░   11%
```

### `devtime`

Running `devtime` with no arguments shows today's breakdown.

```
$ devtime

  Today: 4h 23m

  Projects:
    my-app   2h 45m  ██████████████░░░░░░   63%
    devtime  1h 12m  ██████░░░░░░░░░░░░░░   27%
    my-proj  0h 26m  ██░░░░░░░░░░░░░░░░░░   10%

  Languages:
    TypeScript  2h 50m  █████████████░░░░░░░   65%
    Go          1h 07m  █████░░░░░░░░░░░░░░░   26%
    CSS         0h 26m  ██░░░░░░░░░░░░░░░░░░    9%
```

### `devtime week` / `devtime month [mmm-yyyy]` / `devtime year [yyyy]`

`devtime week` shows the current week. `devtime month` and `devtime year` show the current month/year, or pass an argument to look up a specific one:

```
$ devtime month nov-2025

  November 2025: 42h 15m

  Projects:
    my-app   28h 10m  █████████████░░░░░░░   67%
    devtime  14h 05m  ███████░░░░░░░░░░░░░   33%

  Languages:
    TypeScript  30h 00m  ██████████████░░░░░░   71%
    Go          12h 15m  ██████░░░░░░░░░░░░░░   29%
```

```
$ devtime year 2025

  2025: 312h 40m

  Projects:
    my-app   180h 20m  ███████████░░░░░░░░░   58%
    devtime   82h 10m  █████░░░░░░░░░░░░░░░   26%
    my-proj   50h 10m  ███░░░░░░░░░░░░░░░░░   16%

  Languages:
    TypeScript  190h 30m  ████████████░░░░░░░░   61%
    Go           78h 25m  █████░░░░░░░░░░░░░░░   25%
    CSS          43h 45m  ███░░░░░░░░░░░░░░░░░   14%
```

### `devtime status`

```
$ devtime status

  Status: active
  Project:  my-app
  Language: typescript
  Editor:   vscode
  Session:  0h 45m
```

### `devtime project [name]`

Without a name, opens an interactive picker. With a name, shows the breakdown directly.

```
$ devtime project
? Select a project:
> my-app
  devtime
  wannee
```

```
$ devtime project my-app

  All time:    80h 20m
  This month:  24h 10m
  This week:    8h 30m

  Languages:
    TypeScript  18h 30m  ███████████████░░░░░   77%
    CSS          3h 40m  ███░░░░░░░░░░░░░░░░░   15%
    JSON         2h 00m  ██░░░░░░░░░░░░░░░░░░    8%
```

### `devtime lang [name]`

Without a name, opens an interactive picker. With a name, shows the breakdown directly.

```
$ devtime lang
? Select a language:
> Go
  TypeScript
  CSS
```

```
$ devtime lang go

  All time:    38h 25m
  This month:  12h 15m
  This week:    8h 30m

  Projects:
    devtime  5h 20m
    my-proj  3h 10m
```

### `devtime wakatime-import <file>`

Import your coding history from a WakaTime JSON export. Data is converted into devtime's compacted monthly format and merged with any existing data.

```
$ devtime wakatime-import ~/Downloads/wakatime-export.json

Parsing WakaTime export: /Users/you/Downloads/wakatime-export.json
  2025-03: 18h 29m
  2025-04: 93h 40m
  2025-05: 70h 54m
  2025-06: 138h 44m
  ...

Imported 13 month(s), total: 1051h 23m
```

To export your data from WakaTime, go to [wakatime.com/settings/account](https://wakatime.com/settings/account) and export **Daily Totals** (not Heartbeats — heartbeat exports use a different format and won't work).

### `devtime doctor`

Check if the VS Code extension is sending events correctly.

```
$ devtime doctor

  Last event:
    Type:     heartbeat
    Project:  my-app
    Language: typescript
    Editor:   vscode
    Time:     Mar 16 14:32:05 (3m ago)
```

## Data

All data is stored locally in `~/.devtime/` as plain text JSONL files, one per month.

```
~/.devtime/
├── events-2026-01.jsonl
├── events-2026-02.jsonl
└── events-2026-03.jsonl
```
