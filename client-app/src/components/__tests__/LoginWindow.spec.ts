import {describe, it, expect, vi, type Mock} from 'vitest'; // Импортируем Vitest
import { mount } from '@vue/test-utils'; // Импортируем mount из @vue/test-utils
import LoginWindow from '../LoginWindow.vue'; // Импортируем компонент
import { ElMessage } from 'element-plus'; // Импортируем ElMessage для моков

import * as auth from '../../service/auth.ts'

vi.spyOn(auth, 'login');

vi.mock('element-plus', () => ({
  ElMessage: vi.fn(),
}));

import { login } from '../../service/auth.ts';

describe('LoginWindow.vue', () => {
  it('рендеринг формы входа', () => {
    const wrapper = mount(LoginWindow);
    expect(wrapper.find('h1').text()).toBe('dengovie');
    expect(wrapper.findAll('.el-input').length).toBe(2);
  });

  it('ввод данных в форму и отправка', async () => {
    const wrapper = mount(LoginWindow);
    const aliasInput = wrapper.find('input[type="text"]');
    const passwordInput = wrapper.find('input[type="password"]');
    const submitButton = wrapper.find('.el-button');

    await aliasInput.setValue('testuser');
    await passwordInput.setValue('testpass');
    await submitButton.trigger('click');

    expect(login).toHaveBeenCalledWith({ alias: 'testuser', password: 'testpass' });
  });

  it('успешный вход и вызов события login', async () => {
    (login as Mock).mockResolvedValue({ type: 'success' });
    const wrapper = mount(LoginWindow);
    await wrapper.find('.el-button').trigger('click');
    expect(wrapper.emitted()).toHaveProperty('login');
  });

  it('обработка ошибки invalid_password', async () => {
    (login as Mock).mockResolvedValue({ type: 'error', error: 'invalid_password' });
    const wrapper = mount(LoginWindow);
    await wrapper.find('.el-button').trigger('click');
    expect(ElMessage).toHaveBeenCalledWith({
      message: 'Неправильный логин или пароль',
      type: 'error',
    });
  });

  it('обработка ошибки server_error', async () => {
    (login as Mock).mockResolvedValue({ type: 'error', error: 'server_error' });
    const wrapper = mount(LoginWindow);
    await wrapper.find('.el-button').trigger('click');
    expect(ElMessage).toHaveBeenCalledWith({
      message: 'Внутренняя ошибка сервера',
      type: 'error',
    });
  });

  it('обработка неизвестной ошибки', async () => {
    (login as Mock).mockResolvedValue({ type: 'error', error: 'unknown' });
    const wrapper = mount(LoginWindow);
    await wrapper.find('.el-button').trigger('click');
    expect(ElMessage).toHaveBeenCalledWith({
      message: 'Неизвестная ошибка',
      type: 'error',
    });
  });
});