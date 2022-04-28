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

## youtube playlist watching

One method of triggering `ytdlwatch` to download videos is by watching a
playlist for new videos on YouTube.

### setup

1. Create a project and an API key for accessing the Youtube API:
   https://console.cloud.google.com/apis/credentials

2. Enable access to YouTube API at
   https://console.developers.google.com/apis/api/youtube.googleapis.com/overview

3. Create a new playlist on Youtube. The playlist can be public or unlisted, but
   not private. In Playlist Settings, enable the setting "Add new videos to top
   of playlist"
   ([screenshot](https://user-images.githubusercontent.com/11027/162623093-046a8400-8438-4261-b2c5-e4517dc28be7.png)).
   Needed so that the program only ever has to read the first page of playlist's
   items, and API quota does not have to be used for deleting videos.

   Could the personal "Watch later" playlist be used? Unfortunately not, since
   it's not available through the official API.

4. Run `ytdlwatch` with the environment variables listed below set up.

See also the example systemd service:
[ytdlwatch.service.example][example-systemd-service]

## http api

The program can be made to listen for HTTP requests that include video URL to be
downloaded.

1. Define the `PORT` environment variable
2. Send a POST request to path `/download` with video URL as the body.

### curl example

```sh
curl -X POST http://localhost:8080/download -d'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

## plex integration (optional)

One use case for ytdlwatch is to download videos into a Plex library. That part
is doable by just creating a new library in Plex, and adding the directory at
`OUTPUT_PATH` as media source. What maybe remains is a need for mechanism to
delete watched videos. You can make ytdlwatch to delete watched videos
periodically by configuring the Plex related environment variables below.

## env vars

- `LOG_LEVEL`: Log level of the program. Defaults to `info`.
- `PORT`: The port http server will bind to.
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
- `YOUTUBE_DL_ARGS`: Parameters passed to the executable configured above. This
  can be used to override the default parameters used by ytdlwatch. If set,
  `OUTPUT_PATH` has no effect, and `--paths <output path>` has to be manually
  included in the args, if desired.
- `PLEX_API_TOKEN`: Plex API Token. This official support article has details on
  how to obtain it: [Finding an authentication token / X-Plex-Token
  ][plex-api-token]
- `PLEX_API_URL`: URL to Plex server. This might be something like
  `http://localhost:32400`.
- `PLEX_LIBRARY_KEY`: Key to the Plex library where YouTube videos are added to.
  The program will delete watched videos from this library. You can find the key
  by opening a library with browser in Plex UI and checking page URL for a query
  parameter like `source=7`.

## development

The project uses [`just`](https://github.com/casey/just) as a command runner (or
make alternative).

See `just -l` for recipes.

[yt-dlp]: https://github.com/yt-dlp/yt-dlp
[youtube-dl]: https://github.com/ytdl-org/youtube-dl
[example-systemd-service]:
  https://github.com/raine/ytdlwatch/blob/master/ytdlwatch.service.example
[plex-api-token]:
  https://support.plex.tv/articles/204059436-finding-an-authentication-token-x-plex-token/
