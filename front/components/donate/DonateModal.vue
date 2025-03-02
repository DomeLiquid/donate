<script setup lang="ts">
import { ref, computed } from "vue";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from "@/components/ui/form";
import { useForm } from "vee-validate";
import { toTypedSchema } from "@vee-validate/zod";
import * as z from "zod";
import { v4 as uuidv4 } from "uuid";
import type { Asset } from "~/types";

// 修改 props 定义
const props = defineProps<{
  open: boolean;
  assets: Asset[];
  pid: string; // 添加 pid
}>();

const emit = defineEmits<{
  "update:open": [value: boolean];
  submit: [assetId: string, amount: string];
}>();

const selectedAsset = ref<Asset | null>(null);

const donateFormSchema = toTypedSchema(
  z.object({
    assetId: z.string().min(1, "请选择资产"),
    amount: z.string().min(1, "请输入金额"),
  })
);

const form = useForm({
  validationSchema: donateFormSchema,
  initialValues: {
    assetId: "",
    amount: "",
  },
});

// 预设的捐赠数量选项
const donateAmounts = [0.003, 0.012, 0.024, 0.096, 0.768, 3.072];

// 计算美元价值
const usdValue = computed(() => {
  if (!selectedAsset.value || !form.values.amount) return "0.00";
  const amount = parseFloat(form.values.amount);
  const price = parseFloat(selectedAsset.value.priceUsd);
  return (amount * price).toFixed(2);
});

// 选择预设数量
// 修改选择预设数量的函数
const selectPresetAmount = (amount: number) => {
  const formattedAmount = amount.toString(); // 不使用 toFixed 以避免不必要的格式化
  form.setFieldValue("amount", formattedAmount);
};

// 添加处理输入的函数
const handleAmountInput = (e: Event) => {
  const target = e.target as HTMLInputElement;
  let value = target.value;

  // 移除前导零
  if (value.length > 1 && value.startsWith("0") && !value.startsWith("0.")) {
    value = value.replace(/^0+/, "");
  }

  // 限制小数位数为 4 位
  if (value.includes(".")) {
    const [integer, decimal] = value.split(".");
    if (decimal && decimal.length > 4) {
      value = `${integer}.${decimal.slice(0, 4)}`;
    }
  }

  form.setFieldValue("amount", value);
};

const handleAssetSelect = (assetId: string) => {
  selectedAsset.value =
    props.assets.find((asset) => asset.assetId === assetId) || null;
  // 设置表单的 assetId 值
  form.setFieldValue("assetId", assetId);
};

// 修改 handleSubmit 函数
const handleSubmit = form.handleSubmit((values) => {
  const botId = useRuntimeConfig().public.botId;
  const trace = uuidv4();
  if (!props.pid || props.pid.length === 0) {
    console.log("pid is empty");
    return;
  }
  const payUrl = `https://mixin.one/pay/${botId}?asset=${values.assetId}&amount=${values.amount}&trace=${trace}&memo=${props.pid}`;
  window.location.href = payUrl;
  emit("update:open", false);
});
</script>

<template>
  <Dialog :open="open" @update:open="(val) => emit('update:open', val)">
    <DialogContent class="sm:max-w-[425px]">
      <DialogHeader>
        <DialogTitle class="text-xl font-bold">赞赏项目</DialogTitle>
      </DialogHeader>

      <form @submit="handleSubmit" class="space-y-6">
        <!-- 资产选择部分 -->
        <FormField v-slot="{ field }" name="assetId" :control="form.control">
          <FormItem>
            <FormLabel>选择资产</FormLabel>
            <Select
              v-model="field.value"
              @update:model-value="handleAssetSelect"
            >
              <FormControl>
                <SelectTrigger class="w-full">
                  <SelectValue placeholder="选择赞赏资产">
                    <div v-if="selectedAsset" class="flex items-center gap-3">
                      <div class="relative">
                        <NuxtImg
                          :src="selectedAsset.iconUrl"
                          :alt="selectedAsset.symbol"
                          class="w-6 h-6 rounded-full"
                        />
                        <NuxtImg
                          :src="selectedAsset.chainIconUrl"
                          class="absolute -bottom-1 -right-1 w-4 h-4 rounded-full border-2 border-background"
                        />
                      </div>
                      <span>{{ selectedAsset.symbol }}</span>
                    </div>
                  </SelectValue>
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                <SelectItem
                  v-for="asset in assets"
                  :key="asset.assetId"
                  :value="asset.assetId"
                >
                  <div class="flex items-center gap-2">
                    <div class="relative">
                      <NuxtImg
                        :src="asset.iconUrl"
                        :alt="asset.symbol"
                        class="w-6 h-6 rounded-full"
                      />
                      <NuxtImg
                        :src="asset.chainIconUrl"
                        class="absolute -bottom-1 -right-1 w-4 h-4 rounded-full border-2 border-background"
                      />
                    </div>
                    <div class="flex flex-col">
                      <span class="font-medium">{{ asset.symbol }}</span>
                      <span class="text-xs text-muted-foreground"
                        >≈ ${{ asset.priceUsd }}</span
                      >
                    </div>
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        </FormField>

        <!-- 金额选择部分 -->
        <FormField v-slot="{ field }" name="amount" :control="form.control">
          <FormItem>
            <FormLabel>选择数量</FormLabel>
            <div class="grid grid-cols-3 gap-2 mb-4">
              <Button
                v-for="amount in donateAmounts"
                :key="amount"
                type="button"
                variant="outline"
                :class="{
                  'border-primary': form.values.amount === amount.toString(),
                }"
                @click="selectPresetAmount(amount)"
              >
                {{ amount }}
              </Button>
            </div>
            <!-- 修改金额选择部分的 Input 组件 -->
            <!-- 修改金额输入框部分 -->
            <FormControl>
              <div class="relative">
                <Input
                  :value="field.value"
                  @input="handleAmountInput"
                  type="text"
                  inputmode="decimal"
                  pattern="[0-9]*[.]?[0-9]*"
                  placeholder="输入自定义数量"
                  autocomplete="off"
                />
                <div
                  v-if="selectedAsset"
                  class="absolute right-3 top-1/2 -translate-y-1/2 flex items-center gap-2"
                >
                  <NuxtImg
                    :src="selectedAsset.iconUrl"
                    :alt="selectedAsset.symbol"
                    class="w-5 h-5 rounded-full"
                  />
                </div>
              </div>
            </FormControl>
            <div class="mt-2 text-sm text-muted-foreground">
              <template v-if="selectedAsset && form.values.amount">
                预计金额：${{ usdValue }}
              </template>
              <template v-else-if="!selectedAsset"> 请先选择资产 </template>
              <template v-else> 请输入捐赠数量 </template>
            </div>
            <FormMessage />
          </FormItem>
        </FormField>

        <div class="flex justify-end gap-4">
          <Button
            type="button"
            variant="outline"
            @click="$emit('update:open', false)"
          >
            取消
          </Button>
          <Button type="submit">确认赞赏</Button>
        </div>
      </form>
    </DialogContent>
  </Dialog>
</template>
