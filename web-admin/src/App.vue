<script setup lang="ts">
import { ConfigProvider } from 'ant-design-vue'
import { theme } from './config/theme'
</script>

<template>
  <ConfigProvider :theme="theme">
    <router-view />
  </ConfigProvider>
</template>

<style>
#app {
  height: 100vh;
}
</style>
