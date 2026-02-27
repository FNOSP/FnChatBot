<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { locale, t } = useI18n();

// Language handling
const languageOptions = [
  { value: 'en', label: 'English' },
  { value: 'zh', label: '中文' },
  { value: 'ja', label: '日本語' },
];

const currentLanguage = ref(locale.value);

const handleLanguageChange = (value: string | number | any) => {
  const val = String(value);
  locale.value = val;
  currentLanguage.value = val;
  localStorage.setItem('locale', val);
};
</script>

<template>
  <div class="h-full w-full p-4">
    <t-card :title="t('settings.general')" :bordered="false" class="bg-transparent shadow-none">
      <div class="flex flex-col gap-6">
        <div class="flex items-center justify-between p-4 border border-border rounded-lg bg-bg-card">
          <div class="flex flex-col">
            <span class="text-base font-medium text-text-primary">{{ t('settings.language') }}</span>
            <span class="text-sm text-text-secondary">{{ t('settings.languageDesc') }}</span>
          </div>
          <t-select 
            v-model="currentLanguage" 
            style="width: 200px" 
            @change="handleLanguageChange"
            :options="languageOptions"
          />
        </div>
      </div>
    </t-card>
  </div>
</template>

