# devtime

A local-first coding time tracker for the terminal. Reads event files written by the [devtime VS Code extension](../devtime-vscode/) and displays coding time stats with bar charts.

No daemon, no database, no cloud — just files and a CLI.

## How It Works

1. The **VS Code extension** sends coding events (heartbeat, focus, blur) to a local agent
2. The **agent** (`devtime serve`) writes events as JSONL to `~/.devtime/events-YYYY-MM.jsonl`
3. The **CLI** reads those files, computes sessions on the fly, and displays stats

### Session Rules

- Events within 5 minutes of each other belong to the same session
- A `blur` event ends the current session immediately
- A project or language change starts a new session
- A single isolated event counts as 30 seconds

### Event File Format

One JSON object per line in `~/.devtime/events-YYYY-MM.jsonl`:

```json
{"ts":"2026-03-11T09:00:00+01:00","event":"heartbeat","project":"my-app","lang":"typescript","editor":"vscode"}
```

## Installation

Install the official release binary by copy/pasting the appropriate command below into a terminal.

On a **Mac with an Apple Silicon (M) processor**:

```bash
curl -sSL https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_darwin_arm64.tar.gz | tar xz -C /usr/local/bin devtime
```

On an **older Mac with an Intel processor**:

```bash
curl -sSL https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_darwin_amd64.tar.gz | tar xz -C /usr/local/bin devtime
```

On a **Linux machine with an Intel/AMD processor**:

```bash
curl -sSL https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_linux_amd64.tar.gz | tar xz -C /usr/local/bin devtime
```

On a **Linux machine with an ARM processor**:

```bash
curl -sSL https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_linux_arm64.tar.gz | tar xz -C /usr/local/bin devtime
```

On a **Windows machine** (in PowerShell):

```powershell
cd ~
curl https://github.com/arnaudhrt/devtime/releases/latest/download/devtime_windows_amd64.zip -OutFile devtime.zip
Expand-Archive devtime.zip -Force -DestinationPath $env:LOCALAPPDATA\Microsoft\WindowsApps
Remove-Item devtime.zip
```

You should now be able to run:

```
$ devtime --help
Track your coding time from the terminal
```

### From source

If you have Go installed, you can also install with:

```bash
go install github.com/arnaudhrt/devtime@latest
```

## Commands

### `devtime status`

Show whether you're currently coding and info about the active/last session.

```
$ devtime status

  Status: active
  Project:  my-app
  Language: typescript
  Editor:   vscode
  Session:  0h 45m
```

If inactive:

```
$ devtime status

  Status: not active
  Last session: 0h 32m on my-app (17:04)
```

### `devtime today`

Show today's total coding time with breakdowns by project and language.

```
$ devtime today

  Today: 4h 23m

  Projects:
    befitwithjess  2h 45m  ██████████████░░░░░░   63%
    wannee         1h 12m  ██████░░░░░░░░░░░░░░   27%
    devtime        0h 26m  ██░░░░░░░░░░░░░░░░░░   10%

  Languages:
    TypeScript  2h 50m  █████████████░░░░░░░   65%
    Go          1h 07m  █████░░░░░░░░░░░░░░░   26%
    CSS         0h 26m  ██░░░░░░░░░░░░░░░░░░    9%
```

### `devtime week`

Show this week's coding time (Monday through today).

```
$ devtime week

  This Week: 18h 42m

  Projects:
    ...

  Languages:
    ...
```

### `devtime month`

Show this month's coding time.

```
$ devtime month

  This Month: 62h 15m

  Projects:
    ...

  Languages:
    ...
```

### `devtime project <name> <all|month|week>`

Show coding time for a specific project, broken down by language.

```
$ devtime project my-app month

  my-app — This Month: 24h 10m

  Projects:
    my-app  24h 10m  ████████████████████  100%

  Languages:
    TypeScript  18h 30m  ███████████████░░░░░   77%
    CSS          3h 40m  ███░░░░░░░░░░░░░░░░░   15%
    JSON         2h 00m  ██░░░░░░░░░░░░░░░░░░    8%
```

### `devtime lang <name> <all|month|week>`

Show coding time for a specific language, broken down by project.

```
$ devtime lang Go week

  Go — This Week: 8h 30m

  Projects:
    devtime  5h 20m  ████████████░░░░░░░░   63%
    wannee   3h 10m  ███████░░░░░░░░░░░░░   37%

  Languages:
    Go  8h 30m  ████████████████████  100%
```

## VS Code Extension

Install the [devtime extension](https://marketplace.visualstudio.com/items?itemName=arnaudhrt.devtime) from the VS Code Marketplace.

The extension will:
- Send **heartbeat** events every 30 seconds while you code
- Send **focus**/**blur** events when the VS Code window gains/loses focus
- Auto-start the `devtime serve` agent if it's not running

## Data Storage

All data lives in `~/.devtime/`:

```
~/.devtime/
├── events-2026-01.jsonl
├── events-2026-02.jsonl
└── events-2026-03.jsonl
```

One file per month. Plain text, easy to inspect, back up, or delete.
