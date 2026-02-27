<script setup lang="ts">
import Sidebar from './Sidebar.vue'
</script>

<template>
  <t-layout class="h-screen w-screen overflow-hidden">
    <Sidebar />
    <t-layout class="flex-1 flex flex-col min-w-0 overflow-hidden bg-bg-primary">
      <slot />
    </t-layout>
  </t-layout>
</template>
