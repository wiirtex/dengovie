<template>
  <div class="grid grid-cols-4 gap-4 h-full">
    <DebtsIncomeList/>
    <section class="col-span-2 w-full flex flex-col gap-2">
      <UserProfilePanel class="w-full" :name="user.user_name" :tg="user.user_alias" @logout="$emit('logout')"/>
      <CashSplitter/>
    </section>
    <DebtsOutcomeList/>
  </div>
</template>

<script setup lang="ts">
import DebtsIncomeList from "./DebtsIncomeList.vue";
import UserProfilePanel from "./UserProfilePanel.vue";
import CashSplitter from "./CashSplitter.vue";
import DebtsOutcomeList from "./DebtsOutcomeList.vue";
import type {Me} from "../service/user.ts";

defineProps<{
  user: Me
}>()

defineEmits<{logout: []}>()
</script>