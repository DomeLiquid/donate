<template>
  <AuroraBackground>
    <Motion
      as="div"
      :initial="{ opacity: 0, y: 40, filter: 'blur(10px)' }"
      :in-view="{ opacity: 1, y: 0, filter: 'blur(0px)' }"
      :transition="{ delay: 0.3, duration: 0.8, ease: 'easeInOut' }"
      class="relative flex flex-col items-center justify-center min-h-screen gap-4 px-4"
    >
      <!-- <div
        class="text-center text-3xl font-bold md:text-5xl dark:text-white"
      >
        创建新项目
      </div> -->
      <!-- 修改 Card 和其内部布局 -->
      <Card class="w-[350px] md:w-[450px] relative">
        <!-- 添加 relative 定位 -->
        <CardHeader>
          <CardTitle>New Project</CardTitle>
          <CardDescription
            >Please fill in the donation project information. Fields marked with
            * are required.</CardDescription
          >
        </CardHeader>
        <form @submit="onSubmit">
          <CardContent class="space-y-6">
            <!-- 修改间距 -->
            <div class="grid w-full items-start gap-6">
              <!-- 修改 gap 和 items-start -->
              <FormField v-slot="{ componentField }" name="identityNumber">
                <FormItem>
                  <FormLabel>Identity Number *</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="eg: 39427696"
                      :maxlength="16"
                      v-bind="componentField"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="title">
                <FormItem>
                  <FormLabel>Title *</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="请作者嚯coffee☕️"
                      :maxlength="32"
                      v-bind="componentField"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <Collapsible>
                <CollapsibleTrigger class="w-full">
                  <div class="flex items-center justify-between rounded-lg p-2">
                    <span class="text-sm">显示更多选项</span>
                    <ChevronDown
                      class="h-4 w-4 transition-transform duration-200"
                      :class="{ 'transform rotate-180': isOpen }"
                    />
                  </div>
                </CollapsibleTrigger>
                <CollapsibleContent class="space-y-4 pt-4">
                  <FormField v-slot="{ componentField }" name="link">
                    <FormItem>
                      <FormLabel>链接</FormLabel>
                      <FormControl>
                        <Input
                          type="url"
                          placeholder="请输入链接地址"
                          v-bind="componentField"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  </FormField>

                  <FormField v-slot="{ componentField }" name="imgUrl">
                    <FormItem>
                      <FormLabel>图片链接</FormLabel>
                      <FormControl>
                        <Input
                          type="url"
                          placeholder="请输入图片链接"
                          v-bind="componentField"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  </FormField>

                  <FormField v-slot="{ componentField }" name="description">
                    <FormItem>
                      <FormLabel>描述</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="请输入描述信息"
                          :maxlength="512"
                          v-bind="componentField"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  </FormField>
                </CollapsibleContent>
              </Collapsible>
            </div>
          </CardContent>
          <CardFooter class="flex justify-between mt-6">
            <!-- 添加上边距 -->
            <Button variant="outline" type="button" @click="handleCancel">
              取消
            </Button>
            <Button type="submit">提交</Button>
          </CardFooter>
        </form>
      </Card>
    </Motion>
  </AuroraBackground>

  <!-- 项目列表部分 -->
  <!-- <div class="container mx-auto p-4">
    <Card class="w-full">
      <CardHeader class="border-b">
        <CardTitle class="text-2xl font-bold">探索项目</CardTitle>
        <CardDescription>
          <div class="flex items-center gap-4 mt-4">
            <div class="relative flex-1 max-w-xs">
              <div class="relative">
                <Input
                  v-model="searchIdentity"
                  placeholder="输入创建者ID搜索"
                  class="pr-10 rounded-r-none border-r-0"
                />
                <Button
                  class="absolute right-0 top-0 h-full rounded-l-none bg-primary hover:bg-primary/90 text-white"
                  @click="handleSearch"
                >
                  搜索
                </Button>
              </div>
            </div>
          </div>
        </CardDescription>
      </CardHeader>

      <CardContent class="pt-6">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div
            v-for="project in projects"
            :key="project.project?.pid"
            class="group relative flex flex-col rounded-xl border hover:shadow-lg transition-all duration-200"
          >
            <div class="flex flex-col flex-1 p-6">
              <div class="flex-1">
                <h3
                  class="text-xl font-semibold hover:text-primary cursor-pointer line-clamp-2"
                  @click="router.push(`/project/${project.project?.pid}`)"
                >
                  {{ project.project?.title }}
                </h3>
              </div>

              <div class="flex items-center justify-between pt-4">
                <div class="flex items-center gap-3">
                  <NuxtImg
                    v-if="project.user?.avatarUrl"
                    :src="project.user.avatarUrl"
                    :alt="project.user?.fullName"
                    class="w-8 h-8 rounded-full object-cover"
                  />
                  <div
                    class="flex flex-col cursor-pointer hover:text-primary"
                    @click="
                      router.push(`/user/${project.user?.identityNumber}`)
                    "
                  >
                    <p class="font-medium text-sm">
                      {{ project.user?.fullName }}
                    </p>
                    <p class="text-xs text-muted-foreground">
                      @{{ project.user?.identityNumber }}
                    </p>
                  </div>
                </div>

                <div class="flex items-center gap-2">
                  <Users2 class="h-4 w-4 text-muted-foreground" />
                  <span class="text-sm text-muted-foreground"
                    ><span class="text-pink-400">{{
                      project.project?.donateCnt || 0
                    }}</span>
                    次捐赠</span
                  >
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="flex justify-center pt-8">
          <Button
            variant="outline"
            size="lg"
            @click="loadMore"
            class="min-w-[200px]"
          >
            加载更多项目
          </Button>
        </div>
      </CardContent>
    </Card>
  </div> -->
</template>

<script setup lang="ts">
import { Motion } from "motion-v";
import { ChevronDown } from "lucide-vue-next";
import { ref } from "vue";
import { useRouter } from "vue-router";
import { toTypedSchema } from "@vee-validate/zod";
import { useForm } from "vee-validate";
import * as z from "zod";
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import type { GetProjectItem } from "~/types";

const isOpen = ref(false);
const router = useRouter();

// 定义表单验证架构
const formSchema = toTypedSchema(
  z.object({
    identityNumber: z
      .string()
      .min(1, { message: "身份编号不能为空" })
      .regex(/^\d+$/, { message: "身份编号必须是纯数字" }),
    title: z
      .string()
      .min(1, { message: "标题不能为空" })
      .max(32, { message: "标题不能超过32个字符" }),
    link: z
      .string()
      .url({ message: "请输入有效的URL" })
      .optional()
      .or(z.literal("")),
    imgUrl: z
      .string()
      .url({ message: "请输入有效的图片URL" })
      .optional()
      .or(z.literal("")),
    description: z
      .string()
      .max(512, { message: "描述不能超过512个字符" })
      .optional()
      .or(z.literal("")),
  })
);

// 使用 vee-validate 的 useForm
const { handleSubmit, resetForm } = useForm({
  validationSchema: formSchema,
  initialValues: {
    identityNumber: "",
    title: "",
    link: "",
    imgUrl: "",
    description: "",
  },
});

// 提交处理
const onSubmit = handleSubmit((values) => {
  const formJson = JSON.stringify(values);
  const encodedData = btoa(unescape(encodeURIComponent(formJson)));
  router.push(`/project/${encodedData}`);
});

// 取消处理
const handleCancel = () => {
  resetForm();
};

// 新增的项目列表相关状态
const projects = ref<GetProjectItem[]>([]);
const searchIdentity = ref("");
const offset = ref(0);
const limit = 10;

// 获取项目列表
const fetchProjects = async (reset = false) => {
  try {
    const params: { limit: number; offset: number; identity_number?: string } =
      {
        limit,
        offset: offset.value,
      };

    if (searchIdentity.value) {
      params.identity_number = searchIdentity.value;
    }

    return await getProjects(params);
  } catch (error) {
    console.error("获取项目列表失败:", error);
    return [];
  }
};

// 初始加载
onMounted(async () => {
  try {
    const data = await fetchProjects();
    if (data) {
      projects.value = data.items;
    }
  } catch (error) {
    console.error("err :", error);
  }
});

// 搜索项目
const handleSearch = async () => {
  offset.value = 0;
  const data = await fetchProjects();
  if (data) {
    projects.value = data.items;
  }
};

// 加载更多
const loadMore = async () => {
  offset.value += limit;
  const data = await fetchProjects();
  if (data) {
    projects.value = [...projects.value, ...data.items];
  }
};
</script>
