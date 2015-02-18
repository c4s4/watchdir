Watchdir
========

- Project: <http://github.com/c4s4/watchdir>
- Downloads: <http://sweetohm.net/watchdir/>

Tool to watch directories and run commands when given file events are triggered (possible events are *CREATE*, *REMOVE*, *WRITE*, *RENAME* and *CHMOD*). Licensed under [GPL V3](http://www.gnu.org/licenses/gpl.html).

Installation
------------

Download binary archive at <http://sweetohm.net/watchdir/>, unzip it and copy the binary executable for your platform (named *watchdir-os-arch*) somewhere in yout *PATH* and rename it *watchdir*. This executable doesn't need any dependency or virtual machine to run.

There are binaries for following platforms:

- Linux 386, amd64 and arm.
- FreeBSD 386, amd64 and arm.
- NetBSD 386, amd64 and arm.
- OpenBSD 386 and amd64.
- Darwin 386 and amd64.
- Windows 386 and amd64.

There are no binaries for Plan9 because *fsnotify* library doesn't build on it.

Usage
-----

To run *watchdir* with *config.yml* configuration file, you should type:

    watchdir config.yml

If no configuration file is passed on command line, default ones are used if found:

- *~/.watchdir.yml*
- */etc/watchdir.yml*

Configuration
-------------

Configuration file is using *YAML* syntax and could look like:

    /tmp:
        CREATE: 'echo "%e: %f"'
        REMOVE: 'echo "%e: %f"'

This is a map with directories for keys. For each watched directory, a command is associated for given events.

In these commands, following replacements are made:

- *%f* is replaced with the absolute file name.
- *%e* is replaced with the name of the event (such as *CREATE* or *REMOVE*).
- *%%* is replaced with a single *%*.

Service
-------

To run this tool as a service, copy *watchdir.init* in */etc/init.d/* directory:

    sudo cp watchdir.init /etc/init.d/watchdir

You can then start and stop service with commands:

    sudo service watchdir start
    sudo service watchdir stop

Logs are written to file */var/log/watchdir.log*.

To make your service start at boot time:

    update-rc.d watchdir defaults

To remove it from boot sequence:

    update-rc.d -f watchdir remove

History
-------

- **1.0.0** (*2015-02-18*): First release.

*Enjoy!*
