export default {
  namespaced: true,
  state: () => ({
    domain: '',
  }),
  mutations: {
    setDomain(state, value) {
      state.domain = value;
    },
  },
};
