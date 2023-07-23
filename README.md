# STBackup

![GitHub last commit (branch)](https://img.shields.io/github/last-commit/Tom5521/STBackup/dev?&label=last%20dev%20commit)
![GitHub last commit (branch)](https://img.shields.io/github/last-commit/Tom5521/STBackup/main?&label=last%20main%20commit)
![GitHub](https://img.shields.io/github/license/Tom5521/STBackup)
![GitHub repo size](https://img.shields.io/github/repo-size/Tom5521/STBackup)
![GitHub release (with filter)](https://img.shields.io/github/v/release/Tom5521/STBackup?logo=go&label=lastest%20release)




This is a source code file written in the Go programming language, which provides a backup and restore tool for SillyTavern. The program uses the `rsync` command to synchronize the application files between the local server and the remote server. It also uses the `rclone` tool to synchronize SillyTavern files with a cloud storage service.
## Requirements

- ### Required
- rsync
- rclone (Optional installation)
- #### Optional
- Go (if you want to compile it yourself to your system)
- Tar (if you want use tarballs)
- unzip (If you do not want to install rclone)

## Installation
### Build Method
1. Clone this repository into the sillytavern folder
2. Open a terminal and navigate to the folder that contains the source code file.
3. Run the following command to compile the program:

```bash
go build .
```
1. Once compiled, you can use the program by running the `STBackup` binary file in the same folder as the source code file.
2. (Optional) You can make a `./[binary] link` to be able to run the script from the root of SillyTavern and not need to enter the binary folder. This process is done automatically in the script
### Script Method
Using the script. Below is how to use it
It is useful for those who do not like to compile and want the binary once and for all.
## Usage

The program is run from the command line and accepts various commands and options.


### Commands

- `make`: Creates necessary folders for backup.
- `save`: Saves files to the backup destination.
- `secure-save`: save without delete any file of Backup folder
- `save tar`: Saves files to a tarball in the backup destination.
- `restore`: Restores files from the backup destination.
- `secure-restore`: Restore without delete any file (Use it in case you have updated sillytavern and want to restore the data without breaking the repo and having to clone it again.)
- `restore tar`: Restores files from a tarball in the backup destination.
- `route <destination>`: Moves the backup folder to a new destination.
- `start`: Starts the SillyTavern application.
- `update ST`: Updates the STBackup application. (not work in ST versions under 1.9.1,check error code 29)
- `update me`: Updates the STBackup application and rebuilds if necessary.
- `ls`: Lists files in the remote backup destination.
- `upload`: Uploads files to the remote backup destination.
- `upload tar`: Uploads a tarball to the remote backup destination.
- `download`: Downloads files from the remote backup destination.
- `download tar`: Downloads a tarball from the remote backup destination.
- `init`: Initializes the SillyTavern application.
- `rebuild`: This is a high priority command. It will rebuild the program (if you have the source code at hand) rather than run any other function than the logs and change to the root directory. As soon as it finishes executing the program it will terminate with error code 0.
- `link`: Creates a link to the backup program in the SillyTavern root directory.
- `version`: Displays the version of STBackup.
- `remote`: Configures the rclone remote server.
- `cleanlog`: Clears the log file.
- `log`: Displays the content of the log file.
- `help`: Displays the help message.
- `printconf`:Print config.json file
- `test`: Only works in the dev branch. I use it to debug the code
- `resetconf`: It is used to delete all program settings. Does not delete backups
- `download-rclone`: Download and extract the rclone binary, it is useful if the rclone installed in your distro does not work correctly or if you do not want to install it.
- `Setloglevel`: is used to set the log level in app.log interactively.
### Configuration
The configuration is located in the config.json file.
its parameters are:
1. `remote`: This parameter determines the path to the remote rclone server.
2. `include-folders`: This parameter adds folders to be included in the backup.
3. `exclude-folders`: This parameter adds folders to exclude in the backup (you are free to use it if it takes me too long to update SillyTavern)
4. `local-rclone`: Is used to determine whether to use a local rclone binary or the one that comes installed with the system, its possible values are yes and no.
5. `log-level`: Determines the log level to be printed in app.log, there are 3 levels, minimum, medium and high which are 0,1,2 respectively.

### Log

The program logs all the actions in the `app.log` file.

# STBackup Installer

This script is an installer for the STBackup program. This is a simple script for those who do not want to compile.
## Usage

To use this script, follow these steps:

1. Copy the installation file to the SillyTavern folder.

2. Run the installation file with the following command:

```
bash install.sh <platform>
```

Where `<platform>` is the platform on which you want to install the program. You can use the following values:

- `arm`: if you want to install the program on an Android device using the Termux app.
- `x64`: if you want to install the program on a computer with x86-64 architecture.

3. If you want to install the program on a platform other than `arm` or `x64`, you can modify the script to add support for that platform.

## Functionality

The script works as follows:

- If a platform is specified as an argument and it is not `clone`, the script creates a folder called "STBackup" and downloads the latest version of the program corresponding to the specified platform using the GitHub API. Then, it renames the downloaded file to "backup" and executes it.
- If "clone" is specified as an argument, the script clones the GitHub repository and compiles the program using the `go build` command.
- If no arguments are specified, the script does nothing.

In any case, the program is  and the backup or restore process is started as appropriate.

## Notes

- This script requires the prior installation of the `curl` tool.
- To use this script on other platforms, it is necessary to modify the file to add support for the specific platform.
- This is for termux and linux friends. If you use windows i recommend you to learn how to use [Rclone](https://rclone.org/) and [SillyTavernSimpleLauncher](https://github.com/BlueprintCoding/SillyTavernSimpleLauncher).
- I don't plan to make a windows version... Unless I'm bored. We'll see what I feel like doing.
- *DO NOT MOVE THE BINARY OUTSIDE THE "STBackup" FOLDER.*
## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.


### At what point did this project go from being a personal bullshit to a moderately serious project?
