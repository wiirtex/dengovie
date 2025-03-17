<template>

  <div class="grid grid-cols-4 gap-4 h-full">
    <DebtsIncomeList/>
    <section class="col-span-2 w-full flex flex-col gap-2">
      <UserProfilePanel class="w-full" name="Кадырин Вадим" tg="galaxyshad"/>
      <CashSplitter/>
    </section>
    <DebtsOutcomeList/>
  </div>

</template>

<script setup lang="ts">
import UserProfilePanel from "./components/UserProfilePanel.vue";
import DebtsOutcomeList from "./components/DebtsOutcomeList.vue";
import DebtsIncomeList from "./components/DebtsIncomeList.vue";
import CashSplitter from "./components/CashSplitter.vue";
</script>