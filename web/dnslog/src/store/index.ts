import { createStore } from 'vuex';
import text from './modules/text';
import system from './modules/system';
import apiKey from './modules/apiKey';
import runtimeConfig from './modules/config';
import persistPlugin from './persistPlugin';

const store = createStore({
  modules: {
    text,
    system,
    apiKey,
    runtimeConfig,
  },
  plugins: [persistPlugin],
});

export default store;
