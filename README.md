# cfgrr (Configurer)

`cfgrr` is a tool to help manage the meriade of config files usually found on your Linux setup. If you're familiar with [GNU stow](https://www.gnu.org/software/stow/), you'll find `cfgrr` to be an opinionated replacement for `stow` with a much more guided setup.

At the time of typing this, Jan 2023, `cfgrr` should be able to find all your config files, add them to some kind of directory structure, track them with git, and give you the option to push them to your private git remote, and resore them at will.

## Installation

Download the binary from the [releases](https://github.com/osamaadam/cfgrr/releases) page, and put it somewhere in your `$PATH`.

## Current progress

Currently, `cfgrr` is able to find config files matching certain patterns (currently _"\*\*/.\*"_, and _"\*\*/\*config\*"_). It is also able to restore backed up files.

## Future plans

If I don't actually get hired soon, I'll have to finish this project. If I do get hired, consider this my will for whoever's brave enough to carry the torch. Here's a list of things I'd like to do:

- [ ] `backup` subcommand should initialize git in the backup directory.
- [ ] Create a `git` subcommand to allow committing, and pushing changes to a remote git repository.
- [ ] Write tests :eyes:
- [ ] Make a proper release.

## Usage

For now run `cfgrr -h` for help.

```
TODO: Write usage instructions here.
```

## TL;DR

To backup:

```sh
cfgrr b ~/
```

Then choose the files you'd like to backup.

:sparkles: Poof! A backup of your files will be at `~/.config/cfgrr`.

---

To restore:

> (assuming you have the backup folder at `~/.config/cfgrr`)

```sh
cfgrr r
```
