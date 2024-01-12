# ftp-uploader

It is a tool to upload files to your ftp server with `gitignore` like logic.

## Getting started

Install the latest build

```sh
curl https://google.com
```

Then init project with:

```sh
ftp-uploader init
```

It should generate your `ftp-uploader.yaml` with the following content:

```yaml
configs:
  default:
    root: .
    uploadRoot: my-relative-path/
    name: default
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

## Usage

If you set up everything properly, you can start your upload via the following command:

```sh
ftp-uploader upload -c default
```

If config is not passed it uploads all configs. Example:

```
ftp-uploader upload
```

You can also pass authentication credentials via terminal. Like in the example below

```sh
ftp-uploader upload --username MY_USER --pasword MY_PASSWORD --host MY_HOST --config example
```

Get more info with

```sh
ftp-uploader upload -h
```

