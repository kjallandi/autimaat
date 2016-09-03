## autimaat

This program is an IRC bot, specifically written for a private IRC channel.


## Install

    $ go get github.com/monkeybird/autimaat


## Usage

First, create a new profile directory and configuration file:

	$ autimaat -new /path/to/profile

Edit the newly created configuration file to your liking:

	$ nano /path/to/profile/profile.cfg

Relaunch the bot to use the new profile:

	$ autimaat /path/to/profile

In order to have the bot automatically re-launch after shutdown, an external
supervisor like systemd is required. The bot will create a PID file at
`/path/to/profile/app.pid`, in case the supervisor requires it.

When dealing with systemd, the bot may have to be forked at least once,
after it has been launched. Otherwise, systemd will keep killing it and
re-launching it in a never ending loop. Forking the bot is done through
the following command:

	$ kill -s USR1 `pidof autimaat`

This tells the bot to fork itself, while passing along any existing connections.
The old process then shuts itself down. This mechanism allows the bot to be binary-
patched, without downtime.


### Weather API

The `weather` module provides bindings for some APIs at
https://www.wunderground.com/weather/api/

This service requires the registration of a free account in order to get a
valid API key. The API key you receive should be assigned to the
`WeatherApiKey` field in the bot profile.


### Youtube API

The `url` module uses the `YouTube Data API v3` to fetch playback durations
for videos being linked in a channel. This API requires the registration of
a Google Developer API key at: [https://console.developers.google.com/apis](https://console.developers.google.com/apis)

The API key you receive should be assigned to the `YoutubeApiKey` field in
the bot profile.


## Versioning

The bot version is made up of 3 numbers:

* Major version: This number only changes if the bot itself changes in a way
  that makes it incompatible with previous versions. This does not include
  modules implementing commands.
* Minor version: This number changes whenever one of the module APIs change,
  or commands are added/removed.
* Revision: This is the build number. It is a current unix timestamp, which
  is updated whenever the bot is recompiled. This happenes whenever any kind
  of change occurs in any of the code. Including bug fixes. This number is
  updated through a go build flag. E.g.: `go install -ldflags "-X app.VersionRevision=12345"`



## license

Unless otherwise noted, the contents of this project are subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.
