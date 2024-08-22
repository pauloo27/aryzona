# ARYZONA

The one and only Discord Bot that is proudly Flamenguista.

> ⚠ I am mentally unstable, so I may `git rebase` on the master branch. Deal with it

## Features

- multi-language support (currently only english and portuguese)
- play Flamengo's matches live (using the `!radio globo-rj`)
- fetch live soccer scores (`!live`)
- random animals (`!dog`, `!cat`, `!fox`)
- play radios: lofi, rock, pisadinha, pop, cafe, swing and more (`!radio`)
- play youtube videos (`!play`)
- xkcd (`!xkcd latest`, `!xkcd random`)
- and a few more stuff

## How to add to your server

Aryzona is not public yet.

## How to self host 

_The bot is written in GoLang, so you need the Go compiler._

_To store data, the bot uses Postgres, so you need to have a Postgres server running._


The bot supports loading the config from a `config.yml` file or from environment variables.
You need to pick one.

If you want to use a YAML file, create a `config.yml` file and copy the content of the 
`config.example.yml` file into it.

If you want to use environment variables, you can set them manually or in a `.env` file.
See the `.env.example` file for an example.

To start the bot, run `make run` or `go run .`.

## License

<img src="https://i.imgur.com/AuQQfiB.png" alt="GPL Logo" height="100px" />

This project is licensed under [GNU General Public License v2.0](./LICENSE).

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GUwU General Public License for more details.
