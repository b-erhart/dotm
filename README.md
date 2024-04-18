<h1 align="center">dotm</h1>
<h3 align="center">A simple dotfiles manager.</h3>

`dotm` is a command-line utility for managing [dotfiles](https://wiki.archlinux.org/title/Dotfiles) accross multiple devices and operating systems. While numerous tools exist for this purpose, they tend to be feature-heavy, requiring non-trivial configuration, or lack cross-platform operability. `dotm` addresses these issues by offering a straightforward, pragmatic, and cross-platform solution that can be set up within a few minutes.

> **Note:** This is just a little experiment. Please use with caution if you want to try it.

---

## Usage
### How it Works
The fundamental concept behind `dotm` is to maintain a Git repository containing all dotfiles. Config files, written in TOML format, define mappings between repository files/directories and their corresponding locations on systems. These mapped files are synchronized via the `dotm` CLI tool.

#### Entry Format
A single line entry is created for every file/directory in the repository that should be mapped to a system location. 
```toml
"<repository_location>" = "<system_location>"
```
The `<repository_location>` should be relative to the repository's root directory.

#### Example Repository
Repository folder structure:
```
DOTFILES REPOSITORY
│
├── bash/
│   └── .bashrc
├── nvim/
│   └── ...
├── ssh/
│   └── ...
└── windows-config.toml
```

Config file `windows-config.toml`:
```toml
"nvim" = "${LOCALAPPDATA}/nvim"
"bash/.bashrc" = "${HOMEPATH}/.bashrc"
"ssh" = "C:/Users/User/.ssh"
```
> Note: using `~` for the home directory is currently not supported. Please use the `${HOMEPATH}` environment variable on Windows and `${HOME}` on Unix-based systems instead.

#### Commands
`dotm` provides two commands using config files: `dotm fetch <config_file>` copies specified entries from the system to the repository. `dotm distribute <config_file>` performs the opposite function, copying all entries from the repository to the system. Refer to `dotm -h`, `dotm fetch -h`, and `dotm distribute -h` for further information.

> Warning: Existing files/directories targeted by the copies will be fully overwritten. This means any existing data at these locations will be erased.

### Setting up a Dotfiles Repository
1. Obtain the `dotm` binary from the [releases page](https://github.com/b-erhart/dotm/releases) or install it via `go install github.com/b-erhart/dotm`.
2. Create an empty Git repository at any locaion on your system.
3. Create a config file defining mappings between your system files/directories and the repository.
4. Run `dotm fetch <config_file>` at the repo root to copy the files from their respective locations to the repository.
5. Commit the files to the repository.

### Deploying Dotfiles to a Machine
1. Obtain `dotm` for your system and clone your dotfiles repository.
2. Run `dotm distribute <config_file>` at the repo root to copy files/directories from the repository to specified locations. Create a config for the system if one doesn't exist yet.

### Changing dotfiles
There are two possible approaches for changing dotfiles:
1. Modify the files at their system locations. Then run `dotm fetch` at the repository root to pull the changes into the repository.
2. Modify the files inside the repository. Then run `dotm distribute` to push the changes to the respective system files.