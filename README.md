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
* **--no-search-known-hosts**

  omit adding hosts from well-known file locations.

### environment variables

environment variables related to the application
are prefixed with *SSH_SELECT_*:

* **SSH_SELECT_NO_SEARCH_KNOWN_HOSTS**

  if defined, has the same effect as
  **--no-search-known-hosts**.
* **SSH_SELECT_SSH_BINARY**

  same as **--ssh**
* **SSH_SELECT_KNOWN_HOSTS_FILE_$**

  environment variables with this prefix have the
  same effect as **--known-hosts**; the suffix has
  no effect on the result.

[known_hosts]: http://man.openbsd.org/sshd.8#SSH_KNOWN_HOSTS_FILE_FORMAT
