# Silly Tavern Backup and Cloud Upload

This is a source code file written in the Go programming language, which provides a backup and restore tool for SillyTavern. The program uses the `rsync` command to synchronize the application files between the local server and the remote server. It also uses the `rclone` tool to synchronize SillyTavern files with a cloud storage service.
## Requirements

- Go (if you want to compile it yourself to your system)
- rsync
- rclone

## Installation
### Build Method
1. Clone this repository into the sillytavern folder
2. Open a terminal and navigate to the folder that contains the source code file.
3. Run the following command to compile the program:

```bash
go build backup.go
```
1. Once compiled, you can use the program by running the `backup` binary file in the same folder as the source code file.
2. (Optional) You can make a `./backup link` to be able to run the script from the root of SillyTavern and not need to enter the binary folder. This process is done automatically in the script
### Script Method
Using the script. Below is how to use it
It is useful for those who do not like to compile and want the binary once and for all.
## Usage

The program is run from the command line and accepts various commands and options.

### Commands

- `make`: creates the backup folder in the application's main directory.
- `save`: performs a backup of the application files to the backup folder.
- `restore`: restores the backup files to the application folder.
- `route`: moves the backup folder to a different location.
- `start`: starts the application server.
- `update`: updates the application from a Git repository. Use `update me` to update the script and `update ST` to update the SillyTavern app.
- `ls`: lists the files in the cloud storage service.
- `upload`: uploads the application files to the cloud storage service.
- `download`: downloads the application files from the cloud storage service.
- `init`: initializes and configures the cloud storage service.
- `rebuild`: rebuilds the program binary file.
- `link`:Create a link in the root folder of SillyTavern so you don't have to go into the folder to make your backups. The input of the commands is the same, but now you can do it from the SillyTavern folder.

### Options
- `name-of-remote.txt`: specifies the name of the cloud storage service.


# SillyTavernBackup Installer

This script is an installer for the SillyTavernBackup program. This is a simple script for those who do not want to compile.
## Usage

To use this script, follow these steps:

1. Copy the installation file to the SillyTavern folder.

2. Run the installation file with the following command:

```
bash install.sh <platform>
```

Where `<platform>` is the platform on which you want to install the program. You can use the following values:

- `termux`: if you want to install the program on an Android device using the Termux app.
- `pc`: if you want to install the program on a computer with x86-64 architecture.

3. If you want to install the program on a platform other than `termux` or `pc`, you can modify the script to add support for that platform.

## Functionality

The script works as follows:

- If a platform is specified as an argument and it is not `clone`, the script creates a folder called "SillyTavernBackup" and downloads the latest version of the program corresponding to the specified platform using the GitHub API. Then, it renames the downloaded file to "backup" and executes it.
- If "clone" is specified as an argument, the script clones the GitHub repository and compiles the program using the `go build` command.
- If no arguments are specified, the script does nothing.

In any case, the program is executed and the backup or restore process is started as appropriate.

## Notes

- This script requires the prior installation of the `curl` tool.
- To use this script on other platforms, it is necessary to modify the file to add support for the specific platform.
- This is for termux and linux friends. If you use windows i recommend you to learn how to use [Rclone](https://rclone.org/) and [SillyTavernSimpleLauncher](https://github.com/BlueprintCoding/SillyTavernSimpleLauncher).
- I don't plan to make a windows version... Unless I'm bored. We'll see what I feel like doing.
- *DO NOT MOVE THE BINARY OUTSIDE THE "SillyTavernBackup" FOLDER.*
## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
