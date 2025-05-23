# cfgrr (Configurer)

![Demo](./assets/cfgrr-demo.gif)

`cfgrr` is a tool to help manage the myriad of config files usually found on your Linux setup. If you're familiar with [GNU stow](https://www.gnu.org/software/stow/), you'll find `cfgrr` to be an opinionated replacement for `stow` with a much more guided setup.

At the time of typing this, Jan 2023, `cfgrr` should be able to find all your config files, add them to some kind of directory structure, track them with git, and give you the option to push them to your private git remote, and restore them at will.

## Disclaimer

This tool was tested on Linux systems only, though it should also work with macOS. It also confusingly works on Windows too (though I don't recommend using it there).

## Installation

### Using `go install`

```sh
go install github.com/osamaadam/cfgrr@latest
```

### Download the binary

Download the latest binary for your respective platform from the [latest release](https://github.com/osamaadam/cfgrr/releases/latest)

**-or-**

### Install from source

> Note: You need to have [Go](https://golang.org/) installed.

First, clone the repository:

```sh
git clone https://github.com/osamaadam/cfgrr.git
```

Then run the installation script:

```sh
make install
```

The binary would be installed at `${HOME}/go/bin/cfgrr`

## TL;DR

#### To back up:

```sh
cfgrr b ~/
```

Then choose the files you'd like to back up.

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

This will back up all the config files found in the given directory.

```sh
cfgrr backup [root_path] [...files]
```

> :bell: You'll be prompted to choose the files you'd like to back up.

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

This will back up all the config files found in `~/.config` matching the pattern `**/.*` or `**/*config*` (default patterns).

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

This is an interface to set up `cfgrr` for the first time.

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

In case no files were provided to the argument, the user will be prompted to choose the files to delete.

:mag: For more info, run `cfgrr delete --help`.

#### Replicate:

Creates a browsable replica of the files. Could be useful if the user may want to share the files in a manner that's human-readable.

```sh
cfgrr replicate
```

By default, this creates the replica at `.config/cfgrr/home/` (assuming the user is using default configurations) or generally at `BACKUP_DIR/home/`.

To generate the replica at a different location:

```sh
cfgrr replicate ~/path/to/replica
```

> The path must be absolute, otherwise the tool would assume the path is relative to the backup directory.

You know the drill by this point, to skip the prompt, use `--all` flag.

```sh
cfgrr replicate --all
```

:mag: For more info, run `cfgrr replicate --help`.

#### Push:

This subcommand allows the user to push the backed up files to a remote git repository.

```sh
cfgrr push
```

By default, this will push the changes to the remote `origin` on the current branch.

To push to a different remote or branch:

```sh
cfgrr push origin backup-branch
```

:mag: For more info, run `cfgrr push --help`.

#### Clone:

This subcommand allows the user to clone the backed up files from a remote git repository.
If a repository already exists, the latest changes will be pulled instead.

```sh
cfgrr clone <remote_url>
```

By default, this will clone the branch `master`. To clone a different branch:

```sh
cfgrr clone git@github.com:osamaadam/cfgrr.git --branch backup-branch
```

:mag: For more info, run `cfgrr clone --help`.

## Configuration Details

### MapFile Format Support

As of the latest version, `cfgrr` now supports both YAML and JSON formats for the map file:

- **YAML** (default): Files end with `.yaml` or `.yml`

  ```sh
  cfgrr set map_file cfgrrmap.yaml
  ```

- **JSON**: Files end with `.json`
  ```sh
  cfgrr set map_file cfgrrmap.json
  ```

The mapfile type is automatically determined by the file extension. When creating a new mapfile, you can specify the path with your preferred extension:

```sh
cfgrr backup ~/.bashrc -m /path/to/your/cfgrrmap.json
```

You can switch between formats at any time by changing the mapfile path in your configuration.
