import { defineStore } from 'pinia';
import { postToken } from 'src/client/services.gen'
import { useErrorStore } from './error-store';
import { Message, TokenResponse } from 'src/client';

const key = 'token'

export const useTokenStore = defineStore('token', {
  state: () => ({
    haveToken: !!localStorage.getItem(key)
  }),
  getters: {
    authorization: () => {
      return `Bearer ${localStorage.getItem(key) || ''}`;
    },
  },
  actions: {
    clear() {
      localStorage.clear();
      this.haveToken = false;
    },
    async login(password: string) {
      try {
        const response = await postToken({
          requestBody: {
            password, 
          },
        });
        if ((response as Message).ok !== undefined) {
          throw new Error((response as Message).message)
        }
        localStorage.setItem(key, (response as TokenResponse).token);
        this.haveToken = true;
      } catch (err) {
        useErrorStore().set(new Error(`${err}`));
      }
    }
  },
});
