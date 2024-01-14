# ftp-uploader

It is a tool to upload files to your ftp server with `gitignore` like logic.

I have a lot of projects that use shared hosting with ftp access. 
I wanted to automate the process of uploading them to my ftp server, with ability to ignore certain paths and certain files. 
So, I wrote this cli tool to solve my problems. 

This tool also makes transparent on what would be uploaded, and I believe would make CI/CD less cumbersome. 

Currently, I created binaries for linux and a docker image that can be used in your CI/CD pipelines.

## Install with docker

```
docker pull ghcr.io/ragueel/ftp-uploader:main
```

## Getting started

Install the latest build

### Linux

```shell
sudo curl -fsSL -o /usr/local/bin/ftp-uploader https://github.com/Ragueel/ftp-uploader/releases/latest/download/ftp-uploader-amd64
sudo chmod +x /usr/local/bin/ftp-uploader
```

Then init project with:

```shell
ftp-uploader init
```

It should generate your `ftp-uploader.yaml` with the following content:

```yaml
configs:
  default:
    root: . # local root directory where upload happens
    uploadRoot: my-relative-path/ # directory where files will be uploaded
    ignorePaths:
      - ftp-uploader.yaml
```

`ignorePaths` follows the same structure as ignore lines of any normal `.gitignore` file

You can also provide `ignoreFile` variable in the config. It will merge lines from the file with `ignorePaths`

By default, to authenticate the command uses the following environment variables

```
FTP_UPLOADER_USERNAME
FTP_UPLOADER_PASSWORD
FTP_UPLOADER_HOST
```

There are also some optional environment variables that might be useful
```
FTP_UPLOADER_CONNECTION_COUNT # controls how many parallel connections are created
ROOT_CONFIG_PATH # path to your ftp-upload.yaml file
```

## Usage

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
