import { createI18n } from 'vue-i18n';
import en from './locales/en.json';
import zh from './locales/zh.json';
import ja from './locales/ja.json';

const messages = {
  en,
  zh,
  ja
};

const savedLocale = localStorage.getItem('locale');
const defaultLocale = savedLocale || 'en';

const i18n = createI18n({
  legacy: false, // you must set `false`, to use Composition API
  locale: defaultLocale, // set locale
  fallbackLocale: 'en', // set fallback locale
  messages, // set locale messages
  globalInjection: true,
});

export default i18n;
