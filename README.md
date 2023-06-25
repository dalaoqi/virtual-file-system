# Virtual File System

The Virtual File System is a command-line application that simulates a file system with users, folders, and files. It provides functionality for creating users, managing folders, and manipulating files within the system.

## Features

The Virtual File System offers the following features:

- User Management: Register user in the file system.
- Folder Management: Create, delete, and list folders for each user.
- File Management: Create, delete, and list files within user folders.
- Error Handling: Proper error messages for invalid commands ,non-existent entities or duplicated folder/file creating.

## Requirements

![Golang](https://img.shields.io/badge/Golang-1.20.5-blue)  

## Installation

1. Clone the repository
```
git clone https://github.com/dalaoqi/virtual-file-system.git
```
2. Change to the project directory:
```
cd virtual-file-system
```
3. Build and RUN the project:
```
make
```

Note: `make help` displays the help message with available make commands

## Usage

The Virtual File System supports the following commands:

- `register [username]`: Create a new user with the specified username.
- `create-folder [username] [foldername] [description (optional)]`: Create a new folder for the specified user.
- `delete-folder [username] [foldername]`: Delete a folder and all files within it.
- `rename-folder [username] [foldername] [new-folder-name]`: Rename a folder with a new name.
- `list-folders [username] [--sort-name|--sort-created] [asc|desc]`: List all folders for the specified user, optionally sorting by name or creation date. The default sorting order is by name in ascending order.
- `create-file [username] [foldername] [filename] [description (optional)]`: create a file to the specified user's folder.
- `delete-file [username] [foldername] [filename]`: Delete a file from the specified user's folder.
- `list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]`: List all files in the specified user's folder, optionally sorting by name or creation date. The default sorting order is by name in ascending order.

Note: If any of the arguments `[username]`, `[foldername]`, or `[filename]` contain whitespace characters, you can enclose them in double quotes.

## Examples

Here are some example commands and their usage:

- Register a user: `register dalaoqi`, `register "dalaoqi is awesome"`
- Create a folder: `create-folder dalaoqi docs`, `create-folder dalaoqi "meeting docs"`
- Delete a folder: `delete-folder dalaoqi docs`
- List folders: `list-folders dalaoqi --sort-name asc`

- Create a file: `create-file dalaoqi docs test`, `create-file dalaoqi docs "test file"`
- Delete a file: `delete-file dalaoqi docs test`
- List files: `list-files dalaoqi docs --sort-created desc`
