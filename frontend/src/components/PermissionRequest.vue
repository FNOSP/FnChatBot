<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { Modal, Button } from '@kousum/semi-ui-vue';

const props = defineProps<{
  visible: boolean;
  requestId: string;
  requestedPath: string;
  command: string;
}>();

const emit = defineEmits<{
  close: [];
  respond: [payload: { requestId: string; approved: boolean; remember: boolean }];
}>();

const { t } = useI18n();

const handleAllow = () => {
  emit('respond', {
    requestId: props.requestId,
    approved: true,
    remember: false,
  });
  emit('close');
};

const handleAllowAndRemember = () => {
  emit('respond', {
    requestId: props.requestId,
    approved: true,
    remember: true,
  });
  emit('close');
};

const handleDeny = () => {
  emit('respond', {
    requestId: props.requestId,
    approved: false,
    remember: false,
  });
  emit('close');
};
</script>

<template>
  <Modal
    :visible="visible"
    :title="t('sandbox.permissionTitle')"
    :footer="null"
    @cancel="emit('close')"
  >
    <div class="flex flex-col gap-4">
      <p class="text-sm text-zinc-600 dark:text-zinc-400">
        {{ t('sandbox.permissionDesc') }}
      </p>

      <div class="flex flex-col gap-2">
        <span class="text-sm font-medium text-zinc-700 dark:text-zinc-300">
          {{ t('sandbox.path') }}
        </span>
        <div class="px-3 py-2 bg-amber-100 dark:bg-amber-900/30 border border-amber-300 dark:border-amber-700 rounded-md">
          <code class="text-sm text-amber-800 dark:text-amber-200 break-all">
            {{ requestedPath }}
          </code>
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <span class="text-sm font-medium text-zinc-700 dark:text-zinc-300">
          {{ t('sandbox.command') }}
        </span>
        <div class="px-3 py-2 bg-zinc-100 dark:bg-zinc-800 border border-zinc-300 dark:border-zinc-700 rounded-md overflow-x-auto">
          <pre class="text-sm text-zinc-800 dark:text-zinc-200 font-mono whitespace-pre-wrap break-all">{{ command }}</pre>
        </div>
      </div>

      <div class="flex justify-end gap-3 mt-4">
        <Button
          type="tertiary"
          @click="handleDeny"
        >
          {{ t('sandbox.deny') }}
        </Button>
        <Button
          type="primary"
          @click="handleAllow"
        >
          {{ t('sandbox.allow') }}
        </Button>
        <Button
          type="primary"
          theme="solid"
          @click="handleAllowAndRemember"
        >
          {{ t('sandbox.allowAndRemember') }}
        </Button>
      </div>
    </div>
  </Modal>
</template>

<style scoped>
</style>
