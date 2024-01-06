# ftp-uploader

It is a tool to upload files to your ftp server with some conditions.

## To get started
Install the latest build

```sh
curl https://google.com
```

Then init project with: 
```sh
ftp-uploader init
```

It should generate your `ftp-uploader.config.yaml` with the following content:

```yaml
configs:
  example:
    rootDir: .
    uploadRootDir: /your_folder
    ignore:
      - some_folder/
      - *.txt
```

Ignore follows the same structure as in any normal `.gitignore` file

By default to authenticate the command uses the following environment variables

```
FTP_UPLOADER_USERNAME
FTP_UPLOADER_PASSWORD
FTP_UPLOADER_HOST
```

Then you can start your upload via the following command:

```sh
ftp-uploader --config example upload
```

If config is not passed it uploads all configs

```
ftp-uploader upload
```

You can also pass authentication credentials via terminal. Like in the example below

```sh
ftp-uploader --username MY_USER --config example upload
```
