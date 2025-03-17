<script setup lang="ts">
import {computed} from "vue";

const props = defineProps<{
  name: string
}>();

const displayName = computed(() => {
  const [left, right] = props.name.split(" ");

  if (left !== undefined && right !== undefined) {
    return left[0] + right[0];
  } else if (left !== undefined) {
    return left[0];
  }
  return "?";
})

</script>

<template>
  <div class="flex items-center justify-center font-bold text-2xl min-h-16 w-fit aspect-square rounded-full bg-amber-500 text-white">
    {{displayName}}
  </div>
</template>

<style scoped>

</style>