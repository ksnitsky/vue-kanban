<script setup lang="ts">
import Card from './Card.vue';
import draggable from 'vuedraggable/src/vuedraggable';
import newBlock from './_newBlock.vue';
import { ColumnType } from '../types/column.type';
import { useDataStore } from '../stores/DataStore';
import { ref, computed } from 'vue';

const props = defineProps<{
  column: ColumnType,
}>();

const drag = ref(false);

const createCard = (content: string): void =>
  useDataStore().createCard(props.column, content);

const dragOptions = computed(() => {
  return {
    animation: 250,
    group: "columns",
    disabled: false,
    ghostClass: "ghost"
  };
});

</script>

<template>
  <div class="column">
    <h2>
      {{ column.title }}
    </h2>

    <draggable
      tag="ul"
      @start="drag = true"
      @end="drag = false"
      v-model="column.cards"
      v-bind="dragOptions"
      item-key="id"
    >
      <template #item="{ element }">
        <Card
          :card="element"
          :key="element.id"
          :class="
            element.fixed ? 'fa fa-anchor' : 'glyphicon glyphicon-pushpin'
          "
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
.column {
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

.column h2 {
  padding: .75rem .75rem .375rem;
  font-weight: bold;

  background-color: #DFE3E6;
  z-index: 100;
  overflow-wrap: break-word;
}

.column ul {
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

.column ul {
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.column ul::-webkit-scrollbar {
  display: none;
}

.ghost {
  opacity: 0.5;
  background: #c8ebfb;
}

.flip-list-move {
  transition: transform 0.5s;
}
.no-move {
  transition: transform 0s;
}
</style>
