# ftp-uploader

[![Release](https://img.shields.io/github/release/ragueel/ftp-uploader.svg)](https://github.com/ragueel/ftp-uploader/releases)
![GitHub CI](https://github.com/ragueel/ftp-uploader/actions/workflows/go.yml/badge.svg)

Tired of the manual grind when it comes to uploading projects to a shared hosting with FTP access? Enter 

`ftp-uploader`

a robust command-line tool engineered to simplify your workflow and automate the process of deploying projects to your FTP server.

## Key Features:

- **Seamless Automation:** Bid farewell to manual uploads. Effortlessly deploy your projects with a single command. 
- **Selective Uploads:** Tailor your uploads by leveraging Gitignore logic. Easily ignore specific paths and files, putting you in control of your project uploads.
- **Time-Efficient:** Save precious time on repetitive FTP uploads, allowing you to concentrate on refining your projects.

## Installation

Install the latest build

### Linux

```shell
sudo curl -fsSL -o /usr/local/bin/ftp-uploader https://github.com/Ragueel/ftp-uploader/releases/latest/download/ftp-uploader-linux-amd64
sudo chmod +x /usr/local/bin/ftp-uploader
```

### Docker

```
docker pull ghcr.io/ragueel/ftp-uploader:main
```

## Creating config

Init project with:

```shell
ftp-uploader init
```

It should generate your `ftp-uploader.yaml` with the following content:

```yaml
configs:
  default:
    root: . # local root directory, in which files you want to upload lie
    uploadRoot: my-relative-path/ # directory where files will be uploaded
    ignorePaths:
      - ftp-uploader.yaml
```

*note: You can have as many configs as you want!*

`ignorePaths` follows the same structure as ignore lines of any normal `.gitignore` file

You can also provide `ignoreFile` variable in the config. It will merge lines from the file with `ignorePaths`

You can configure the behaviour of the tool via the following environment variables

|name                          | description                                                |
|-------------------------------------------------------------------------------------------|
|FTP_UPLOADER_CONNECTION_COUNT | controls how many parallel connections are created         |
|ROOT_CONFIG_PATH              | path to your ftp-upload.yaml file                          |
|FTP_UPLOADER_USERNAME         | username                                                   |
|FTP_UPLOADER_PASSWORD         | password                                                   |
|FTP_UPLOADER_HOST             | host (port should be included)                             |

```

## Uploading to ftp server

If you set up everything properly, you can start your upload via the following command:

```shell
ftp-uploader upload -c default
```

If config is not passed it uploads all configs. Example:

```
ftp-uploader upload
```

You can also pass authentication credentials via terminal. Like in the example below

```shell
ftp-uploader upload --username MY_USER --pasword MY_PASSWORD --host MY_HOST --config example
```

It is also possible to use multiple connections to speed up the process of uploading. To do that pass `-t` flag with integer. 
The tool will create connections equal to that amount, and will use it to upload your files.

```shell
ftp-uploader upload -c default -t 10 # will use 10 connections to upload files
```


Get more info with

```shell
ftp-uploader -h
```

## Note:
This tool was created to simplify the FTP uploading process for projects hosted on shared servers. Your feedback and contributions will make it even better!
