<script setup lang="ts">
import { computed } from '@vue/reactivity';
import { ref } from 'vue';

const props = defineProps<{
  isCol: Boolean | true,
  placeholder: String | 'Добавить колонку',
}>();

const emit = defineEmits<{
  (e: 'createCol', value: string): void,
  (e: 'createCard', value: string): void,
}>();

const create = ref(false);
const newValue = ref('');

const toggle = () => create.value = !create.value;

const anotherOne = computed(() => {
  const arr: string[] = props.placeholder?.split(' ');
  return `${arr[0]} еще одну ${arr[1]}`;
});

const emitEvent = (value: string): void => {
  props.isCol ? emit('createCol', value) : emit('createCard', value);
  newValue.value = '';
}

</script>

<template>
  <div class="new-card">
    <div 
      class="standart"
      v-if="!create"
      @click="toggle"
    >
      <span class="new-card__plus">
      </span>

      <span class="new-card__text">
        {{ anotherOne }}
      </span>
    </div>

    <div class="create-form" v-if="create">
        <input
          v-if="isCol"
          type="text"
          placeholder="Введите название колонки"
          v-model="newValue"
        />

        <textarea
          v-else
          placeholder="Введите название карточки"
          v-model="newValue"
        >
        </textarea>

      <div class="create-form__bottom">
          <button
            class="create-form__create"
            @click="emitEvent(newValue);   toggle();"
          >
            {{ placeholder }}
          </button>

        <span
          class="create-form__discard"
          @click="toggle"
        >
        </span>
      </div>
    </div>
  </div>
</template>

<style>
.new-card {
  margin: .75rem;

  color: #6B808C;
  background-color: #DFE3E6;
}

.standart {
  align-items: center;
  display: flex;
  cursor: pointer;

  gap: .5rem;
  width: 100%;
  height: 100%;
}

.new-card__plus, .create-form__discard {
  width: .9375rem;
  height: .9375rem;

  background: url('../assets/plus.svg') no-repeat center; 
  background-size: cover;
}

.create-form__discard {
  -webkit-transform: rotate(45deg);
  transform: rotate(45deg);
  cursor: pointer;
}

.create-form textarea, input {
  width: 100%;
  resize: none;
  padding: .5rem .75rem;

  border: none;
  background-color: #FFFFFF;
  box-shadow: 0 .0625rem .25rem rgba(9, 45, 66, 0.25);
  border-radius: .1875rem;
}

.create-form__bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: .75rem;
}

.create-form__create {
  padding:  .25rem .5rem;
  font-weight: bold;

  color: #FFFFFF;
  background-color: #39C071;

  border: none;
  border-radius: .1875rem;

  cursor: pointer;
}

textarea, input {
  box-sizing: border-box;
  outline: none ;
}

.create-form__create:active {
  background-color: #188c49;
}

</style>
