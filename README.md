![build](https://github.com/uhppoted/uhppoted-app-wild-apricot/workflows/build/badge.svg)

# uhppoted-app-wild-apricot

```cron```'able command line utility to transfer access control lists from a [Wild Apricot](https://www.wildapricot.com) 
organisational account to a set of UHPPOTE UTO311-L0x access access controllers. 

Supported operating systems:
- Linux
- MacOS
- Windows
- ARM7

## Releases

| *Version* | *Description*                                                                  |
| --------- | ------------------------------------------------------------------------------ |
|           |                                                                                |

## Installation

Executables for all the supported operating systems are packaged in the [releases](https://github.com/uhppoted/uhppoted-app-wild-apricots/releases). The provided archives contain the executables for all the operating systems - OS specific tarballs can be 
found in the [uhppoted](https://github.com/uhppoted/uhppoted/releases) releases.

Installation is straightforward - download the archive and extract it to a directory of your choice and then place the executable in a directory in your PATH. The `uhppoted-app-wild-apricot` utility requires the following additional 
information:

- `uhppoted.conf` configuration file
- Wild Apricot API key with read permission for  contact lists and member groups

### `uhppoted.conf`

`uhppoted.conf` is the communal configuration file shared by all the `uhppoted` project modules and is (or will 
eventually be) documented in [uhppoted](https://github.com/uhppoted/uhppoted). `uhppoted-app-wild-apricot` requires the 
_devices_ section to resolve non-local controller IP addresses and door to controller door identities.

A sample _[uhppoted.conf](https://github.com/uhppoted/uhppoted/blob/master/app-notes/wild-apricot/uhppoted.conf)_ file is included
in the `uhppoted` distribution.

### Building from source

Assuming you have `Go` and `make` installed:

```
git clone https://github.com/uhppoted/uhppoted-app-wild-apricot.git
cd uhppoted-app-wild-apricot
make build
```

If you prefer not to use `make`:
```
git clone https://github.com/uhppoted/uhppoted-app-wild-apricot.git
cd uhppoted-app-wild-apricot
mkdir bin
go build -o bin ./...
```

The above commands build the `uhppoted-app-wild-apricot` executable to the `bin` directory.

#### Dependencies

| *Dependency*                                                                 | *Description*                              |
| ---------------------------------------------------------------------------- | ------------------------------------------ |
| [com.github/uhppoted/uhppote-core](https://github.com/uhppoted/uhppote-core) | Device level API implementation            |
| [com.github/uhppoted/uhppoted-api](https://github.com/uhppoted/uhppoted-api) | Common API for external applications       |
| golang.org/x/lint/golint                                                     | Additional *lint* check for release builds |

## uhppoted-app-wild-apricot

Usage: ```uhppoted-app-wild-apricot [--debug] [--config <configuration file>] <command> [options]```

Supported commands:

- `help`
- `version`
- `get`
- `load-acl`
- `compare-acl`

### `help`

Displays a summary of the command usage and options.

Command line:

- ```uhppoted-app-wild-apricot help``` displays a short summary of the command and a list of the available commands

- ```uhppoted-app-wild-apricot help <command>``` displays the command specific information.

### `version`

Displays the current version of the command.

Command line:

```uhppoted-app-wild-apricot version```

### `get`

Fetches contact lists and membership groups from a Wild Apricot membership database stores it as a TSV file. Intended for use in
a `cron` task that routinely transfers information from the worksheet for scripts on the local host managing the access control
system. 

Command line:

```uhppoted-app-wild-apricot get``` 

```uhppoted-app-wild-apricot [--debug] get [--file <TSV>] [--workdir <dir>]```

```
  --file        File path for the destination TSV file. Defaults to <yyyy-mm-dd HHmmss>.tsv
  
  --workdir     Directory for working files, in particular the tokens, revisions, etc
                that provide access to Wild Apricot. Defaults to:
                - /var/uhppoted on Linux
                - /usr/local/var/com.github.uhppoted on MacOS
                - ./uhppoted on Microsoft Windows
  
  --debug       Displays verbose debugging information, in particular the communications
                with the UHPPOTE controllers
```

### `load-acl`

Fetches an ACL file from a Wild Apricot membership database and downloads it to the configured UHPPOTE controllers.
Intended for use in a `cron` task.

The command writes an operation summary to a _log_ file and a summary of changes to a _report_ .

Unless the `--force` option is specified, the command will only download and update changes since the last update. 

Command line:

```uhppoted-app-wild-apricot load-acl```

```uhppoted-app-wild-apricot [--debug] [--config <file>] load-acl [--force] [--strict] [--dry-run] [--workdir <dir>] [--no-log] [--no-report]```

```
  --force            Retrieves and updates the access control lists unconditionally.
  --strict           Fails with an error if the contacts and/or membership groups contains 
                     errors e.g. duplicate card numbers
  --dry-run          Executes the load-acl command but does not update the access
                     control lists on the controllers. Used primarily for testing 
                     scripts, crontab entries and debugging. 

  --workdir          Directory for working files, in particular the tokens, revisions,
                     etc, that provide access to Google Sheets. Defaults to:
                     - /var/uhppoted on Linux
                     - /usr/local/var/com.github.uhppoted on MacOS
                     - ./uhppoted on Microsoft Windows
  --no-log           Disables the creation of log entries on the 'log' worksheet
  
  --no-report        Disables the creation of report entries on the 'report' worksheet
    
  --config           File path to the uhppoted.conf file containing the access
                     controller configuration information. Defaults to:
                     - /etc/uhppoted/uhppoted.conf (Linux)
                     - /usr/local/etc/com.github.uhppoted/uhppoted.conf (MacOS)
                     - ./uhppoted.conf (Windows)

  --debug            Displays verbose debugging information, in particular the 
                     communications with the UHPPOTE controllers
```

### `compare-acl`

Fetches an ACL from a Wild Apricot membership database and compares it to the cards stored in the configured
access controllers. Intended for use in a `cron` task that routinely audits the controllers against an
authoritative source.

Command line:

```uhppoted-app-wild-apricot compare-acl ```

```uhppoted-app-wild-apricot [--debug] [--config <file>] compare-acl [--workdir <dir>]```
```
  --workdir       Directory for working files, in particular the tokens, revisions, etc, 
                  that provide access to Google Sheets. Defaults to:
                  - /var/uhppoted on Linux
                  - /usr/local/var/com.github.uhppoted on MacOS
                  - ./uhppoted on Microsoft Windows

  --config        File path to the uhppoted.conf file containing the access controller 
                  configuration information. Defaults to:
                  - /etc/uhppoted/uhppoted.conf (Linux)
                  - /usr/local/etc/com.github.uhppoted/uhppoted.conf (MacOS)
                  - ./uhppoted.conf (Windows)

  --debug         Displays verbose debugging information, in particular the 
                  communications with the UHPPOTE controllers

```
```
