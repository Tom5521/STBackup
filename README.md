## Silly Tavern Backup and Cloud Upload

This script is used to create a backup of the Silly Tavern program and upload it to the cloud. The backup includes various directories and files related to the program's data, such as chats, characters, themes, worlds, backgrounds, and configuration files.

### Prerequisites

- A sh shell (for the friends of termux or really old pc)
- If you are using termux or some other rare architecture you will need a go compiler.
- Rclone installed and properly configured for cloud storage access
- Rsync

### Usage

To use this script, follow these steps:

1. Clone in the Silly Tavern repository from GitHub.
2. Navigate to the project's root directory.
3. Open a terminal or command prompt.
4. Run the following command:

```
./backup [option] [destination]
```

Replace `[option]` with one of the following values:

- `start`: Starts the Silly Tavern server.
- `update`: Updates the Silly Tavern program.
- `make`: Creates the necessary backup directories.
- `save`: Saves a backup of the Silly Tavern program to the specified destination.
- `restore`: Restores a previously saved backup to the Silly Tavern program.

If you choose the `update` option, you can provide an additional parameter to specify what to update:

- `SillyTavern`: Updates the Silly Tavern program from the remote repository.
- `me`: Updates the current script.

### Starting the Server

When executing the script with the `start` option, the following steps are performed:

1. The script starts the Silly Tavern server by running `node server.js`.

### Updating the Program

When executing the script with the `update` option, the following steps are performed:

1. If no additional parameter is provided, the script displays an error message stating that nothing is selected.
2. If the parameter is `SillyTavern`, the script performs a `git pull` command in the current directory to update the Silly Tavern program from the remote repository.
3. If the parameter is `me`, the script navigates to the parent directory, performs a `git pull` command, and then returns to the original directory to update the current working directory.
If you choose the `save` or `restore` options, you need to provide a `[destination]` parameter specifying the location where the backup should be saved or restored.

### Backup Process

When executing the script with the `save` option, the following steps are performed:

1. The script creates the backup directories if they don't already exist.
2. Various directories and files related to the Silly Tavern program are copied to the backup directories, including chats, characters, themes, worlds, user avatars, backgrounds, group chats, groups, thumbnails, configuration files, and important JSON files.

### Restore Process

When executing the script with the `restore` option, the following steps are performed:

1. The script restores the backed-up directories and files to their original locations in the Silly Tavern program.

### Cloud Upload

The script also provides functionality to upload the backup to a cloud storage location using Rclone.

- To list the contents of the cloud storage directory, use the following command:

```
./backup ls
```

- To upload the backup to the cloud storage, use the following command:

```
./backup upload
```

- To download the backup from the cloud storage, use the following command:

```
./backup download
```

### Changing Backup Location

If you want to change the backup location, you can use the `route` option followed by the new destination path. This will move the entire backup directory to the specified location.

To change the backup location, run the following command:

```
./backup route [destination]
```

Replace `[destination]` with the desired path for the backup directory.

Note: The backup directory should already exist at the specified location.

### Important Notes

- Make sure that Rclone is properly configured with the desired cloud storage account before using the upload and download options.
- Always provide the correct option when running the script to avoid unexpected behavior or data loss.
- Backup files can be large, so ensure you have sufficient storage space and a stable internet connection when performing backup and cloud upload operations.
- Make sure that Node.js and Git are properly installed and configured before using this script.
- Always provide the correct option when running the script to avoid unexpected behavior.
- The `start` option assumes that the server script is named `server.js` in the current directory.
- The `update` option assumes that the Silly Tavern repository is already cloned in the current or parent directory.

For further assistance or information, please refer to the Silly Tavern documentation or contact the project maintainers.

### Unimportant notes
I made this program to synchronize my sillytavern files from my pc with those of my mobile phone with termux.
