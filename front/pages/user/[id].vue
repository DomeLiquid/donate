<template>
  <div class="container p-4 space-y-8">
    <Card class="w-full max-w-3xl mx-auto">
      <!-- 用户基本信息 -->
      <CardHeader class="space-y-4">
        <CardTitle class="text-2xl md:text-3xl flex items-center gap-4">
          <NuxtImg
            v-if="userData.avatarUrl"
            :src="userData.avatarUrl"
            :alt="userData.fullName"
            class="w-16 h-16 rounded-full object-cover"
          />
          <div class="flex flex-col">
            <h3>{{ userData.fullName }}</h3>
            <p class="text-base text-muted-foreground">
              @{{ userData.identityNumber }}
            </p>
          </div>
        </CardTitle>
        <CardDescription v-if="userData.biography" class="text-lg">
          {{ userData.biography }}
        </CardDescription>
      </CardHeader>

      <CardContent class="space-y-8">
        <!-- 捐赠统计 -->
        <div class="border-b pb-6">
          <h3 class="text-lg font-medium mb-4">捐赠统计</h3>
          <div class="grid grid-cols-2 gap-4">
            <div class="text-center">
              <p class="text-2xl font-bold">{{ donateStats.projectCount }}</p>
              <p class="text-sm text-muted-foreground">已捐赠项目</p>
            </div>
            <div class="text-center">
              <p class="text-2xl font-bold">${{ donateStats.totalAmount }}</p>
              <p class="text-sm text-muted-foreground">总捐赠金额</p>
            </div>
          </div>
        </div>

        <!-- 捐赠历史 -->
        <div class="space-y-4">
          <h3 class="text-lg font-medium">捐赠历史</h3>
          <div
            v-for="action in donateActions"
            :key="action.project.pid + action.assetId"
            class="flex items-start gap-4 p-4 rounded-lg bg-secondary/30"
          >
            <!-- 项目信息 -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between">
                <div class="flex flex-col">
                  <p
                    class="font-medium truncate cursor-pointer hover:text-primary"
                    @click="router.push(`/project/${action.project.pid}`)"
                  >
                    {{ action.project.title }}
                  </p>

                  <div
                    class="flex items-center gap-3 text-sm text-muted-foreground"
                  >
                    <NuxtImg
                      v-if="action.user?.avatarUrl"
                      :src="action.user.avatarUrl"
                      :alt="action.user?.fullName"
                      class="w-10 h-10 rounded-full object-cover"
                    />
                    <div
                      class="flex flex-col cursor-pointer hover:text-primary"
                      @click="
                        router.push(`/user/${action.user?.identityNumber}`)
                      "
                    >
                      <p class="font-medium">{{ action.user?.fullName }}</p>
                      <p class="text-xs">@{{ action.user?.identityNumber }}</p>
                    </div>
                  </div>
                </div>
                <div class="flex items-center gap-2">
                  <div class="relative">
                    <NuxtImg
                      :src="action.asset.iconUrl"
                      :alt="action.asset.symbol"
                      class="w-6 h-6 rounded-full"
                    />
                    <NuxtImg
                      :src="action.asset.chainIconUrl"
                      :alt="action.asset.chainSymbol"
                      class="absolute -bottom-1 -right-1 w-4 h-4 rounded-full border-2 border-background"
                    />
                  </div>
                  <div class="flex flex-col items-end">
                    <p class="text-sm text-muted-foreground">
                      {{ formatAmount(action.amount) }}
                      {{ action.asset.symbol }}
                    </p>
                    <p class="text-xs text-muted-foreground">
                      ≈ ${{
                        (
                          parseFloat(action.amount) *
                          parseFloat(action.asset.priceUsd)
                        ).toFixed(2)
                      }}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import type { User, UserAction } from "~/types";

const route = useRoute();
const router = useRouter();

const userData = ref<User>({
  identityNumber: "",
  fullName: "",
  mixinUid: "",
  avatarUrl: "",
  biography: "",
  mixinCreatedAt: "",
  createdAt: "",
  updatedAt: "",
});

const donateStats = ref({
  projectCount: 0,
  totalAmount: "0.00",
});

const donateActions = ref<UserAction[]>([]);

const formatAmount = (amount: string) => {
  return parseFloat(amount).toFixed(4);
};

onMounted(async () => {
  try {
    // 获取用户信息
    const user = await getUser(route.params.id as string);
    if (user) {
      userData.value = user;

      // 获取用户捐赠历史
      const donateHistory = await getUserDonateProjects(user.identityNumber);
      if (donateHistory && donateHistory.length > 0) {
        donateActions.value = donateHistory;

        // 计算统计信息
        const uniqueProjects = new Set(donateHistory.map((a) => a.project.pid));
        donateStats.value.projectCount = uniqueProjects.size;

        // 计算总金额（USD）
        const total = donateHistory.reduce((sum, action) => {
          return (
            sum + parseFloat(action.amount) * parseFloat(action.asset.priceUsd)
          );
        }, 0);
        donateStats.value.totalAmount = total.toFixed(2);
      }
    }
  } catch (error) {
    console.error("获取用户信息失败:", error);
  }
});
</script>
