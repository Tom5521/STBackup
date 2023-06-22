# Silly Tavern Backup and Cloud Upload

This is a source code file written in the Go programming language, which provides a backup and restore tool for SillyTavern. The program uses the `rsync` command to synchronize the application files between the local server and the remote server. It also uses the `rclone` tool to synchronize SillyTavern files with a cloud storage service.
## Requirements

- Go (if you want to compile it yourself to your system)
- rsync
- rclone

## Installation

1. Clone this repository into the sillytavern folder
2. Open a terminal and navigate to the folder that contains the source code file.
3. Run the following command to compile the program:

```bash
go build backup.go
```

4. Once compiled, you can use the program by running the `backup` binary file in the same folder as the source code file.

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

### Options
- `name-of-remote.txt`: specifies the name of the cloud storage service.
## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
