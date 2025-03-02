<script setup lang="ts">
interface User {
  identityNumber: string;
  fullName: string;
  avatarUrl?: string;
}

interface Props {
  user: User;
  class?: string;
}

const props = defineProps<Props>();
</script>

<template>
  <div
    class="flex items-center gap-3 p-4 rounded-lg bg-secondary/50"
    :class="props.class"
  >
    <NuxtImg
      v-if="props.user.avatarUrl"
      :src="props.user.avatarUrl"
      :alt="props.user.fullName"
      class="h-12 w-12 rounded-lg object-cover"
    />
    <div
      v-else
      class="flex h-12 w-12 items-center justify-center rounded-lg bg-muted"
    >
      <svg
        class="h-6 w-6 text-muted-foreground"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
        />
      </svg>
    </div>

    <div class="min-w-0">
      <p class="truncate text-lg font-medium">{{ props.user.fullName }}</p>
      <p class="text-sm text-muted-foreground">
        @{{ props.user.identityNumber }}
      </p>
    </div>
  </div>
</template>
