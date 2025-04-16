<template>
  <el-table v-loading="isFetching" @selection-change="handleSelectionChange" class="w-full h-full shadow-lg rounded-md"
            :data="tableData" style="width: 100%">
    <el-table-column type="selection"/>
    <el-table-column property="Name" label="Имя"/>
    <el-table-column
        label="Телеграм"
        property="Alias"
    />
    <el-table-column
        label="Будет должен"
    >
      <template #default="scope">{{ multipleSelection.includes(scope.row) ? debt : 0 }}</template>
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
  <el-button :loading="isSharing" @click="handleOnSplit" size="large" type="primary" class="w-full"> Поделить</el-button>
</template>

<script setup lang="ts">
import {computed, onBeforeMount, ref} from "vue";
import {type GroupUser, users} from "../service/groups.ts";
import {share} from "../service/debts.ts";

const multipleSelection = ref<GroupUser[]>([])
const num = ref(100)
const isCountMe = ref(true);
const isFetching = ref(true);

const isSharing = ref(false);

const props = defineProps<{
  meId: number
}>();
const emit = defineEmits<{ split: [number, number[]] }>()

onBeforeMount(async () => {
  isFetching.value = true;

  const userList = await users();

  if (userList.type === 'success') {
    tableData.value = userList.data;
  }

  isFetching.value = false;
})

const debt = computed(() => ((multipleSelection.value.length) === 0) ? 0 : Math.ceil(num.value / (multipleSelection.value.length + +isCountMe.value)));

const handleSelectionChange = (val: GroupUser[]) => {
  multipleSelection.value = val
}

const handleOnSplit = async () => {
  isSharing.value = true;

  let usrIds = multipleSelection.value.map(item => +item.ID);
  const amount = num.value;

  if (isCountMe.value) {
    usrIds = [props.meId, ...usrIds];
  }

  await share(amount, usrIds);

  emit('split', amount, usrIds);

  isSharing.value = false;
}

const tableData = ref<GroupUser[]>([]);
</script>

