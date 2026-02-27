<script setup lang="ts">
import Sidebar from './Sidebar.vue'
</script>

<template>
  <t-layout class="h-screen w-screen overflow-hidden bg-bg-primary">
    <Sidebar />
    <t-layout class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Shell container keeps content centered like TDesign console -->
      <div class="flex-1 flex flex-col min-w-0 overflow-hidden px-6 py-4">
        <div class="flex-1 flex flex-col min-w-0 overflow-hidden rounded-lg bg-bg-card border border-border shadow-md">
          <slot />
        </div>
      </div>
    </t-layout>
  </t-layout>
</template>
