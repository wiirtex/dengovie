<script setup lang="ts">
import {reactive} from "vue";
import {login} from "../service/auth.ts";
import {ElMessage} from "element-plus";

const formData = reactive({
  alias: "",
  password: "",
});

const emit = defineEmits<{(e: 'login'): void}>()

async function onSubmit() {
  const {alias, password} = formData;

  const res = await login({alias, password});

  if (res.type === "success") {
    emit('login');
  } else {
    switch (res.error) {
      case "invalid_password":
        ElMessage({
          message: 'Неправильный логин или пароль',
          type: 'error'
        });
        break;
      case "server_error":
        ElMessage({
          message: 'Внутренняя ошибка сервера',
          type: 'error'
        });
        break;
      case 'unknown':
        ElMessage({
          message: 'Неизвестная ошибка',
          type: 'error'
        });
        break;
    }
  }
}
</script>

<template>

  <div class="text-left h-full w-full flex items-center justify-center">
    <el-form v-model="formData" label-position="top" class="my-auto px-4 py-4 bg-white rounded-md shadow-2xl">
      <h1 class="text-center text-4xl font-bold mb-4">dengovie</h1>
      <el-form-item label="Alias">
        <el-input v-model="formData.alias"/>
      </el-form-item>
      <el-form-item label="Password">
        <el-input v-model="formData.password" type="password"/>
      </el-form-item>
      <el-form-item class="mb-0 mt-8">
        <el-button size="large" class="w-full" type="primary" @click="onSubmit">Войти</el-button>
      </el-form-item>
    </el-form>
  </div>

</template>

<style scoped>

</style>