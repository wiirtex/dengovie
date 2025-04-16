// @ts-nocheck

import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import {ElTable, ElTableColumn, ElSwitch, ElInputNumber, ElButton, ElCheckbox} from 'element-plus';
import CashSplitter from "../CashSplitter.vue";

describe('CashSplitter.vue', () => {
  it('рендеринг компонента и таблицы', () => {
    const wrapper = mount(CashSplitter, {
      global: {
        components: {
          ElTable,
          ElTableColumn,
          ElSwitch,
          ElInputNumber,
          ElButton,
        },
      },
    });

    // Проверяем, что компонент отрендерился
    expect(wrapper.exists()).toBe(true);

    // Проверяем, что таблица отрендерилась
    const table = wrapper.findComponent(ElTable);
    expect(table.exists()).toBe(true);

    // Проверяем, что все колонки отрендерились
    const columns = wrapper.findAllComponents(ElTableColumn);
    expect(columns.length).toBe(4); // 4 колонки: selection, name, telegram, debt
  });

  it('логика выбора строк в таблице', async () => {
    const wrapper = mount(CashSplitter, {
      global: {
        components: {
          ElTable,
          ElTableColumn,
          ElSwitch,
          ElInputNumber,
          ElButton,
        },
      },
    });

    // Получаем таблицу
    const table = wrapper.findComponent(ElTable);

    // Эмулируем выбор строки
    const firstRow = { name: 'Aleyna Kutzner', telegram: '@aleyna' };
    await table.vm.$emit('selection-change', [firstRow]);

    // Проверяем, что строка добавлена в multipleSelection
    expect(wrapper.vm.multipleSelection).toEqual([firstRow]);
  });

  it('логика вычисления долга (debt)', async () => {
    const wrapper = mount(CashSplitter, {
      global: {
        components: {
          ElTable,
          ElTableColumn,
          ElSwitch,
          ElInputNumber,
          ElButton,
        },
      },
    });

    // Устанавливаем сумму
    await wrapper.vm.$nextTick();
    wrapper.vm.num = 1000;

    // Выбираем одну строку
    const firstRow = { name: 'Aleyna Kutzner', telegram: '@aleyna' };
    await wrapper.vm.handleSelectionChange([firstRow]);

    // Проверяем, что долг вычислен правильно (1000 / 2 = 500)
    expect(wrapper.vm.debt).toBe(500);

    // Включаем переключатель "посчитать вместе со мной"
    wrapper.vm.isCountMe = false;

    // Проверяем, что долг вычислен правильно (1000 / 1 = 1000)
    expect(wrapper.vm.debt).toBe(1000);
  });

  it('логика переключателя isCountMe', async () => {
    const wrapper = mount(CashSplitter, {
      global: {
        components: {
          ElTable,
          ElTableColumn,
          ElSwitch,
          ElInputNumber,
          ElButton,
        },
      },
    });

    // Проверяем начальное состояние переключателя
    expect(wrapper.vm.isCountMe).toBe(true);

    // Эмулируем изменение состояния переключателя
    const switchComponent = wrapper.findComponent(ElSwitch);
    await switchComponent.setValue(false);

    // Проверяем, что состояние изменилось
    expect(wrapper.vm.isCountMe).toBe(false);
  });

  it('логика ввода суммы (num)', async () => {
    const wrapper = mount(CashSplitter, {
      global: {
        components: {
          ElTable,
          ElTableColumn,
          ElSwitch,
          ElInputNumber,
          ElButton,
        },
      },
    });

    // Проверяем начальное значение суммы
    expect(wrapper.vm.num).toBe(100);

    // Эмулируем изменение суммы
    const inputNumber = wrapper.findComponent(ElInputNumber);
    await inputNumber.setValue(500);

    // Проверяем, что значение изменилось
    expect(wrapper.vm.num).toBe(500);
  });
});