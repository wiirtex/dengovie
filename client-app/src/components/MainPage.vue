<template>
  <div class="grid grid-cols-4 gap-4 h-full">
    <DebtsIncomeList v-loading="isDebtsFetching" :debts="debtsIncome"/>
    <section class="col-span-2 w-full flex flex-col gap-2">
      <UserProfilePanel class="w-full" :name="user.name" :tg="user.alias" @logout="$emit('logout')"/>
      <CashSplitter @split="updDebts" :me-id="+user.user_id"/>
    </section>
    <DebtsOutcomeList v-loading="isDebtsFetching" :debts="debtsOutcome" @pay="updDebts"/>
  </div>
</template>

<script setup lang="ts">
import DebtsIncomeList from "./DebtsIncomeList.vue";
import UserProfilePanel from "./UserProfilePanel.vue";
import CashSplitter from "./CashSplitter.vue";
import DebtsOutcomeList from "./DebtsOutcomeList.vue";
import type {Me} from "../service/user.ts";
import {computed, onBeforeMount, ref} from "vue";
import {all, type DebtItem} from "../service/debts.ts";

const props = defineProps<{
  user: Me
}>()

const debtsOutcome = computed<DebtItem[]>(() => debtList.value.filter(x => x.amount < 0));
const debtsIncome = computed<DebtItem[]>(() => debtList.value.filter(x => x.amount > 0));

const isDebtsFetching = ref(true);
const debtList = ref<DebtItem[]>([]);

async function updDebts() {
  isDebtsFetching.value = true;

  const data = await all();

  if (data.type === 'success') {
    debtList.value = data.data.Debts;
  }

  isDebtsFetching.value = false;
}

onBeforeMount(updDebts)

defineEmits<{ logout: [] }>()
</script>