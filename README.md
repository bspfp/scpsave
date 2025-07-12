# scpsave

Cloud Save with SCP

- [Usage](#usage)
  - [Clone](#clone)
  - [Build](#build)
  - [Create Configuration File](#create-configuration-file)
  - [Rename Configuration File](#rename-configuration-file)
  - [Run](#run)
- [Configuration File Contents](#configuration-file-contents)

## Usage

### Clone

```powershell
cd X:\Working\Folder
git clone https://github.com/bspfp/scpsave.git
```

### Build

In `scpsave/cmd/scpsave`

```powershell
go build -o scpsave.exe ./...
```

Or `Run Task > build scpsave` in VSCode(open the project as a workspace file.)

### Create Configuration File

In `scpsave/cmd/scpsave`

```powershell
.\scpsave.exe -c
```

### Rename Configuration File

In `scpsave/cmd/scpsave`

After editing the file, rename it with the command below:

```powershell
Rename-Item -Path .\config.sample.yaml -NewName config.yaml
```

### Run

In `scpsave/cmd/scpsave`

```powershell
.\scpsave.exe
```

## Configuration File Contents

| Item                | Format             | Description                                                                                |
| ------------------- | ------------------ | ------------------------------------------------------------------------------------------ |
| server_address      | host:port          | SSH server address                                                                         |
| username            | username           | SSH username                                                                               |
| private_key_path    | file_path          | SSH user private key file path                                                             |
| remote_root         | absolute_path      | Absolute path to upload                                                                    |
| games               | game settings      | Game synchronization settings                                                              |
| games.name          | game name          | Must be unique                                                                             |
| games.local_dir     | save_file_folder   | Absolute path to save files                                                                |
| games.file_patterns | save_file_patterns | Be careful with backslashes and special character escaping                                 |
| games.program_name  | program_name       | (Optional) Absolute path to the game executable, or just the filename (e.g., filename.exe) |
