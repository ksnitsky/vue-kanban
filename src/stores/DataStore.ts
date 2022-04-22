import { ColumnType } from './../types/column.type';
import { defineStore } from 'pinia';

export const useDataStore = defineStore('DataStore', {
  state: () => ({
    columns: <ColumnType[]>[],
  }),
  
  getters: {
    nextColumnId: (state): number => state.columns.length,

    // If I don't use the api, it could be a big performance problem
    nextCardId: (state): number =>
      Math.max(...state.columns.flatMap(column =>
        column.cards.flatMap(card => card.id)
      )) + 1
  },

  actions: {
    createCol(title: string): void {
      this.columns.push({
        id: this.nextColumnId,
        title: title,
        cards: [],
      });
    },

    createCard(column: ColumnType, content: string): void {
      column.cards.push({ id: this.nextCardId, content });
    },
  }
});
