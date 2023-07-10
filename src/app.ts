import axios from 'axios'
import { load } from 'cheerio'
import TelegramBot from 'node-telegram-bot-api'

require('dotenv').config();

interface Episod {
  title: string
  date: string
  relized: boolean
}

interface AnimeGoResp {
  episods: Episod[]
  image: string | null
  title: string | null
}

async function fetchAndParseAnimeData(url: string): Promise<AnimeGoResp> {

  const { data } = await axios.get(url, {
    headers: {
      'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3'
    }
  });

  const $ = load(data);
  const releasedEpisodesContainer = $('.released-episodes-container');

  const episods: Episod[] = []
  
  for (let i = 0; i < releasedEpisodesContainer.children().length; i++) {
    const cheld: any = releasedEpisodesContainer.children()[i];    

    if (cheld?.children[0]?.children[0]?.children[0]?.children[0]?.children[0]?.data !== undefined) {
      episods.push({
        title: cheld?.children[0]?.children[0]?.children[0]?.children[0]?.children[0]?.data,
        date: cheld?.children[0]?.children[0]?.children[2]?.children[0].attribs['data-label'],
        relized: cheld?.children[0]?.children[0]?.children[3]?.children[0] !== undefined
      })
    } 
  }

  const image = $('img');

  return { 
    episods: episods, 
    image: image[2]?.attribs?.src, 
    title: image[2]?.attribs?.title 
  }
}

const token = process.env.BOT_TOKEN;

if (token === undefined) {
  throw new Error("Not a telegram token ")
}

const bot = new TelegramBot(token, { polling: true });

const urlPattern = new RegExp('^https://animego.org/anime/.*$');

bot.on('message', async (msg: any) => {
  const chatId = msg.chat.id;

  if (urlPattern.test(msg.text)) {
    const { episods, image, title } = await fetchAndParseAnimeData(msg.text)
    
    if (episods[0].relized === true) {
      const message = `${title}\n\nУже вышла последняя ${episods[0].title} (${episods[0].date})`
      if (image && image !== null) {
        bot.sendPhoto(chatId, image, { caption: message });
      } else {
        bot.sendMessage(chatId, message);
      }
    } else {
      let lastEpisod = null;

      if (episods[0].relized === false && episods[1].relized === true) {
        lastEpisod = episods[0];
      } else if (episods[1].relized === false && episods[2].relized === true) {
        lastEpisod = episods[1];
      } else {
        lastEpisod = episods[2];
      }

      const message = `${title}\n\nАниме сохраненно, вы будете получать уведомления когда выйдут новые серии. ${lastEpisod.title} выйдет ${lastEpisod.date}.`;
      if (image && image !== null) {      
        bot.sendPhoto(chatId, image, { caption: message });
      } else {
        bot.sendMessage(chatId, message);
      }
    }
  } else {
    bot.sendMessage(chatId, 'Не похоже что это ссылка на animego. :(');
  }
});
