# kk

`kk` is a small command-line utility for common package-manager and systemd tasks. Native package operations are supported through APT, DNF/YUM, and pacman on Linux.

## Requirements

- Go 1.27 or newer to build from source.
- APT (`apt` or `apt-get`) on Debian-family distributions, DNF/YUM on RPM-family distributions, or pacman on Arch-family Linux distributions.
- Flatpak for Flathub installs.
- Linux with systemd for service commands.
- Administrative privileges where the selected system utility requires them. `kk` does not invoke `sudo` or elevate privileges automatically.

## Build

```sh
go build -o kk .
```

Run `./kk` without arguments to print the command summary.

## Install

With Go and Git installed, build the binary and install it manually:

```sh
sudo install -m 0755 kk /usr/local/bin/kk
```

Or use the project installer script:

```sh
curl -fsSL https://raw.githubusercontent.com/Asterlenist/kk/main/install.py | sudo python3
```

The installer builds `kk` and places it in `/usr/bin/` or `/usr/local/bin/` depending on the script implementation. Piping a downloaded script to `sudo python3` runs that script with administrative privileges; only use this command if you trust the source.

## Commands

```text
kk system
kk pkg-inst <package> [package...]
kk pkg-upd
kk pkg-upg
kk pkg-rem <package> [package...]
kk pkg-pur <package> [package...]
kk flatpak-inst <app-id> [app-id...]
kk service-add <service> [service...]
kk service-rem <service> [service...]
kk update
```

The application also accepts longer aliases:

```text
kk package-install <package> [package...]
kk package-update
kk package-upgrade
kk package-remove <package> [package...]
kk package-purge <package> [package...]
kk flatpak-install <app-id> [app-id...]
```

- `system` prints the detected operating system, Linux distribution, and available package manager.
- `pkg-inst`/`package-install` install one or more packages using the detected native package manager.
- `pkg-upd`/`package-update` refresh package data. On pacman systems this also upgrades installed packages, as Arch Linux does not support a safe standalone package-index refresh.
- `pkg-upg`/`package-upgrade` upgrade installed packages. On pacman systems this performs a full system synchronization and upgrade.
- `pkg-rem`/`package-remove` remove packages.
- `pkg-pur`/`package-purge` remove packages and clean configuration/dependencies where the package manager supports it.
- `flatpak-inst` installs application IDs from the Flathub remote. Flatpak must already be installed and configured with Flathub.
- `service-add` enables and starts systemd services immediately.
- `service-rem` stops and disables systemd services. It does not delete unit files.
- `update` pulls the latest project version and installs the updated binary to `/usr/local/bin/kk`.

Package operations use the available APT, DNF/YUM, or pacman executable. Administrative privileges may be required; `kk` does not invoke `sudo` or elevate privileges automatically. Review package names and system changes before running install, remove, purge, or upgrade commands.

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

The detector recognizes `apt-get`/`apt`, `dnf`/`yum`, and `pacman`. On Linux distributions using other package managers, and on macOS, Windows, and FreeBSD, `kk system` still reports the operating system but does not select a package manager; native package commands return an error. Flatpak and systemd service commands are separate from native package-manager support.
