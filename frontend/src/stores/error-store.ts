import { defineStore } from 'pinia';
import { Notify } from 'quasar';

export const useErrorStore = defineStore('error', {
  state: () => ({
    error: undefined as undefined|Error
  }),
  actions: {
    set(err: Error) {
      console.error(err)
      this.error = err;
      Notify.create({
        message: `${err}`,
        color: 'negative'
      })
    },
    clear() {
      this.error = undefined;
    }
  },
});
