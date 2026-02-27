<script setup lang="ts">
import { useI18n } from 'vue-i18n';

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
  <t-dialog
    :visible="visible"
    :header="t('sandbox.permissionTitle')"
    :footer="false"
    @close="emit('close')"
  >
    <div class="flex flex-col gap-4">
      <p class="text-sm text-text-secondary">
        {{ t('sandbox.permissionDesc') }}
      </p>

      <div class="flex flex-col gap-2">
        <span class="text-sm font-medium text-text-primary">
          {{ t('sandbox.path') }}
        </span>
        <div class="px-3 py-2 bg-warning/10 border border-warning/30 rounded-md">
          <code class="text-sm text-warning break-all">
            {{ requestedPath }}
          </code>
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <span class="text-sm font-medium text-text-primary">
          {{ t('sandbox.command') }}
        </span>
        <div class="px-3 py-2 bg-bg-secondary border border-border rounded-md overflow-x-auto">
          <pre class="text-sm text-text-primary font-mono whitespace-pre-wrap break-all">{{ command }}</pre>
        </div>
      </div>

      <div class="flex justify-end gap-3 mt-4">
        <t-button
          theme="default"
          variant="text"
          @click="handleDeny"
        >
          {{ t('sandbox.deny') }}
        </t-button>
        <t-button
          theme="primary"
          variant="outline"
          @click="handleAllow"
        >
          {{ t('sandbox.allow') }}
        </t-button>
        <t-button
          theme="primary"
          @click="handleAllowAndRemember"
        >
          {{ t('sandbox.allowAndRemember') }}
        </t-button>
      </div>
    </div>
  </t-dialog>
</template>

<style scoped>
</style>
