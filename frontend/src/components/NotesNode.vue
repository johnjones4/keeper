<template>
  <q-btn 
    v-if="note.name !== ''" 
    @click="() => expandNode(note.fullpath)"
    size="sm"
    unelevated
    padding="2px"
    align="left"
  >
    {{ note.name }}
  </q-btn>
  <ul>
    <li 
      v-for="dir in note.dirs" 
      :key="dir.fullpath" 
    >
      <notes-node :note="dir" />
    </li>
    <li 
      v-for="n in note.notes" 
      :key="n.key" 
      dense
    >
      <q-btn 
        :to="`/note/${n.key}`"
        size="sm"
        unelevated
        padding="2px"
        align="left"
      >
        {{ n.name }}
      </q-btn>
    </li>
  </ul>
</template>

<script setup lang="ts">
import { DirNode, useNotesStore } from 'src/stores/notes-store'

const notesStore = useNotesStore();

const expandNode = async (path: string) => {
  await notesStore.loadDirectory(path)
}

defineProps<{
  note: DirNode
}>()

</script>

<style>
  .q-item {
    flex-direction: column;
    flex-grow: 1;
  }

  ul {
    padding-inline-start: 20px;
  }
</style>