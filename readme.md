Watchdir
========

Tool to run a command when a file is modified in a given directory:

    watchdir /tmp 'echo "%s" >> /tmp/.liste'

This command will watch directory */tmp* and will run command 
`echo "%s" >> /tmp/.liste` when a new file is created. String *%s* is replaced
with created file name. If you want to write a *%* character in a command, you
can escape it with *%%*.
