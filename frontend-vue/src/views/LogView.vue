<script setup>
import { computed, ref, onMounted, onUnmounted } from "vue";
import dayjs from "dayjs";
import * as logService from "@/services/logService";
import { useLog } from "@/composables/useLog";

// === STATE MODAL & DATA TERPILIH ===

// === COMPOSABLE (DATA SOURCE) ===
// 🔥 useLog → centralized state (clean architecture)
const { logList, loading, fetchLogs } = useLog();

// === SEARCH ===
// 🔥 reactive input search
const search = ref("");

// === REALTIME CLOCK (UI tambahan) ===
const now = ref(dayjs().format("HH:mm:ss"));
let interval;

// 🔥 Lifecycle mount
onMounted(() => {
    fetchLogs(); // ambil data awal

    // 🔥 update jam tiap detik
    interval = setInterval(() => {
        now.value = dayjs().format("HH:mm:ss");
    }, 1000);
});

// 🔥 Cleanup interval (penting untuk mencegah memory leak)
onUnmounted(() => clearInterval(interval));

// === TABLE HEADERS ===
// 🔥 key harus sesuai dengan data item
const headers = [
    { title: 'Username', key: 'user_id' },
    { title: 'Method', key: 'method' },
    { title: 'Endpoint', key: 'endpoint' },
    { title: 'IP Address', key: 'ip_address' },
    { title: 'Timestamp', key: 'created_at' },
];

// === PAGINATION ===
// 🔥 v-data-table options (controlled pagination)
const options = ref({
    page: 1,
    itemsPerPage: 5,
});

// === FILTER SEARCH (CLIENT SIDE) ===
// 🔥 computed → reactive & efisien
const filteredLogs = computed(() => {
    if (!search.value) return logList.value;

    const keyword = search.value.toLowerCase();

    return logList.value.filter(log => {
        const username = (log.user_id || '').toLowerCase(); // sesuai nama field di db (user_id)
        const method = (log.method || '').toLowerCase(); // memastikan tidak ada perbedaan GET dengan input search get
        const endpoint = (log.endpoint || '').toLowerCase();
        const ipAddress = (log.ip_address || '').toLowerCase();
        const timestamp = (log.created_at || '').toLowerCase();

        return (
            username.includes(keyword) ||
            method.includes(keyword) ||
            endpoint.includes(keyword) ||
            ipAddress.includes(keyword) ||
            timestamp.includes(keyword)
        );
    });
});
</script>

<template>
    <v-container fluid> <!--fluid membuat jadi lebih lebar tabelnya-->
        <v-card elevation="6" rounded="xl">

            <!-- === HEADER (2 BARIS - BEST PRACTICE UI) === -->
            <v-card-title class="d-flex flex-column gap-2">

                <!-- BARIS 1: TITLE -->
                <!-- 🔥 dipisah supaya tidak cramped -->
                <div class="text-h3 font-weight-bold d-flex gap-2">
                    <v-icon size="28">mdi-history</v-icon>
                    System Log
                </div>

                <!-- BARIS 2: ACTION -->
                <div class="d-flex align-center gap-2">

                    <!-- SEARCH -->
                    <!-- 🔥 flex-grow → otomatis melebar -->
                    <!-- clearable → UX lebih baik -->
                    <v-text-field
                        v-model="search"
                        placeholder="Cari data log ..."
                        prepend-inner-icon="mdi-magnify"
                        density="compact"
                        variant="outlined"
                        hide-details
                        clearable
                        class="flex-grow-1"
                        autocomplete="off"
                    />
                </div>
            </v-card-title>

            <v-divider />

            <!-- === TABLE === -->
            <v-data-table
                v-model:options="options"
                :headers="headers"
                :items="filteredLogs"
                :loading="loading"
                :items-per-page-options="[5, 10, 25, 50, { title: 'All', value: -1 }]"
                class="modern-table"
                hover
                density="comfortable"
                rounded="lg"
            >
                <!-- Loading -->
                <!-- 🔥 skeleton loader → lebih smooth daripada spinner -->
                <template #loading>
                    <v-skeleton-loader type="table-row@5" />
                </template>

                <!-- No (Auto numbering) -->
                <!-- 🔥 support pagination -->
                <template #item.no="{ index }">
                <v-chip size="small" color="primary" variant="tonal">
                    {{
                    options.itemsPerPage === -1
                        ? index + 1
                        : index + 1 + (options.page - 1) * options.itemsPerPage
                    }}
                </v-chip>
                </template>
            </v-data-table>
        </v-card>
    </v-container>
</template>