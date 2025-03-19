import { config } from '@vue/test-utils';
import {
  ElInput,
  ElButton,
  ElForm,
  ElFormItem,
  ElMessage,
  ElTable,
  ElTableColumn,
  ElSwitch,
  ElInputNumber
} from 'element-plus';

// Глобальная регистрация компонентов Element Plus
config.global.components = {
  ElInput,
  ElButton,
  ElForm,
  ElFormItem,
  ElTable,
  ElTableColumn,
  ElSwitch,
  ElInputNumber,
};

// Глобальная регистрация плагинов Element Plus
config.global.plugins = [ElMessage];