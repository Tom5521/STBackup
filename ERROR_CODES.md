# Official List of Error Codes

## Errcodes:
- **1:** It means that the corresponding and/or required data was not inserted or no option was specified in most cases.
- **2:** The backup destination is not specified in the route function.
- **3:** An application was not selected for update in the update function.
- **4:** There was an error creating the file for download.
- **5:** There was an error when requesting the file in the corresponding links.
- **6:** Error copying download request data to the corresponding file.
- **7:** The required JSON file was not found.
- **8:** The required variable could not be found in the JSON file.
- **9:** The value of the remote directory is "" or null (is solved by configuring it in config.json or by running `./backup remote`).
- **10:** Rclone was not found installed. You can fix it by downloading the binary using `./backup download-rclone`.
- **11:** No rsync was found installed. There is no option to use the binary locally because... It's RSYNC!!! It's not even 7MB!!! It's in all the damn repositories everywhere!!!
- **12:** There was an error reading the config.json file.
- **13:** Error decoding the config.json file.
- **14:** Error encoding the config.json file.
- **15:** Error writing to the config.json file.
- **16:** Error creating the binary file for download.
- **17:** Error making the download request.
- **18:** Error copying request data to file.
- **19:** No source code found when recompiling.
- **20:** No Go compiler found in the rebuild function.
- **21:** There was an unknown error in the rebuild function.
- **22:** Error when serializing the structure in the update json function
- **23:** Error oppening the json file
