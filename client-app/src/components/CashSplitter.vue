<template>
  <el-table @selection-change="handleSelectionChange" class="w-full h-full" :data="tableData" style="width: 100%">
    <el-table-column type="selection"/>
    <el-table-column property="name" label="Имя"/>
    <el-table-column
        label="Телеграм"
        property="telegram"
    />
    <el-table-column
        label="Будет должен"
    >
      <template #default="scope">{{ debt }}</template>
    </el-table-column>
  </el-table>
  <div>
    <el-switch
        v-model="isCountMe"
        class="w-full"
        active-text="посчитать вместе со мной"
        inactive-text="посчитать без меня"
    />
    <p class="text-left text-sm">Сумма</p>
    <el-input-number size="large" class="w-full" v-model="num" :min="100" :max="10000000" :step="100">
      <template #suffix>
        <span>₽</span>
      </template>
    </el-input-number>
  </div>
  <el-button size="large" type="primary" class="w-full"> Поделить</el-button>
</template>

<script setup lang="ts">
import {computed, ref} from "vue";

interface User {
  name: string
  telegram: string
}

const multipleSelection = ref<User[]>([])
const num = ref(100)
const isCountMe = ref(true);

const debt = computed(() => ((multipleSelection.value.length) === 0) ? 0 : Math.ceil(num.value / (multipleSelection.value.length + +isCountMe.value)));

const handleSelectionChange = (val: User[]) => {
  multipleSelection.value = val
}

const tableData: User[] = [
  {
    name: 'Aleyna Kutzner',
    telegram: '@aleyna',
  },
  {
    name: 'Helen Jacobi',
    telegram: '@halen',
  },
  {
    name: 'Brandon Deckert',
    telegram: '@dec',
  },
  {
    name: 'Margie Smith',
    telegram: '@margo',
  },
]
</script>

