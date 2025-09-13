import { describe, it, expect } from 'vitest' // или jest
import { mount } from '@vue/test-utils'
import DebtsIncomeList from '../DebtsIncomeList.vue'

describe('DebtsIncomeList.vue', () => {
  const mockDebts = [
    {
      another_user_id: 1,
      another_user_name: 'Иван Иванов',
      amount: 1500
    },
    {
      another_user_id: 2,
      another_user_name: 'Петр Петров',
      amount: 2500
    },
    {
      another_user_id: 3,
      another_user_name: 'Мария Сидорова',
      amount: 3500
    }
  ]

  it('отображает заголовок "Вам должны"', () => {
    const wrapper = mount(DebtsIncomeList, {
      props: {
        debts: []
      }
    })

    expect(wrapper.find('h1').text()).toBe('Вам должны')
    expect(wrapper.find('h1').classes()).toContain('text-green-500')
  })

  it('отображает правильное количество элементов долгов', () => {
    const wrapper = mount(DebtsIncomeList, {
      props: {
        debts: mockDebts
      }
    })

    const debtItems = wrapper.findAll('.border.rounded')
    expect(debtItems).toHaveLength(mockDebts.length)
  })

  it('корректно отображает информацию о каждом долге', () => {
    const wrapper = mount(DebtsIncomeList, {
      props: {
        debts: mockDebts
      }
    })

    const debtItems = wrapper.findAll('.border.rounded')

    mockDebts.forEach((debt, index) => {
      const item = debtItems[index]
      expect(item.text()).toContain(debt.another_user_name)
      expect(item.text()).toContain(`${debt.amount}₽`)
    })
  })

  it('отображает суммы долгов зеленым цветом', () => {
    const wrapper = mount(DebtsIncomeList, {
      props: {
        debts: mockDebts
      }
    })

    const amountElements = wrapper.findAll('.text-green-500')
    expect(amountElements).toHaveLength(mockDebts.length + 1) // +1 для заголовка

    mockDebts.forEach((debt, index) => {
      // Пропускаем первый элемент (заголовок)
      const amountElement = amountElements[index + 1]
      expect(amountElement.text()).toBe(`${debt.amount}₽`)
    })
  })

  it('отображает правильные CSS классы для контейнера', () => {
    const wrapper = mount(DebtsIncomeList, {
      props: {
        debts: mockDebts
      }
    })

    const container = wrapper.find('.overflow-y-auto')
    expect(container.classes()).toContain('max-h-[650px]')
    expect(container.classes()).toContain('h-full')
    expect(container.classes()).toContain('flex')
    expect(container.classes()).toContain('flex-col')
    expect(container.classes()).toContain('gap-1')
  })

  it('работает с пустым массивом долгов', () => {
    const wrapper = mount(DebtsIncomeList, {
      props: {
        debts: []
      }
    })

    const debtItems = wrapper.findAll('.border.rounded')
    expect(debtItems).toHaveLength(0)
  })

  it('корректно принимает и отображает props', () => {
    const wrapper = mount(DebtsIncomeList, {
      props: {
        debts: mockDebts
      }
    })

    expect(wrapper.props().debts).toEqual(mockDebts)
    expect(wrapper.vm.debts).toEqual(mockDebts)
  })

  it('имеет правильную структуру каждого элемента долга', () => {
    const wrapper = mount(DebtsIncomeList, {
      props: {
        debts: [mockDebts[0]]
      }
    })

    const debtItem = wrapper.find('.border.rounded')
    expect(debtItem.classes()).toContain('flex')
    expect(debtItem.classes()).toContain('flex-row')
    expect(debtItem.classes()).toContain('justify-between')
    expect(debtItem.classes()).toContain('px-2')
    expect(debtItem.classes()).toContain('py-1')
  })
})