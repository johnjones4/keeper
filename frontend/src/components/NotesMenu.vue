<template>
  <q-scroll-area class="fit">
    <notes-node :note="notes" />
  </q-scroll-area>
</template>

<script setup lang="ts">

import NotesNode from 'src/components/NotesNode.vue'
import { useNotesStore } from 'src/stores/notes-store';
import { storeToRefs } from 'pinia';
import { onMounted } from 'vue';

const notesStore = useNotesStore();
const {notes} = storeToRefs(notesStore);

onMounted(async () => {
  await notesStore.loadRootDirectory();
});

</script>