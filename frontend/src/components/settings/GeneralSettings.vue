<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { Layout, Card, Select, SelectOption } from '@kousum/semi-ui-vue';

const { locale, t } = useI18n();

// Language handling
const languageOptions = [
  { value: 'en', label: 'English' },
  { value: 'zh', label: '中文' },
  { value: 'ja', label: '日本語' },
];

const currentLanguage = ref(locale.value);

const handleLanguageChange = (value: string | number | any[]) => {
  const val = String(value);
  locale.value = val;
  currentLanguage.value = val;
  localStorage.setItem('locale', val);
};
</script>

<template>
  <Layout class="h-full w-full">
    <Card :title="t('settings.general')" :bordered="false" class="h-full bg-transparent shadow-none">
      <div class="flex flex-col gap-6">
        <div class="flex items-center justify-between p-4 border border-zinc-200 dark:border-zinc-800 rounded-lg bg-card">
          <div class="flex flex-col">
            <span class="text-base font-medium">{{ t('settings.language') }}</span>
            <span class="text-sm text-zinc-500">{{ t('settings.languageDesc') }}</span>
          </div>
          <Select 
            :value="currentLanguage" 
            style="width: 200px" 
            @change="handleLanguageChange"
          >
            <SelectOption v-for="opt in languageOptions" :key="opt.value" :value="opt.value" :label="opt.label">
              {{ opt.label }}
            </SelectOption>
          </Select>
        </div>

      </div>
    </Card>
  </Layout>
</template>

<style scoped>
/* Ensure Semi UI components blend with existing Tailwind styles if needed */
:deep(.semi-card-header) {
  padding-left: 0;
  padding-right: 0;
}
:deep(.semi-card-body) {
  padding-left: 0;
  padding-right: 0;
}
</style>
