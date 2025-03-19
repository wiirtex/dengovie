<template>
  <div class="w-full h-full" v-loading="isFetching">
    <MainPage v-if="isLogin" :user="user!" @logout="onLogout"/>
    <LoginPage v-else @login="getProfile"/>
  </div>
</template>

<script setup lang="ts">
import LoginPage from "./components/LoginWindow.vue";
import MainPage from "./components/MainPage.vue";
import {computed, onBeforeMount, ref} from "vue";
import {type Me, me} from "./service/user.ts";
import {logout} from "./service/auth.ts";

const isLogin = computed(() => user.value !== undefined)
const isFetching = ref(true);
const user = ref<Me | undefined>()

async function getProfile() {
  isFetching.value = true;
  const res = await me();

  if (res.type === 'failure') {
    user.value = undefined;
  } else {
    user.value = res.data;
  }
  isFetching.value = false;
}

async function onLogout() {
  isFetching.value = true;
  await logout();
  user.value = undefined;
  isFetching.value = false;
}

onBeforeMount(getProfile)
</script>