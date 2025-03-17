<template>
  <section class="card flex flex-col">
    <h1 class="font-bold text-lg text-green-500 mb-4"> Вам должны </h1>
    <div class="overflow-y-auto max-h-[650px] h-full flex flex-col gap-1">
      <div class="flex flex-row justify-between border rounded px-2 py-1" v-for="u of userList" :key="u.telegram">
        <span class="text-sm ">{{ u.name }}</span>
        <span class="text-sm font-bold text-green-500">{{ u.debt }}₽</span>
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
    name: "Иван Иванов",
    telegram: "ivan_ivanov",
    debt: 150
  },
  {
    name: "Мария Петрова",
    telegram: "maria_p",
    debt: 300
  },
  {
    name: "Алексей Смирнов",
    telegram: "alex_smirnov",
    debt: 450
  },
  {
    name: "Екатерина Васильева",
    telegram: "katya_v",
    debt: 200
  },
  {
    name: "Дмитрий Козлов",
    telegram: "dima_kozlov",
    debt: 100
  },
  {
    name: "Ольга Новикова",
    telegram: "olga_nov",
    debt: 50
  },
  {
    name: "Сергей Морозов",
    telegram: "serg_moroz",
    debt: 600
  },
  {
    name: "Анна Кузнецова",
    telegram: "anna_kuz",
    debt: 250
  },
  {
    name: "Павел Волков",
    telegram: "pavel_volk",
    debt: 700
  },
  {
    name: "Татьяна Лебедева",
    telegram: "tanya_leb",
    debt: 400
  },
  {
    name: "Николай Соколов",
    telegram: "kolya_sokol",
    debt: 800
  },
  {
    name: "Елена Павлова",
    telegram: "lena_pavlova",
    debt: 350
  },
  {
    name: "Артем Федоров",
    telegram: "artem_fed",
    debt: 900
  },
  {
    name: "Виктория Михайлова",
    telegram: "vika_mikh",
    debt: 500
  },
  {
    name: "Александр Белов",
    telegram: "sasha_belov",
    debt: 1000
  },
  {
    name: "Юлия Ковалева",
    telegram: "yulia_kov",
    debt: 150
  },
  {
    name: "Михаил Николаев",
    telegram: "misha_nik",
    debt: 200
  },
]);

</script>

<style scoped>

</style>