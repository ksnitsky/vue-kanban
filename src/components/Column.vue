<script setup lang="ts">
import { ref } from 'vue';
import { ColumnType } from '../types/column.type';
import Card from './Card.vue';
import newBlock from './_newBlock.vue';

const props = defineProps<{
  column: ColumnType,
}>();


const createCard = (value: string): void => {
  props.column?.cards.push({
    content: value
  });
};

</script>

<template>
  <div class="column">
    <h2>
      {{ column.title }}
    </h2>

    <ul>
      <Card
        v-for="(card, index) in column.cards"
        :card="card"
        :key="index"
      ></Card>
    </ul>

    <newBlock
      @create-card="createCard"
    >
    </newBlock>
  </div>
</template>

<style >
.column {
  z-index: 0;
  width: 17rem;
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
  padding: .75rem;
  font-weight: bold;

  background-color: #DFE3E6;
  z-index: 100;
  overflow-wrap: break-word;
}

.column ul {
  width: 100%;
  list-style: none;
  padding: 0 .75rem .75rem;
  height: fit-content;
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
</style>
