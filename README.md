# VSCode clean recently opened list

A simple script to clean your VSCode "recently opened" list. This is the list you get when you run the command "File: Open Recent". It includes recently opened folders, workspaces, files and can be used to quickly search and open them. There's a builtin command called "File: Clear Recently Opened" but that command empties the whole list whereas this script simplify removes the entries that are no longer present on your disk.

Usage:

```bash
go run main.go
```

By the way, I've also made a Pop Launcher plugin to have quick access to this list without having to open VSCode first. Here's the link if you're interested: [jeusto/pop-launcher-plugins](https://github.com/Jeusto/pop-launcher-plugins)
