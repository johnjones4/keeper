<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
        />

        <q-toolbar-title>
          Notes
        </q-toolbar-title>

        <q-btn @click="newNote" color="positive">New</q-btn>
        <q-btn @click="logout" color="negative">Logout</q-btn>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
    >
      <notes-menu />
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import NotesMenu from 'components/NotesMenu.vue';
import { useQuasar } from 'quasar'
import { useNotesStore } from 'src/stores/notes-store';
import { useRouter } from 'vue-router';
import { useTokenStore } from 'src/stores/token-store';

const $q = useQuasar()
const notesStore = useNotesStore();
const router = useRouter();
const tokenStore = useTokenStore();

const leftDrawerOpen = ref(false);

const toggleLeftDrawer = () => {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}
const newNote = () => {
  $q.dialog({
    title: 'New Note',
    message: 'Note filepath',
    prompt: {
      model: '',
      type: 'text' // optional
    },
    cancel: true,
    persistent: true
  }).onOk(async data => {
    let path = data
    if (!path) {
      return
    }
    if (path[0] !== '/') {
      path = '/' + path
    }
    if (path.split('.').length < 2) {
      path = path + '.md'
    }

    const key = await notesStore.newNote(path)
    router.push(`/note/${key}`);
  })
}

const logout = () => {
  tokenStore.clear();
  notesStore.clear();
  router.push('/login')
}

</script>

<style>

  .q-page {
    justify-content: flex-start;
    align-items: flex-start;
    flex-direction: column;
  }

</style>