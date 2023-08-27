# Anime Schedule Bot

## Operation principle

- You send the bot a link to an anime (mandatory from the animego.org website).
- The bot mentions the release time of the episode (in Japan).
- The bot checks every hour, according to the cron schedule, if a new episode of the anime has been released (with voice acting or subtitles).

### Supported websites

#### English-language sites

- [4anime.is](https://4anime.is/)

#### Russian-language sites

- [animego.org](https://animego.org/)
- [amedia.online](https://amedia.online/)
- [animevost.org](https://animevost.org/)

## Local setup

- make env file (see .env.example)
- run script from Makefile

```sh
make dev 
```

## Production setup 

- run script from Makefile

```sh
make start
```

> it runs docker containers