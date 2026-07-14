<script setup>
import { onMounted } from "vue";
import { useAppTheme } from "@/composables/useTheme";
// 2 import di bawah untuk session expired
import SessionExpiredDialog
from "@/components/SessionExpiredDialog.vue";

// composable
const { loadTheme } = useAppTheme();

// load saat app start
onMounted(() => {
    loadTheme();
});
</script>

<template>
  <v-app>
    <router-view />
    <!-- 🔥 global session expired -->
        <SessionExpiredDialog />
  </v-app>
</template>