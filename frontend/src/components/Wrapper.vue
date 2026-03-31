<script setup lang="ts">
import draggable from 'vuedraggable';
import newBlock from './_newBlock.vue';
import Column from './Column.vue';
import { useDataStore } from '../stores/DataStore';
import { computed } from 'vue';

const dataStore = useDataStore();
const store = dataStore;

const dragOptions = computed(() => {
  return {
    animation: 250,
    group: "column",
    disabled: false,
    ghostClass: "ghost"
  };
});

</script>

<template>
  <div class="wrapper">
    <draggable
      class="columns"
      handle=".handle"
      :list="store.columns"
      v-bind="dragOptions"
      itemKey="id"
    >
      <template #item="{ element }">
        <Column
          :column="element"
          :key="element.id"
        >
        </Column>
      </template>
    </draggable>

    <div class="column-item">
      <newBlock
        :is-col="true"
        placeholder="Добавить колонку"
        @create-col="store.createCol"
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

  background: url('../assets/background.webp') no-repeat center;
  background-size: cover;

  display: flex;
  flex-direction: row;
  gap: .75rem;

  overflow: hidden;
}

.columns {
  display: flex;
  flex-direction: row;
  gap: .75rem;
  flex-wrap: wrap;
}
</style>
