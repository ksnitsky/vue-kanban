<script setup lang="ts">
import Card from './Card.vue';
import draggable from 'vuedraggable';
import newBlock from './_newBlock.vue';
import { ColumnType } from '../types/column.type';
import { useDataStore } from '../stores/DataStore';
import { computed } from 'vue';

const props = defineProps<{
  column: ColumnType,
}>();

const createCard = (content: string): void =>
  useDataStore().createCard(props.column, content);

const dragOptions = computed(() => {
  return {
    animation: 250,
    group: "cards",
    disabled: false,
    ghostClass: "ghost"
  };
});

</script>

<template>
  <div class="column-item">
    <h2 class="handle">
      {{ column.title }}
    </h2>

    <draggable
      tag="ul"

      v-model="column.cards"
      v-bind="dragOptions"
      itemKey="id"
    >
      <template #item="{ element }">
        <Card
          :card="element"
          :key="element.id"
        ></Card>
      </template>
    </draggable>

    <newBlock
      @createCard="createCard"
    >
    </newBlock>
  </div>
</template>

<style >
.column-item {
  z-index: 0;
  width: 18rem;
  height: fit-content;
  max-height: 100%;

  overflow-y: hidden;
  background-color: #DFE3E6;

  border-radius: .1875rem;

  position: relative;
  display: flex;
  flex-direction: column;
}

.column-item h2 {
  padding: .75rem .75rem .375rem;
  font-weight: bold;

  background-color: #DFE3E6;
  z-index: 100;
  overflow-wrap: break-word;
}

.column-item ul {
  width: 100%;
  list-style: none;
  padding: .375rem .75rem .375rem;
  max-height: 100vh;

  overflow-y: auto;
  scroll-behavior: smooth;

  display: flex;
  flex-direction: column;
  gap: .75rem;
}

.column-item ul {
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.column-item ul::-webkit-scrollbar {
  display: none;
}

.ghost {
  opacity: 0.5;
  background-color: #C9D2D9 !important;
}

.flip-list-move {
  transition: transform 0.5s;
}

.no-move {
  transition: transform 0s;
}
</style>
