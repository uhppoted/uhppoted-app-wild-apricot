![build](https://github.com/uhppoted/uhppoted-app-wild-apricot/workflows/build/badge.svg)

# uhppoted-app-wild-apricot

```cron```'able command line utility to manage the access control lists of a set of UHPPOTE UTO311-L0x access access controllers from a [Wild Apricot](https://www.wildapricot.com) organisational account.

Supported operating systems:
- Linux
- MacOS
- Windows
- ARM7

## Releases

| *Version* | *Description*                                                                   |
| --------- | ------------------------------------------------------------------------------- |
| v0.7.3    | Maintenance release to update dependencies on `uhppote-core` and `uhppoted-lib` |
| v0.7.2    | Maintenance release to update dependencies on `uhppote-core` and `uhppoted-lib` |
| v0.7.1    | Maintenance release to update dependencies on `uhppote-core` and `uhppoted-lib` |
| v0.7.0    | Added support for time profiles from the extended API                           |
| v0.6.12   | Implemented prepended facility code for card numbers retrieved without          |
| v0.6.10   | Initial release                                                                 |

## Installation

Executables for all the supported operating systems are packaged in the [releases](https://github.com/uhppoted/uhppoted-app-wild-apricot/releases). The provided archives contain the executables for all the operating systems - operating system specific tarballs can be found in the [uhppoted](https://github.com/uhppoted/uhppoted/releases) releases.

Installation is straightforward - download the archive and extract it to a folder of your choice and then place the executable in a filder in your PATH. The `uhppoted-app-wild-apricot` utility requires the following additional information:

- `uhppoted.conf` configuration file
- a Wild Apricot account ID
- a Wild Apricot API key with read permission for  contact lists and member groups
- a [Grule](http://hyperjumptech.viewdocs.io/grule-rule-engine) rules file that defines the member access permissions

### `uhppoted.conf`

`uhppoted.conf` is the communal configuration file shared by all the `uhppoted` project modules and is (or will eventually be) documented in [uhppoted](https://github.com/uhppoted/uhppoted). `uhppoted-app-wild-apricot` requires the 
_devices_ section to resolve non-local controller IP addresses and door to controller door identities. 

It also uses the following additional configuration items:

| *Key* | *Default value* | *Description*                                             |
| ----- | --------------- | --------------------------------------------------------- |
| `wild-apricot.http.client-timeout` | 10s | Wild Apricot API request timeout           |
| `wild-apricot.http.retries`        | 3   | Number of times retry a failed API request | 
| `wild-apricot.http.retry-delay`    | 5s  | Interval between retries of a failed API request |
| `wild-apricot.facility-code`  | Facility code | Facility code prepended to card numbers that are 5 digits or less |
| `wild-apricot.fields.card-number`  | Card Number | Contact field name to use for card number  |
| `wild-apricot.display-order.groups` | _(alphabetic)_ | Optional output ordering for the member list groups | 
| `wild-apricot.display-order.doors`  | _(alphabetic)_ | Optional output ordering for the ACL doors |

A sample _[uhppoted.conf](https://github.com/uhppoted/uhppoted/blob/master/app-notes/wild-apricot/uhppoted.conf)_ file is included in the `uhppoted` distribution.

### `credentials.json`

A _credentials_ file should be a valid JSON file that contains the Wild Apricot account ID and API key e.g.:
```
  { 
    "account-id": 615252,
    "api-key": "8dhuwyeb7262jdufhde87bhbdehdes"
  }
```

### Access rules file

The _rules_ file is a text file containing the [Grule](http://hyperjumptech.viewdocs.io/grule-rule-engine) rules that define the member access e.g.:
```
// *** GRULES ***

rule Teacher "Grants a teacher access to common areas and Hogsmeade" {
     when
         member.HasGroup("Teacher")
     then
         record.Grant("Great Hall");
         record.Grant("Gryffindor");
         record.Grant("Hufflepuff");
         record.Grant("Ravenclaw");
         record.Grant("Slytherin");
         record.Grant("Hogsmeade");
         Retract("Teacher");
}

rule Staff "Grants ordinary staff access to common areas, Hogsmeade and kitchen" {
     when
         member.HasGroup(601422)
     then
         record.Grant("Great Hall");
         record.Grant("Gryffindor");
         record.Grant("Hufflepuff");
         record.Grant("Ravenclaw");
         record.Grant("Slytherin");
         record.Grant("Hogsmeade");
         record.Grant("Kitchen");
         Retract("Staff");
}

rule Gryffindor "Grants a Gryffindor student access to common areas and Gryffindor" {
     when
         member.HasGroup("Student") && member.HasGroup("Gryffindor")
     then
         record.Grant("Great Hall");
         record.Grant("Gryffindor");
         Retract("Gryffindor");
}

rule Pets "Grants a Pet access to the Kitchen at mealtimes" {
     when
         member.HasGroup("Pet")
     then
         record.Grant("Kitchen:100");
         Retract("Pets");
}

rule DeathEaters "Denies Hogwarts access to any known members of the Death Eaters" {
     when
         member.HasGroup("Death Eaters)
     then
         record.Revoke"Great Hall");
         record.Revoke("Gryffindor");
         record.Revoke("Hufflepuff");
         record.Revoke("Ravenclaw");
         record.Revoke("Slytherin");
         record.Revoke("Dungeon);
         record.Revoke("Kitchen");
         Retract("DeathEaters");
}

// *** GRULES ***
```
_Notes:_

1. The card associated with a member has access to a door if access has been `granted` and NOT `revoked`. If a card's door access has been both granted and revoked then it will not have access e.g. a _Slytherin_ student could have access to the _Great Hall_, unless of course he/she is a known associate of _Voldemort_.
2. To grant time based access to a door:
   ```
   record.Grant("Kitchen:100");
   ```
3. The grules file must have markers at the start and end of the file as a basic validity check for downloaded files e.g. Google Drive
   shares occasionally download as empty files without error (ref. https://github.com/uhppoted/uhppoted-app-wild-apricot/issues/2)

   - The first non-blank line should be 
   ```
   // *** GRULES ***
   ```
   
   - The last non-blank line should be:
   ```
   // *** END GRULES ***
   ```

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
go build -trimpath -o bin ./...
```

The above commands build the `uhppoted-app-wild-apricot` executable to the `bin` directory.

#### Dependencies

| *Dependency*                                                                 | *Description*                              |
| ---------------------------------------------------------------------------- | ------------------------------------------ |
| [com.github/uhppoted/uhppote-core](https://github.com/uhppoted/uhppote-core) | Device level API implementation            |
| [com.github/uhppoted/uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) | Shared application library                 |
| [github.com/hyperjumptech/grule-rule-engine](https://github.com/hyperjumptech/grule-rule-engine) | Grule rule engine for processing ACL rules |
| github.com/sirupsen/logrus                                                   | Indirect dependency from [grule-rule-engine](https://github.com/hyperjumptech/grule-rule-engine) |
| golang.org/x/sys                                                             | Library for Windows system calls |

## uhppoted-app-wild-apricot

Usage: ```uhppoted-app-wild-apricot [--debug] [--config <configuration file>] <command> [options]```

Supported commands:

- `help`
- `version`
- `get-members`
- `get-groups`
- `get-doors`
- `get-acl`
- `compare-acl`
- `load-acl`

### `help`

Displays a summary of the command usage and options.

Command line:

`uhppoted-app-wild-apricot help`

`uhppoted-app-wild-apricot help <command>`

### `version`

Displays the current version of the command.

Command line:

`uhppoted-app-wild-apricot version`

### `get-members`

Combines the contacts list and membership groups from a Wild Apricot membership database to create the table of members that is used to generate an access control list. The resulting member summary list can be optionally stored to a file. Intended for use in a `cron` task that routinely retrieves information from the Wild Apricot database for use by scripts on the local host. 

Command line:

```uhppoted-app-wild-apricot get-members --credentials <file>``` 

```uhppoted-app-wild-apricot [--debug] [--config <file>] get-members [--credentials <file>] [--workdir <dir>] [--file <file>]```

```
  --credentials <file> File path for the credentials file with the Wild Apricot account ID and API key. 
                       Defaults to <config dir>/.wild-apricot/credentials.json

  --workdir      Directory for working files, in particular the tokens, revisions, etc. Defaults to:
                 - /var/uhppoted on Linux
                 - /usr/local/var/com.github.uhppoted on MacOS
                 - ./uhppoted on Microsoft Windows

  --file <file> Optional file path for the destination TSV file. Displays a formatted member list on console if not provided. If the file has a .tsv extension, the output is formatted as a TSV file.
    
  --config      File path to the uhppoted.conf file containing the access
                controller configuration information. Defaults to:
                - /etc/uhppoted/uhppoted.conf (Linux)
                - /usr/local/etc/com.github.uhppoted/uhppoted.conf (MacOS)
                - ./uhppoted.conf (Windows)

  --debug       Displays verbose debugging information, in particular the communications
                with the UHPPOTE controllers
```

### `get-groups`

Retrieves the membership groups from a Wild Apricot membership database to displays as a table (or optionally stores it to a file). Intended as a convenience to assist when creating the rules that convert a member list into an access control list. 

Command line:

```uhppoted-app-wild-apricot get-groups --credentials <file>``` 

```uhppoted-app-wild-apricot [--debug] [--config <file>] get-groups [--credentials <file>] [--workdir <dir>] [--file <file>]```

```
  --credentials <file> File path for the credentials file with the Wild Apricot account ID and API key. 
                       Defaults to <config dir>/.wild-apricot/credentials.json

  --workdir      Directory for working files, in particular the tokens, revisions, etc. Defaults to:
                 - /var/uhppoted on Linux
                 - /usr/local/var/com.github.uhppoted on MacOS
                 - ./uhppoted on Microsoft Windows

  --file <file> Optional file path to which to write the output. Displays a formatted groups list on console if not provided. If the file has a .tsv extension, the output is formatted as a TSV file.
    
  --config      File path to the uhppoted.conf file containing the access
                controller configuration information. Defaults to:
                - /etc/uhppoted/uhppoted.conf (Linux)
                - /usr/local/etc/com.github.uhppoted/uhppoted.conf (MacOS)
                - ./uhppoted.conf (Windows)

  --debug       Displays verbose debugging information, in particular the communications
                with the UHPPOTE controllers
```

### `get-doors`

Extracts the list of doors from the `uhppoted.conf` configuration file. Intended as a convenience to assist when creating the rules that convert a member list into an access control list. 

Command line:

```uhppoted-app-wild-apricot get-doors``` 

```uhppoted-app-wild-apricot [--debug] [--config <file>] get-doors [--workdir <dir>] [--file <file>]```

```
  --workdir      Directory for working files, in particular the tokens, revisions, etc. Defaults to:
                 - /var/uhppoted on Linux
                 - /usr/local/var/com.github.uhppoted on MacOS
                 - ./uhppoted on Microsoft Windows

  --file <file> Optional file path to which to write the output. Displays a formatted door list on console if not provided. If the file has a .tsv extension, the output is formatted as a TSV file.
    
  --config      File path to the uhppoted.conf file containing the access
                controller configuration information. Defaults to:
                - /etc/uhppoted/uhppoted.conf (Linux)
                - /usr/local/etc/com.github.uhppoted/uhppoted.conf (MacOS)
                - ./uhppoted.conf (Windows)

  --debug       Displays verbose debugging information, in particular the communications
                with the UHPPOTE controllers
```

### `get-acl`

Retrieves the contacts list and membership groups from a Wild Apricot membership database and applies the access rules to create an access control list for the configured set of access controllers. The ACL can optionally be stored to a file for use by other scripts. 

Command line:

```uhppoted-app-wild-apricot get-acl --credentials <file> --rules <uri>``` 

```uhppoted-app-wild-apricot [--debug] [--config <file>] get-acl --credentials <file> --rules <uri> [--workdir <dir>] [--file <TSV>]```

```
  --credentials <file> File path for the credentials file with the Wild Apricot account ID and API key.

  --rules <uri>  URI for the Grule file that defines the rules used to grant or
                 revoke access (assumes a local file if the URI does not start with
                 http://, https:// or file://).  
                 Note that for rules files stored on Google Drive, the URI should be
                 of the form:
                 https://drive.google.com/uc?export=download&id=<file ID>

  --workdir      Directory for working files, in particular the tokens, revisions, etc. Defaults to:
                 - /var/uhppoted on Linux
                 - /usr/local/var/com.github.uhppoted on MacOS
                 - ./uhppoted on Microsoft Windows

  --file <file> File path for the optional output file. Displays the ACL on the console
                if not provided. Formats the output as TSV if the provided file has a
                .tsv extension.
    
  --config      File path to the uhppoted.conf file containing the access
                controller configuration information. Defaults to:
                - /etc/uhppoted/uhppoted.conf (Linux)
                - /usr/local/etc/com.github.uhppoted/uhppoted.conf (MacOS)
                - ./uhppoted.conf (Windows)

  --debug       Displays verbose debugging information, in particular the communications
                with the UHPPOTE controllers
```

### `compare-acl`

Retrieves the contacts list and membership groups from a Wild Apricot membership database and applies the access rules to create an access control list for the configured set of access controllers and compares the resulting ACL with the access control lists stored on the configured controllers. Intended for use in a `cron` task that routinely audits the controllers against an authoritative source.

Command line:

```uhppoted-app-wild-apricot compare-acl --credentials <file> --rules <uri>``` 

```uhppoted-app-wild-apricot [--debug] [--config <file>] compare-acl [--credentials <file>] [--rules <uri>] [--strict] [--summary] [--workdir <dir>] [--report <file>]```

```
  --credentials <file> File path for the credentials file with the Wild Apricot account ID and API key.

  --rules <uri>  URI for the Grule file that defines the rules used to grant or
                 revoke access (assumes a local file if the URI does not start with
                 http://, https:// or file://).  
                 Note that for rules files stored on Google Drive, the URI should be
                 of the form:
                 https://drive.google.com/uc?export=download&id=<file ID>

  --strict       Fails with an error if the contacts and/or membership groups contains  
                 errors e.g. duplicate card numbers

  --summary      Reports only a summary of the comparison. Defaults to false.

  --workdir      Directory for working files, in particular the tokens, revisions, etc. Defaults to:
                 - /var/uhppoted on Linux
                 - /usr/local/var/com.github.uhppoted on MacOS
                 - ./uhppoted on Microsoft Windows

  --report <file> File path for the optional output file. Displays the compare report on
                  the console if not provided. Formats the output as TSV if the provided
                  file has a .tsv extension.
    
  --config      File path to the uhppoted.conf file containing the access
                controller configuration information. Defaults to:
                - /etc/uhppoted/uhppoted.conf (Linux)
                - /usr/local/etc/com.github.uhppoted/uhppoted.conf (MacOS)
                - ./uhppoted.conf (Windows)

  --debug       Displays verbose debugging information, in particular the communications
                with the UHPPOTE controllers
```

### `load-acl`

Retrieves the contacts list and membership groups from a Wild Apricot membership database and applies the access rules to create an access control list that is downloaded to the configured set of access controllers.. Intended for use in a `cron` task that routinely updates the controllers from an authoritative source.

The command writes an operation summary to a _log_ file and a summary of changes to a _report_ .

Unless the `--force` option is specified, the command will only download and update changes since the last update. 

Command line:

```uhppoted-app-wild-apricot load-acl```

```uhppoted-app-wild-apricot [--debug] [--config <file>] load-acl [--credentials <file>] [--rules <uri>] [--force] [--strict] [--dry-run] [--workdir <dir>] [--log <file>]```

```
  --credentials <file> File path for the credentials file with the Wild Apricot account ID and API key.

  --rules <uri>  URI for the Grule file that defines the rules used to grant or
                 revoke access (assumes a local file if the URI does not start with
                 http://, https:// or file://).  
                 Note that for rules files stored on Google Drive, the URI should be
                 of the form:
                 https://drive.google.com/uc?export=download&id=<file ID>

  --force        Retrieves and updates the access control lists unconditionally.
  --strict       Fails with an error if the contacts and/or membership groups contains  
                 errors e.g. duplicate card numbers

  --dry-run      Executes the load-acl command but does not update the access
                 control lists on the controllers. Used primarily for testing 
                 scripts, crontab entries and debugging. 

  --workdir      Directory for working files (in particular the version information  
                 for the last update). Defaults to:
                     - /var/uhppoted on Linux
                     - /usr/local/var/com.github.uhppoted on MacOS
                     - ./uhppoted on Microsoft Windows

  --log <file>  Optional output file for a summary of the load operation. Formatted as
                headerless TSV if the file has a .tsv extension. worksheet
  
  --report <file> Optional output file for a detailed report of the load operation. 
                  Formatted as headerless TSV if the file has a .tsv extension. 
  
  --config      File path to the uhppoted.conf file containing the access
                controller configuration information. Defaults to:
                - /etc/uhppoted/uhppoted.conf (Linux)
                - /usr/local/etc/com.github.uhppoted/uhppoted.conf (MacOS)
                - ./uhppoted.conf (Windows)

  --debug       Displays verbose debugging information, in particular the 
                communications with the UHPPOTE controllers
```

