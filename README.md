# Anime Sentry

## Operation principle

- You send the bot a link to an anime.
- The bot mentions the release time of the episode.
- The bot checks every half hour, according to the cron schedule, whether a new anime episode has been released (with voice acting or subtitles).
- You can opt out of notifications for a specific anime, and the bot will stop notifying you about the release of new episodes.

## Supported websites

- [animego.me](https://animego.me/)

![prewiew image](./prewiew.jpg)

## Supported Languages

The bot currently supports the following languages:

- En (English) 🇬🇧
- Ru (Russian) 🇷🇺

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