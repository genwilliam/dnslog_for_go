import { createStore } from 'vuex';
import text from './modules/text';
import persistPlugin from './persistPlugin';

const store = createStore({
  modules: {
    text,
  },
  plugins: [persistPlugin],
});

export default store;
