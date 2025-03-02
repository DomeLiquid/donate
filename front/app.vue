<template>
  <div data-theme="dark">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<script setup>
onMounted(() => {
  document.documentElement.classList.add("dark");
});
</script>
