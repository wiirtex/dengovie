import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import NameAvatar from '../NameAvatar.vue';


describe('NameAvatar.vue', () => {
  it('рендеринг компонента с переданным именем', () => {
    const wrapper = mount(NameAvatar, {
      props: {
        name: 'John Doe',
      },
    });

    expect(wrapper.exists()).toBe(true);

    expect(wrapper.text()).toBe('JD');
  });

  it('корректно вычисляет displayName для имени из двух слов', () => {
    const wrapper = mount(NameAvatar, {
      props: {
        name: 'John Doe',
      },
    });

    expect(wrapper.text()).toBe('JD');
  });

  it('корректно вычисляет displayName для имени из одного слова', () => {
    const wrapper = mount(NameAvatar, {
      props: {
        name: 'John',
      },
    });

    expect(wrapper.text()).toBe('J');
  });

  it('возвращает "?" для пустого имени', () => {
    const wrapper = mount(NameAvatar, {
      props: {
        name: '',
      },
    });

    expect(wrapper.text()).toBe('?');
  });

  it('возвращает "?" для имени, состоящего только из пробелов', () => {
    const wrapper = mount(NameAvatar, {
      props: {
        name: '   ',
      },
    });

    expect(wrapper.text()).toBe('?');
  });

  it('корректно обрабатывает имя с лишними пробелами', () => {
    const wrapper = mount(NameAvatar, {
      props: {
        name: '  John Doe  ',
      },
    });

    expect(wrapper.text()).toBe('JD');
  });
});