import { ColumnType } from './../types/column.type';
import { defineStore } from 'pinia';

export const useDataStore = defineStore('DataStore', {
  state: () => ({
    columns: <ColumnType[]>[],
  }),

  getters: {
    nextId: (state): number => state.columns.length,
  },

  actions: {
    createCol(title: string): void {
      this.columns.push({
        id: this.nextId,
        title: title,
        cards: [],
      });
    },

    createCard(column: ColumnType, content: string): void {
      column.cards.push({ content });
    },

    // getColumnById(id: number): number {
    //   return this.columns.findIndex(
    //     (column) => column.id === id
    //   );
    // }
  }
});
