<script setup lang="ts">
import newBlock from './_newBlock.vue';
import Column from './Column.vue';
import { ref } from 'vue';
import { testdata } from '../assets/testdata';

const columns = ref(testdata);

const createCol = (newCol: string): void => {
  columns.value.push({
    id: columns.value.length,
    title: newCol,
    cards: [],
  });
};

</script>

<template>
  <div class="wrapper">
    <Column
      v-for="column of columns"
      :column="column"
      :key="column.id"
    >
    </Column>

    <div class="column">
      <newBlock
        :is-col="true"
        placeholder="Добавить колонку"
        @create-col="createCol"
      >
      </newBlock>
    </div>
  </div>
</template>

<style scoped>
.wrapper {
  width: 100%;
  height: 100%;
  padding: 1.25rem;

  background: url('../assets/background.png') no-repeat center;
  background-size: cover;

  display: flex;
  flex-direction: row;
  gap: .75rem;
}
</style>
