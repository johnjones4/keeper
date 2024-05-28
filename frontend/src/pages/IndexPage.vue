<template>
  <q-page class="row items-center justify-evenly">
    <MdEditor 
      v-if="!!note.key"
      v-model="note.body"
      @on-change="deferredSaveNote"
      @on-save="saveNote"
      :preview="false"
      language="en-US"
    />
  </q-page>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import { useRoute } from 'vue-router'
import { useNotesStore } from 'src/stores/notes-store';
import { storeToRefs } from 'pinia';
import { MdEditor } from 'md-editor-v3';
import 'md-editor-v3/lib/style.css';
import { Notify } from 'quasar';

const route = useRoute();
const notesStore = useNotesStore();

const {note} = storeToRefs(notesStore);

let timer = undefined as NodeJS.Timeout|undefined;
let justLoaded = true;

const loadNote = async (key: string) => {
  justLoaded = true;
  await notesStore.selectNote(key);
}

const saveNote = async () => {
  if (justLoaded) {
    justLoaded = false;
    return;
  }
  await notesStore.saveNote();
  Notify.create({
    message: 'Note saved',
    color: 'positive'
  })
}

const deferredSaveNote = () => {
  if (timer) {
    clearTimeout(timer);
    timer = undefined;
  }
  timer = setTimeout(() => {
    saveNote();
  }, 1000);
}

defineOptions({
  name: 'IndexPage'
});

watch(() => route.params.key, async (newId) => {
  loadNote(newId as string)
});

if (route.params.key) {
  loadNote(route.params.key as string)
}

</script>

<style>
  .md-editor {
    height: 100%;
    flex-grow: 1;
  }
</style>