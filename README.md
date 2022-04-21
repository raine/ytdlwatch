# ytdlwatch

A long-running Go program that watches a Youtube playlist for new videos, and
downloads them using [`yt-dlp`][yt-dlp] or other preferred tool.

## install

The Go Toolchain is required.

The following will install `ytdlwatch` to `$GOPATH/bin`:

```sh
go install github.com/raine/ytdlwatch@latest
```

To build for another platform, clone the repository and run go build like so:

```sh
env GOOS=target-OS GOARCH=target-architecture go build .
```

## how to use

1. Create a project and an API key for accessing the Youtube API:
   https://console.cloud.google.com/apis/credentials

2. Create a new playlist on Youtube. The playlist can be public or unlisted, but
   not private. In Playlist Settings, enable the setting "Add new videos to top
   of playlist"
   ([screenshot](https://user-images.githubusercontent.com/11027/162623093-046a8400-8438-4261-b2c5-e4517dc28be7.png)).
   Needed so that the program only ever has to read the first page of playlist's
   items, and API quota does not have to be used for deleting videos.

   Could the personal "Watch later" playlist be used? Unfortunately not, since
   it's not available through the official API.

3. Run `ytdlwatch` with the environment variables listed below set up.

See also the example systemd service:
[ytdlwatch.service.example][example-systemd-service]

## env vars

- `YOUTUBE_API_KEY`: Youtube API key. Needed for accessing the given Youtube
  playlist. Create one here (and a new project, if necessary):
  https://console.cloud.google.com/apis/credentials **required**
- `YOUTUBE_PLAYLIST_ID`: An ID for the Youtube playlist that is monitored by the
  program. You can see this in browser's page URL when viewing a playlist on
  Youtube. A long string that begins with characters `PL`. **required**
- `OUTPUT_PATH`: A path to a directory where videos should be downloaded to.
  **required**
- `YOUTUBE_DL_PATH`: Path to executable that is used to download videos from
  Youtube. Defaults to [`yt-dlp`][yt-dlp], which is a fork of
  [`youtube-dl`][youtube-dl].
- `LOG_LEVEL`: Log level of the program. Defaults to `info`. Set to `debug` to
  see all output from whatever program is used to download videos.

## development

The project uses [`just`](https://github.com/casey/just) as a command runner (or
make alternative).

See `just -l` for recipes.

[yt-dlp]: https://github.com/yt-dlp/yt-dlp
[youtube-dl]: https://github.com/ytdl-org/youtube-dl
[example-systemd-service]:
  https://github.com/raine/ytdlwatch/blob/master/ytdlwatch.service.example
