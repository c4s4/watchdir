Watchdir
========

Tool to watch directories and run commands when given file events are triggered (possible events are *CREATE*, *REMOVE*, *WRITE*, *RENAME* and *CHMOD*).

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
