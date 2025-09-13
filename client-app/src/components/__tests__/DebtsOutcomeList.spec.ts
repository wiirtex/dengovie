import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { ElPopconfirm, ElButton, ElInputNumber } from 'element-plus'
import DebtsOutcomeList from '../DebtsOutcomeList.vue'

// Мокаем сервис
vi.mock('@/service/debts', () => ({
  pay: vi.fn()
}))

describe('DebtsOutcomeList.vue', () => {
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
    }
  ]

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('отображает заголовок "Вы должны" красным цветом', () => {
    const wrapper = mount(DebtsOutcomeList, {
      props: {
        debts: []
      },
      global: {
        components: {
          ElPopconfirm,
          ElButton,
          ElInputNumber
        }
      }
    })

    const title = wrapper.find('h1')
    expect(title.text()).toBe('Вы должны')
    expect(title.classes()).toContain('text-red-500')
  })

  it('отображает правильное количество элементов долгов', () => {
    const wrapper = mount(DebtsOutcomeList, {
      props: {
        debts: mockDebts
      },
      global: {
        components: {
          ElPopconfirm,
          ElButton,
          ElInputNumber
        }
      }
    })

    const debtItems = wrapper.findAll('.border.rounded')
    expect(debtItems).toHaveLength(mockDebts.length)
  })

  it('корректно отображает информацию о каждом долге', () => {
    const wrapper = mount(DebtsOutcomeList, {
      props: {
        debts: mockDebts
      },
      global: {
        components: {
          ElPopconfirm,
          ElButton,
          ElInputNumber
        }
      }
    })

    const debtItems = wrapper.findAll('.border.rounded')

    mockDebts.forEach((debt, index) => {
      const item = debtItems[index]
      expect(item.text()).toContain(debt.another_user_name)
      expect(item.text()).toContain(`${debt.amount}₽`)
      expect(item.text()).toContain('отдал')
    })
  })

  it('отображает суммы долгов красным цветом', () => {
    const wrapper = mount(DebtsOutcomeList, {
      props: {
        debts: mockDebts
      },
      global: {
        components: {
          ElPopconfirm,
          ElButton,
          ElInputNumber
        }
      }
    })

    const amountElements = wrapper.findAll('.text-red-500')
    expect(amountElements).toHaveLength(mockDebts.length + 1) // +1 для заголовка

    mockDebts.forEach((debt, index) => {
      const amountElement = amountElements[index + 1]
      expect(amountElement.text()).toBe(`${debt.amount}₽`)
    })
  })

  it('отображает popconfirm с input-number для ввода суммы', async () => {
    const wrapper = mount(DebtsOutcomeList, {
      props: {
        debts: mockDebts
      },
      global: {
        components: {
          ElPopconfirm,
          ElButton,
          ElInputNumber
        }
      }
    })

    const popconfirms = wrapper.findAllComponents(ElPopconfirm)
    expect(popconfirms).toHaveLength(mockDebts.length)

    const inputNumbers = wrapper.findAllComponents(ElInputNumber)
    expect(inputNumbers).toHaveLength(mockDebts.length)

    inputNumbers.forEach(input => {
      expect(input.props('min')).toBe(1)
      expect(input.props('step')).toBe(100)
      expect(input.props('size')).toBe('small')
    })
  })

  it('работает с пустым массивом долгов', () => {
    const wrapper = mount(DebtsOutcomeList, {
      props: {
        debts: []
      },
      global: {
        components: {
          ElPopconfirm,
          ElButton,
          ElInputNumber
        }
      }
    })

    const debtItems = wrapper.findAll('.border.rounded')
    expect(debtItems).toHaveLength(0)

    const popconfirms = wrapper.findAllComponents(ElPopconfirm)
    expect(popconfirms).toHaveLength(0)
  })

  it('корректно принимает props', () => {
    const wrapper = mount(DebtsOutcomeList, {
      props: {
        debts: mockDebts
      },
      global: {
        components: {
          ElPopconfirm,
          ElButton,
          ElInputNumber
        }
      }
    })

    expect(wrapper.props().debts).toEqual(mockDebts)
  })

  it('имеет правильную структуру контейнера', () => {
    const wrapper = mount(DebtsOutcomeList, {
      props: {
        debts: mockDebts
      },
      global: {
        components: {
          ElPopconfirm,
          ElButton,
          ElInputNumber
        }
      }
    })

    const container = wrapper.find('.overflow-y-auto')
    expect(container.classes()).toContain('max-h-[650px]')
    expect(container.classes()).toContain('flex')
    expect(container.classes()).toContain('flex-col')
    expect(container.classes()).toContain('gap-1')
  })
})