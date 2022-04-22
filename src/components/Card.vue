<script setup lang="ts">
import { CardType } from '../types/card.type';

const props = defineProps<{
  card: CardType,
}>();

</script>

<template>
  <li>
    <div>
      <p>
        {{ card.content }}
      </p>
    </div>
  </li>
</template>

<style scoped>
div {
  padding: .75rem;
  border-radius: .1875rem;
  background-color: #FFFFFF;
  box-shadow: 0 .0625rem .25rem rgba(9, 45, 66, 0.25);
}

li {
  cursor: move;
}
</style>
