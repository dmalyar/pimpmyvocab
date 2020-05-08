# PimpMyVocab

Telegram bot which helps learn new words.

At the moment supports only English-Russian dictionary.

The running instance: @PimpMyVocab_bot

## How to run it

### Common steps:
- Create a telegram bot and acquire a token.
- Add next commands to a bot:
    - /repeat
    - /quiz
    - /list
    - /clear
    - /help
- Acquire a yandex.dictionary token.
- Choose your way: dockerized app&db (docker-compose), dockerized app or non-dockerized app

### Dockerized app&db:
- Fill .env file
- Fill config.yaml file
- Run `docker-compose up` from project's root folder

### Dockerized app:
- Run postgres DB instance
- Fill config.yaml file
- Uncomment last line in the Dockerfile
- Run `docker build --tag pmv:1.0.0 .` and then `docker run --name pmv pmv:1.0.0`

### Non-dockerized app:
- Run postgres DB instance
- Fill config.yaml and copy it to `$HOME/.pimpmyvocab` folder
- Build/Install and run it just like any other go app

## TODO
- More languages
- More ways to remove words from a dictionary
- Multiple dictionaries per user
- All messages from resource files
- Cache
- More tests 
