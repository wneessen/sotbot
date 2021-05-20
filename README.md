# SoTbot
Your humble SoT bot

## Build
```shell
$ go build github.com/wneessen/bot/cmd/bot
$ ./bot -t <auth_token>
```

## Running without building
```shell
$ go run github.com/wneessen/bot/cmd/bot -t <auth_token>
```

Converting a MP3 to DCA:
```shell
$ ffmpeg -i file.mp3 -f s16le pipe:1 | dca >./media/audio/file.dca
```

## Attribution
Angry pirate sound: https://soundbible.com/858-Angry-Pirate.html