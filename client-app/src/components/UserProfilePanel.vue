<template>
  <section class="card flex flex-row items-center">
    <NameAvatar class="h-24" :name="props.name"/>
    <section class="w-full flex flex-col ml-4 text-left justify-center">
      <div v-if="!isEditing" class="text-xl font-bold">
        {{name}}
      </div>
      <div v-else class="flex flex-row gap-2 mr-2 text-xl font-bold">
        <el-input v-model="changedName" :placeholder="name"/>
        <el-button @click="onNameChange">Применить</el-button>
      </div>
      <div class="text-black/50">
        @{{tg}}
      </div>
    </section>
    <div class="w-48 justify-end flex flex-row flex-wrap gap-2">
      <el-button @click="isEditing = !isEditing">{{isEditing ? 'Отмена' : '✏️'}}</el-button>
      <el-button v-if="!isEditing" @click="$emit('logout')">Выйти</el-button>

      <el-button @click="onDeleteAccount" type="danger" v-if="isEditing">Удалить профиль</el-button>
    </div>
  </section>
</template>

<script setup lang="ts">

import NameAvatar from "./NameAvatar.vue";
import {ref} from "vue";
import {deleteAccount, updName} from "../service/user.ts";
import {ElMessage} from "element-plus";

const changedName = ref<string>();
const isEditing = ref(false);

const props = defineProps<{
  name: string
  tg: string
}>();

const emit = defineEmits<{
  logout: []
}>()

async function onNameChange() {
  await updName(changedName.value);
  emit('logout');
}

async function onDeleteAccount() {
  const req = await deleteAccount();

  if (req.type === 'failure') {
    ElMessage.error("У вас есть долги");
  } else {
    emit('logout');
  }
}

</script>