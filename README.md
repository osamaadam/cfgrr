# cfgrr (Configurer)

`cfgrr` is a tool to help manage the myriad of config files usually found on your Linux setup. If you're familiar with [GNU stow](https://www.gnu.org/software/stow/), you'll find `cfgrr` to be an opinionated replacement for `stow` with a much more guided setup.

At the time of typing this, Jan 2023, `cfgrr` should be able to find all your config files, add them to some kind of directory structure, track them with git, and give you the option to push them to your private git remote, and resore them at will.

## Disclaimer

This tool was tested on Linux systems only, though it should also work with MacOS. It also confusingly works on Windows too (though I don't recommend using it there).

## Installation

### Download the binary

Download the latest binary for [Linux](https://github.com/osamaadam/cfgrr/releases/latest/download/cfgrr_linux), [MacOS](https://github.com/osamaadam/cfgrr/releases/latest/download/cfgrr_macos), or [Windows](https://github.com/osamaadam/cfgrr/releases/latest/download/cfgrr_windows.exe), and don't forget to [add the file to your `$PATH`](https://linuxize.com/post/how-to-add-directory-to-path-in-linux/).

**-or-**

### Install from source

> Note: You need to have [Go](https://golang.org/) installed.

First, clone the repository:

```sh
git clone https://github.com/osamaadam/cfgrr.git
```

Then run the install script:

```sh
make install
```

The binary would be installed at `${HOME}/go/bin/cfgrr`

## TL;DR

#### To backup:

```sh
cfgrr b ~/
```

Then choose the files you'd like to backup.

:sparkles: Poof! A backup of your files will be at `~/.config/cfgrr`.

#### To restore:

> (assuming you have the backup folder at `~/.config/cfgrr`)

```sh
cfgrr r -a
```

## Current progress

Currently, `cfgrr` is able to find config files matching certain patterns (currently `"**/.*"`, and `"**/*config*"`). It is also able to restore backed up files.

## Usage

#### Backup:

:warning: **WARNING** :warning: `backup` will copy the files to the backup directory, and replace them with symlinks to their equivalent in the backup directory.

This will backup all the config files found in the given directory.

```sh
cfgrr backup [root_path] [...files]
```

> :bell: You'll be prompted to choose the files you'd like to backup.

To skip the prompt, use the `--all` flag.

```sh
cfgrr backup ~/dotfiles/ -a
```

:mag: For more info, run `cfgrr backup --help`.

##### Examples:

```sh
cfgrr backup ~/.config
```

```sh
cfgrr b ~/.bashrc ~/.zshrc
```

This will backup all the config files found in `~/.config` matching the pattern `**/.*` or `**/*config*` (default patterns).

#### Restore:

:warning: **WARNING** :warning: `restore` will replace the files from the described paths (paths in cfgrrmap.yaml) with symlinks to their equivalent in the backup directory.

This will restore all the backed up files to their original locations.

```sh
cfgrr restore
```

:mag: For more info, run `cfgrr restore --help`.

##### Example:

```sh
cfgrr r -d ~/cfgrr/
```

This will prompt you to choose the files you'd like to restore from the backup directory `~/cfgrr/`.

If you'd like to restore all the files, you can use the `-a` flag.

```sh
cfgrr r -a
```

#### Set:

This is an interface to set the config values for `cfgrr`.

> The config values are stored in `~/.cfgrr.yaml` by default.

```sh
cfgrr set [key] [value]
```

:mag: For more info, run `cfgrr set --help`.

##### Example:

```sh
cfgrr set backup_dir ~/cfgrr
```

#### Unset:

This is an interface to unset the config values for `cfgrr`.

```sh
cfgrr unset [key]
```

:mag: For more info, run `cfgrr unset --help`.

##### Example:

```sh
cfgrr unset backup_dir
```

#### Setup:

This is an interface to setup `cfgrr` for the first time.

```sh
cfgrr setup
```

You'll then be prompted to choose the various config values for `cfgrr`.

:mag: For more info, run `cfgrr setup --help`.

#### Delete:

This subcommand allows the user to delete backed up files from the backup directory.

```sh
cfgrr delete [...files]
```

If the user were to use the `--restore` flag, the symlinks created previously would be replaced by the file from the backup directory. This is a safe way of undoing your backups.

```sh
cfgrr delete --restore [...files]
```

By default, the `-r` flag only replaces the files if they are symlinks to the target. To restore the files regardless if they are symlinks or not, use the `--force` flag.

```sh
cfgrr delete --restore --force [...files]
```

or more compactly, `cfgrr d -rf [...files]`.

In case no files were provided to the argument, the user will be prompted to choose the files to delete.

:mag: For more info, run `cfgrr delete --help`.
