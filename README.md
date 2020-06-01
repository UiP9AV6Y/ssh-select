# ssh-select

![Build Status](https://github.com/UiP9AV6Y/ssh-select/workflows/Build%2FTest/badge.svg)

SSH-select is an auto-completion wrapper for SSH.
The completion engine reads hosts from various
known locations for SSH [known_hosts][] files, but
can be instructed to search other locations as well.

SSH-select is a one-off wrapper, i.e. once a host
has been selected, its process is replaced with that
of an SSH client using the *exec* system call. There
are no daemons or other lingering processes involved.

## examples

`ssh-select`

without any arguments, known locations will be searched
for host inputs. the resulting suggestions vary between
execution environments.

`ssh-select --no-search --known-hosts /tmp/ssh_known_hosts`

populates the suggestion database with entries from the
provided [known_hosts][] file. if the file does not
exist, no autocompletion is possible as other sources
have been disabled.

`ssh-select --ssh /opt/bin/rsh`

passes the collected arguments to the provided binary
instead of the regular `ssh` binary.

`ssh-select --no-search --hosts /etc/hosts`

disables the automated import from well known host
sources, but adds the */etc/hosts* file back into
the selection.

`ssh-select --zone /var/named/example.com.zone`

reads hostnames from the provided [zone file][].

`ssh-select --zone '/var/lib/bind/**/db.*'`

use hostnames from each [zone file][] the glob pattern
resolves to.

> **NOTE**: when using globs, wrap the argument in
> quotes to prevent the shell from expanding it,
> which would cause only a single file to be
> used as input and the rest as `ssh` arguments

`ssh-select -i /tmp/id_rsa

arguments not supported by `ssh-select` will be
forwarded to `ssh` as-is.

## configuration

SSH-select can change its behaviour base on
commandline arguments and environment variables.
commandline arguments which are not
understood/supported are passed on to the spawned
SSH client unchanged.

### commandline arguments

commandline arguments take precedence over
environment variables:

* **--ssh**

  path to the SSH client to invoke once a host has
  been selected. if not specified, `ssh` is searched
  for in the users *PATH*.
* **--known-hosts**

  path to a [known_hosts][] file. can be specified
  multiple times. duplicate hosts are filtered.
  supports [glob patterns][].
* **--hosts**

  path to a [hosts file][]. can be specified
  multiple times. duplicate hosts are filtered.
  supports [glob patterns][].
* **--zone**

  path to a domain name [zone file][]. can be
  specified multiple times. duplicate hosts are
  filtered. supports [glob patterns][].

  currently only the following resource records
  are processed:

  *PTR*, *A*, *AAAA*, *CNAME*, *MX*, *NS*
* **--no-search**

  omit adding hosts from well-known file locations.
* **--quiet**

  omit the search summary upon start.

### environment variables

environment variables related to the application
are prefixed with *SSH_SELECT_*:

* **SSH_SELECT_NO_SEARCH**

  if defined, has the same effect as **--no-search**.
* **SSH_SELECT_QUIET**

  if defined, has the same effect as **--quiet**.
* **SSH_SELECT_SSH_BINARY**

  same as **--ssh**
* **SSH_SELECT_KNOWN_HOSTS_FILE_$**

  environment variables with this prefix have the
  same effect as **--known-hosts**; the suffix has
  no effect on the result.
* **SSH_SELECT_HOSTS_FILE_$**

  environment variables with this prefix have the
  same effect as **--hosts**; the suffix has
  no effect on the result.
* **SSH_SELECT_ZONE_FILE_$**

  environment variables with this prefix have the
  same effect as **--zone**; the suffix has
  no effect on the result.

[known_hosts]: http://man.openbsd.org/sshd.8#SSH_KNOWN_HOSTS_FILE_FORMAT
[hosts file]: http://www.tldp.org/LDP/solrhe/Securing-Optimizing-Linux-RH-Edition-v1.3/chap9sec95.html
[zone file]: https://tools.ietf.org/html/rfc1035
[glob patterns]: https://github.com/bmatcuk/doublestar#patterns
