# SoTBot - Your humble SoT based Discord bot
SoTBot is a Sea of Thieves themed Discord bot written in Go (Golang) and makes heavy use of the 
fantastic [discordgo](https://github.com/bwmarrin/discordgo) library.

## Requirements
To run SoTBot you require a Discord bot token. You can create one in the
[Discord developer portal](https://discord.com/developers/applications)

Also the bot uses SQLite3 for it's database. Therefore the libsqlite3 library is required to be present on the
machine running the bot.

To build the bot from the sources, you need to have Go installed

## Installation

## Building from source
There is a `Makefile` included in the project. Just run `make` and your Bot binary will be built in
as `./bin/sotbot`

## Features
SoTBot is heavily influenced by the Eggdrop bots of the olden IRC days. Lots of its commands are SoT-based, but
there are couple of fun non SoT-related commands as well.

### Sea of Thieves related commands
#### Getting access to the API
Unfortunately the SoT API does not offer any kind of OAuth2 for authentication, hence we have to use
a kind of hackish way to get your access cookie from the Microsoft Live login. You can use Simon's
[SoT-RAT-Extractor](https://github.com/echox/sot-rat-extractor) (please read the notes in the project before using
it) to get a current cookie and store it in the database of the bot. Unfortunately the cookie is only valid for 
14 days, so you'll have to renew it every now and then.

#### Important security note
The RAT cookie gives full access to your SoT account page and the bot does not store the cookie in any 
encrypted way. Therefore, before storing your cookie in the bot's DB, please make sure that you know what
you are doing. Maybe at some time, RARE decides to offer an API which offers OAuth2, so we can allow the bot
having access to the API data without having to store/renew the cookie.

## Helpful commands

Converting a MP3 file to a bot-compatible DCA file:
You first need to install [ffmpeg](https://ffmpeg.org/) and [dca](https://github.com/bwmarrin/dca). Then run:
```shell
$ ffmpeg -i file.mp3 -f s16le pipe:1 | dca >./media/audio/file.dca
```

## Attribution
The sounds that are provided in this repository are by the following people:
* Airhorn by kneedrawp: https://www.youtube.com/watch?v=1SVe1D7er-U
* Fart sound by Paula: https://soundbible.com/1605-Blowing-A-Raspberry.html
* Angry pirate sound by Mike Koenig: https://soundbible.com/858-Angry-Pirate.html