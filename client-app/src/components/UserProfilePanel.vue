<script setup lang="ts">

import NameAvatar from "./NameAvatar.vue";

const props = defineProps<{
  name: string
  tg: string
}>();

defineEmits<{
  logout: []
}>()

</script>

<template>
  <section class="card flex flex-row items-center">
    <NameAvatar class="h-24" :name="props.name"/>
    <section class="w-full flex flex-col ml-4 text-left justify-center">
      <div class="text-xl font-bold">
        {{name}}
      </div>
      <div class="text-black/50">
        {{tg}}
      </div>
    </section>
    <el-button @click="$emit('logout')">Выйти</el-button>
  </section>
</template>

<style scoped>

</style>