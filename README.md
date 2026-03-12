# devtime

Track your coding time from the terminal. Local-only, no account, no cloud.

## Prerequisites

Devtime currently works with **VS Code only**. Install the extension first:

[devtime for VS Code](https://marketplace.visualstudio.com/items?itemName=arnaudhrt.devtime-local)

The extension runs in the background and tracks your coding activity locally in `~/.devtime/`.

## Installation

**Mac (Apple Silicon):**

```bash
curl -sSL https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_darwin_arm64.tar.gz | sudo tar xz -C /usr/local/bin devtime
```

**Mac (Intel):**

```bash
curl -sSL https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_darwin_amd64.tar.gz | sudo tar xz -C /usr/local/bin devtime
```

**Linux (amd64):**

```bash
curl -sSL https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_linux_amd64.tar.gz | sudo tar xz -C /usr/local/bin devtime
```

**Linux (arm64):**

```bash
curl -sSL https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_linux_arm64.tar.gz | sudo tar xz -C /usr/local/bin devtime
```

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

### `devtime profile`

```
$ devtime profile

  Devtime Profile

  Tracking since: Mar 11, 2026
  Total time:     124h 10m
  Daily average:  2h 45m
  Days tracked:   45

  Projects:
    my-app   80h 20m  █████████████░░░░░░░   65%
    devtime  28h 40m  ████░░░░░░░░░░░░░░░░   23%
    my-proj  15h 10m  ██░░░░░░░░░░░░░░░░░░   12%

  Languages:
    TypeScript  72h 30m  ███████████░░░░░░░░░   58%
    Go          38h 25m  ██████░░░░░░░░░░░░░░   31%
    CSS         13h 15m  ██░░░░░░░░░░░░░░░░░░   11%
```

### `devtime today` / `week` / `month`

```
$ devtime today

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

### `devtime status`

```
$ devtime status

  Status: active
  Project:  my-app
  Language: typescript
  Editor:   vscode
  Session:  0h 45m
```

### `devtime projects`

Interactive project picker. Select a project to see its breakdown.

```
$ devtime projects
? Select a project:
> my-app
  devtime
  wannee
```

### `devtime project <name>`

```
$ devtime project my-app

  Devtime for my-app

  All time:    80h 20m
  This month:  24h 10m
  This week:    8h 30m

  Languages:
    TypeScript  18h 30m  ███████████████░░░░░   77%
    CSS          3h 40m  ███░░░░░░░░░░░░░░░░░   15%
    JSON         2h 00m  ██░░░░░░░░░░░░░░░░░░    8%
```

### `devtime langs`

Interactive language picker. Select a language to see its breakdown.

```
$ devtime langs
? Select a language:
> Go
  TypeScript
  CSS
```

### `devtime lang <name>`

```
$ devtime lang go

  Devtime for GO

  All time:    38h 25m
  This month:  12h 15m
  This week:    8h 30m

  Projects:
    devtime  5h 20m
    my-proj  3h 10m
```

## Data

All data is stored locally in `~/.devtime/` as plain text JSONL files, one per month.

```
~/.devtime/
├── events-2026-01.jsonl
├── events-2026-02.jsonl
└── events-2026-03.jsonl
```
