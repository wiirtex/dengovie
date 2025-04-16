<template>
  <section class="card flex flex-col">
    <h1 class="font-bold text-lg text-green-500 mb-4"> Вам должны </h1>
    <div class="overflow-y-auto max-h-[650px] h-full flex flex-col gap-1">
      <div class="flex flex-row justify-between border rounded px-2 py-1" v-for="u of props.debts" :key="u.another_user_id">
        <span class="text-sm ">{{ u.another_user_name }}</span>
        <span class="text-sm font-bold text-green-500">{{ u.amount }}₽</span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type {DebtItem} from "../service/debts.ts";

const props = defineProps<{
  debts: DebtItem[]
}>()
</script>