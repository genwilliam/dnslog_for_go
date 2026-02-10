export interface ApiErrorState {
  status: number;
  type: string;
  message: string;
}

export default {
  namespaced: true,
  state: () => ({
    apiError: null as ApiErrorState | null,
  }),
  mutations: {
    setApiError(state, payload: ApiErrorState) {
      state.apiError = payload;
    },
    clearApiError(state) {
      state.apiError = null;
    },
  },
};
