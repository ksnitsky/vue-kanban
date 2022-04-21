<script setup lang="ts">
import { ref } from 'vue';
import { ColumnType } from '../types/column.type';
import Card from './Card.vue';
import newBlock from './_newBlock.vue';

const props = defineProps<{
  column?: ColumnType,
}>();


const createCard = (value: string): void => {
  props.column?.cards.push({
    content: value
  });
};

</script>

<template>
  <div class="column">
    <div >
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
        :is-col="false"
        placeholder="Добавить карточку"
        @create-card="createCard"
      >
      </newBlock>
    </div>
  </div>
</template>

<style >
.column {
  z-index: 0;
  width: 17rem;
  height: fit-content;
  max-height: 100%;

  overflow: hidden;
  background-color: #DFE3E6;

  border-radius: .1875rem;
}

.column h2 {
  padding: .75rem .75rem 0;
  font-weight: bold;
}

.column ul {
  list-style: none;
  padding: .75rem;
  max-height: 25rem;

  overflow-y: auto;

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
