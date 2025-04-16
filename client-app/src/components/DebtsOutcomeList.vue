<template>
  <section class="card flex flex-col">
    <h1 class="font-bold text-lg text-red-500 mb-4"> Вы должны </h1>
    <div class="overflow-y-auto max-h-[650px] flex flex-col gap-1">
      <div class="flex flex-row items-center justify-between border rounded px-2 py-1" v-for="u of props.debts" :key="u.another_user_id">
        <span class="text-sm ">{{ u.another_user_name }}</span>
        <span class="flex flex-row gap-2">
          <span class="text-sm font-bold text-red-500">{{ u.amount }}₽</span>

          <el-popconfirm
              width="220"
              icon-color="#626AEF"
              title="Сколько?"
              @confirm="onPay(u.another_user_id)"
          >
            <template #reference>
              <el-button link type="success" size="small">отдал</el-button>
            </template>
            <template #actions="{ confirm }">
              <el-input-number  class="mb-2 w-full" v-model="amount" size="small" :min="1" :step="100"/>
              <el-button
                  type="primary"
                  size="small"
                  @click="confirm"
              >
                Подтвердить
              </el-button>
            </template>
          </el-popconfirm>



        </span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import {type DebtItem, pay} from "../service/debts.ts";
import {ref} from "vue";

const props = defineProps<{
  debts: DebtItem[]
}>()

const amount = ref<number>(100);

const emit = defineEmits<{ pay: [number, number] }>()

async function onPay(id: number) {
  await pay(amount.value, id);

  emit("pay", id, amount.value);
}
</script>