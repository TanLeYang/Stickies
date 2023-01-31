# Stickies
Telegram stickers are great, but there's no easy way for multiple people to add to a sticker pack using the official Telegram Stickers bot.
Stickies Bot lets you create sticker packs that you can let others add to.

## How it works
1. Create a sticker pack using Stickies Bot
2. A unique sharing code will be generated
3. Share the code with others
4. Anyone with the code can then provide it to Stickies Bot to add stickers to the pack

You can create and share multiple sticker packs, each one will have its own unique sharing code.
You can also view all the sticker packs you have created and retrieve their sharing codes.

To remove stickers or delete the entire sticker pack, simply do so using the offical Telegram Stickers bot. Only the original creator can remove stickers and delete sticker packs.

## Deployment
The bot is not currently deployed anywhere.
You can deploy a version of this bot on your own in the environment of your choosing:
1. Register a new Telegram Bot
2. Set up a PostgreSQL database
3. Set the environment variables i.e. the bot name, and bot token, db credentials

Check the `.env.template` file for what environment variables are required.

A docker-compose yml file is also provided for running a postgres db locally.

## Limitations
Only image stickers are supported currently, no gifs.
