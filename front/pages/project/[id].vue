<template>
  <div class="p-4 space-y-8">
    <!-- 现有的项目卡片内容 -->
    <Card class="w-full max-w-3xl mx-auto">
      <CardHeader class="space-y-4">
        <CardTitle class="text-2xl md:text-3xl">
          <h3 class="">Title</h3>
        </CardTitle>
        <CardDescription class="flex items-center gap-2 text-lg">
          <span class="">{{ projectData.title }}</span>
        </CardDescription>
      </CardHeader>

      <CardContent class="space-y-8">
        <!-- 项目创建者 -->
        <div v-if="projectData.user" class="pb-6">
          <h3 class="text-lg font-medium mb-4">Recipient</h3>
          <UserCard :user="projectData.user" />
        </div>

        <!-- 项目描述 -->
        <div
          v-if="projectData.description"
          class="prose dark:prose-invert max-w-none"
        >
          <h3 class="text-lg font-medium mb-4">Description</h3>
          <p class="text-base leading-relaxed">{{ projectData.description }}</p>
        </div>

        <!-- 项目图片 -->
        <div v-if="projectData.imgUrl" class="space-y-4">
          <h3 class="text-lg font-medium">Image</h3>
          <div class="relative aspect-video w-full overflow-hidden rounded-lg">
            <NuxtImg
              :src="projectData.imgUrl"
              alt="Project Image"
              class="w-full h-full"
            />
          </div>
        </div>

        <!-- 相关链接 -->
        <div v-if="projectData.link" class="space-y-4">
          <h3 class="text-lg font-medium">Related Link</h3>
          <a
            :href="projectData.link"
            target="_blank"
            class="inline-flex items-center gap-2 text-primary hover:underline"
          >
            <span>{{ projectData.link }}</span>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
              />
            </svg>
          </a>
        </div>
      </CardContent>

      <CardHeader class="text-center">
        <CardTitle class="flex items-center justify-center gap-2">
          <div>
            <ShimmerButton
              size="lg"
              class="shadow-2xl gap-2"
              shimmer-size="2px"
              @click="handleDonate"
            >
              <span
                class="whitespace-pre-wrap text-center text-sm font-medium leading-none tracking-tight text-white lg:text-lg dark:from-white dark:to-slate-900/10"
              >
                Support
              </span>
              <span class="text-rose-500">🌹</span>
            </ShimmerButton>
          </div>
        </CardTitle>

        <CardDescription>
          已有 {{ donateStats.userCount }} 人赞赏，共计 ${{
            donateStats.totalAmount
          }}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex justify-center mb-8"></div>

        <!-- 捐赠列表 -->
        <div class="space-y-4">
          <div
            v-for="action in donateActions"
            :key="action.identityNumber + action.assetId + action.amount"
            class="flex items-start gap-4 p-4 rounded-lg bg-secondary/30"
          >
            <!-- 用户头像 -->
            <NuxtImg
              v-if="action.avatarUrl"
              :src="action.avatarUrl"
              :alt="action.fullName"
              class="w-12 h-12 rounded-full object-cover"
            />
            <div
              v-else
              class="w-12 h-12 rounded-full bg-muted flex items-center justify-center"
            >
              <UserIcon class="w-6 h-6 text-muted-foreground" />
            </div>

            <!-- 捐赠信息 -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between">
                <p
                  class="font-medium truncate max-w-[120px] cursor-pointer hover:text-primary"
                  @click="router.push(`/user/${action.identityNumber}`)"
                >
                  {{ action.fullName }}
                </p>
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
              <p
                class="text-sm text-muted-foreground truncate"
                v-if="action.biography"
              >
                {{ action.biography }}
              </p>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>

  <DonateModal
    v-model:open="showDonateDialog"
    :assets="assets"
    :pid="projectData.pid"
  />
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { UserIcon } from "lucide-vue-next";
import type { Asset, UserAction } from "~/types";

interface ProjectData {
  pid: string;
  identityNumber: string;
  title: string;
  link?: string;
  imgUrl?: string;
  description?: string;
  user?: {
    identityNumber: string;
    fullName: string;
    avatarUrl?: string;
  };
}

const route = useRoute();
const router = useRouter();
const projectData = ref<ProjectData>({
  pid: "",
  identityNumber: "",
  title: "",
  link: "",
  imgUrl: "",
  description: "",
  user: undefined,
});

const donateStats = ref({
  userCount: 0,
  totalAmount: "0.00",
});

const donateActions = ref<UserAction[]>([]);

const formatAmount = (amount: string) => {
  return parseFloat(amount).toFixed(4);
};

// 导入和状态定义
const showDonateDialog = ref(false);
const assets = ref<Asset[]>([]);

// 修改 handleDonate 函数
const handleDonate = async () => {
  try {
    const assetsList = await getAssets();
    if (assetsList) {
      assets.value = assetsList;
      showDonateDialog.value = true;
    }
  } catch (error) {
    console.error("获取资产列表失败:", error);
  }
};

// 模板中的 Modal 组件修改为

onMounted(async () => {
  try {
    const data = await getProject(route.params.id as string);

    projectData.value = {
      pid: data.pid,
      identityNumber: data.identityNumber,
      title: data.title,
      link: data.link || "",
      imgUrl: data.imgUrl || "",
      description: data.description || "",
      user: data.user,
    };
    // 获取捐赠信息
    if (data.pid) {
      const donateUsers = await getProjectDonateUsers(data.pid);
      if (donateUsers && donateUsers.length > 0) {
        donateActions.value = donateUsers;

        // 计算统计信息
        const uniqueUsers = new Set(donateUsers.map((a) => a.identityNumber));
        donateStats.value.userCount = uniqueUsers.size;

        // 计算总金额（USD）
        const total = donateUsers.reduce((sum, action) => {
          return (
            sum + parseFloat(action.amount) * parseFloat(action.asset.priceUsd)
          );
        }, 0);
        donateStats.value.totalAmount = total.toFixed(2);
      } else {
        // 如果没有捐赠记录，保持默认值
        donateActions.value = [];
        donateStats.value = {
          userCount: 0,
          totalAmount: "0.00",
        };
      }
    }
  } catch (error) {
    console.error("获取项目信息失败:", error);
  }
});
</script>
