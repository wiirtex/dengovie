<template>
  <section class="card flex flex-col">
    <h1 class="font-bold text-lg text-red-500 mb-4"> Вы должны </h1>
    <div class="overflow-y-auto max-h-[650px] flex flex-col gap-1">
      <div class="flex flex-row items-center justify-between border rounded px-2 py-1" v-for="u of userList" :key="u.telegram">
        <span class="text-sm ">{{ u.name }}</span>
        <span class="flex flex-row gap-2">
          <span class="text-sm font-bold text-red-500">-{{ u.debt }}₽</span>
          <el-button plain type="success" size="small">отдал</el-button>
        </span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">

import {ref} from "vue";

interface User {
  name: string
  telegram: string
  debt: number
}

const userList = ref<User[]>([
  {
    name: "Лев Львович",
    telegram: "lion",
    debt: 60
  },
  {
    name: "Тест тестов",
    telegram: "lion",
    debt: 96
  },
  {
    name: "Крафт",
    telegram: "craft",
    debt: 1200
  },
  {
    name: "Мурз мурзович",
    telegram: "yea",
    debt: 60
  },
  {
    name: "Евгения Степанова",
    telegram: "zhenya_step",
    debt: 700
  },
  {
    name: "Григорий Петров",
    telegram: "grisha_petrov",
    debt: 800
  },
  {
    name: "Людмила Семенова",
    telegram: "lyuda_sem",
    debt: 900
  },
  {
    name: "Роман Кузьмин",
    telegram: "roma_kuz",
    debt: 1000
  },
  {
    name: "Алина Воробьева",
    telegram: "alina_vorob",
    debt: 1100
  },
  {
    name: "Денис Фролов",
    telegram: "denis_frol",
    debt: 1200
  },
  {
    name: "Кристина Егорова",
    telegram: "kristina_egor",
    debt: 1300
  },
  {
    name: "Валерий Дмитриев",
    telegram: "valera_dmitr",
    debt: 1400
  },
  {
    name: "Ирина Алексеева",
    telegram: "irina_alex",
    debt: 1500
  },
  {
    name: "Станислав Борисов",
    telegram: "stas_boris",
    debt: 1600
  },
  {
    name: "Максим Тимофеев",
    telegram: "max_tim",
    debt: 1700
  },
  {
    name: "Оксана Федотова",
    telegram: "oksana_fed",
    debt: 1800
  },
  {
    name: "Аркадий Жуков",
    telegram: "arkadiy_zhuk",
    debt: 1900
  },
  {
    name: "Лидия Макарова",
    telegram: "lida_makar",
    debt: 2000
  },
  {
    name: "Василий Крылов",
    telegram: "vasya_kryl",
    debt: 2100
  },
  {
    name: "Зоя Никитина",
    telegram: "zoya_nik",
    debt: 2200
  },
  {
    name: "Геннадий Савельев",
    telegram: "gena_savel",
    debt: 2300
  },
  {
    name: "Лариса Гусева",
    telegram: "lara_guseva",
    debt: 2400
  },
  {
    name: "Федор Тихонов",
    telegram: "fedya_tikh",
    debt: 2500
  },
  {
    name: "Эльвира Давыдова",
    telegram: "elvira_dav",
    debt: 2600
  },
]);
</script>