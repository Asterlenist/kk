# kk

`kk` is a small command-line utility for common package-manager and systemd tasks. Native package operations are currently supported only through APT and DNF/YUM on Linux.

## Requirements

- Go 1.27 or newer to build from source.
- APT (`apt` or `apt-get`) on Debian-family Linux distributions, or DNF/YUM on RPM-family Linux distributions.
- Flatpak for Flathub installs.
- Linux with systemd for service commands.
- Administrative privileges where the selected system utility requires them. `kk` does not invoke `sudo` or elevate privileges automatically.

## Build

```sh
go build -o kk .
```

Run `./kk` without arguments to print the command summary.

## Commands

```text
kk system
kk pkg-inst <package> [package...]
kk pkg-rem <package> [package...]
kk pkg-pur <package> [package...]
kk flatpak-inst <app-id> [app-id...]
kk service-add <service> [service...]
kk service-rem <service> [service...]
```

- `system` prints the detected operating system, Linux distribution and available package manager.
- `pkg-inst` installs one or more packages using the detected native package manager.
- `pkg-rem` removes packages. `pkg-pur` also requests configuration and dependency cleanup where the package manager supports it.
- `flatpak-inst` installs application IDs from the Flathub remote. Flatpak must already be installed and configured with Flathub.
- `service-add` enables and starts systemd services immediately.
- `service-rem` stops and disables systemd services. It does not delete unit files.

Package operations use the available APT or DNF/YUM executable. Administrative privileges may be required; `kk` does not invoke `sudo` or elevate privileges automatically. Review package names and system changes before running install, remove, or purge commands.

## Command history

`kk` automatically records each invocation in `~/.local/share/kk/history.log`. On first run it creates `~/.config/kk/config.toml` with these settings:

```toml
[logging]
enabled = true
path = "~/.local/share/kk/history.log"
max_size = 0
```

- Set `enabled = false` to disable command history.
- Set `path` to change the log file location. A path beginning with `~/` is resolved relative to your home directory; other relative paths are relative to the current working directory.
- `max_size` is the maximum log size in bytes. `0` means unlimited. Set it to at least `128` to cap the file; when the cap is reached, old complete entries are removed so the newest history is retained. A command record longer than the cap is shortened.

## Supported package managers

The detector recognizes only `apt-get`/`apt` and `dnf`/`yum`. On Linux distributions using other package managers, and on macOS, Windows, and FreeBSD, `kk system` still reports the operating system but does not select a package manager; native package commands return an error. Flatpak and systemd service commands are separate from native package-manager support.
