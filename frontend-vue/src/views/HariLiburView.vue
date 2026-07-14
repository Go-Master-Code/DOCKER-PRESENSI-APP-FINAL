<script setup>
import { computed, ref, onMounted, onUnmounted } from "vue";
import dayjs from "dayjs";
import HariLiburModal from "@/components/HariLiburModal.vue"; // ⬅ pastikan path benar
import * as hariLiburService from "@/services/hariLiburService";
import { useHariLibur } from "@/composables/useHariLibur";
import { useSnackbar } from "@/composables/useSnackbar";

// === STATE MODAL & DATA TERPILIH ===
// 🔥 selectedHariLibur → data yang diklik (edit/delete)
// NOTE: selalu di-clone supaya tidak langsung mutate table (best practice)
const selectedHariLibur = ref(null);

const editModal = ref(false);
const deleteModal = ref(false);
const createModal = ref(false);

const { openSnackbar } = useSnackbar();

// === COMPOSABLE (DATA SOURCE) ===
// 🔥 useHariLibur → centralized state (clean architecture)
const { hariLiburList, loading, fetchHariLibur } = useHariLibur();

// === HANDLER CRUD ===

// 🔥 UPDATE
// NOTE:
// - await penting untuk memastikan request selesai sebelum refresh
// - fetch ulang untuk sinkronisasi data backend
const handleEdit = async ({ payload, onError, onSuccess }) => {
  try {
    await hariLiburService.updateHariLibur(payload.id, payload);
    await fetchHariLibur();

    onSuccess();
  } catch (err) {
    const message =
      err.response?.data?.error ||
      "Gagal update data";

    onError(message);
  }
};

// 🔥 DELETE
const handleDelete = async (data) => {
    try {
        await hariLiburService.deleteHariLibur(data.id);
        await fetchHariLibur();

        openSnackbar("Data berhasil dihapus", "success");
    } catch (err) {
        console.error(err);
        openSnackbar("Gagal hapus data", "error");
    }
};

// 🔥 CREATE
const handleCreate = async ({ payload, onError, onSuccess }) => {
  try {
    // 🔥 kirim ke backend
    await hariLiburService.createHariLibur(payload);
    await fetchHariLibur(); // 🔥 refresh data table

    onSuccess(); // 🔥 kasih tahu modal: sukses
  } catch (err) {
    // 🔥 ambil error dari backend (gin response)
    const message =
      err.response?.data?.error || // 🔥 spesifik (tanggal sudah ada)
      err.response?.data?.message || // 🔥 fallback
      "Terjadi kesalahan";

    onError(message); // 🔥 kirim balik ke modal
  }
};

// === MODAL CONTROL ===

// 🔥 Edit → clone object agar tidak langsung update UI sebelum save
const openEditModal = (item) => {
    selectedHariLibur.value = { ...item };
    editModal.value = true;
};

// 🔥 Delete → cukup clone juga untuk safety
const openDeleteModal = (item) => {
    selectedHariLibur.value = { ...item };
    deleteModal.value = true;
};

// === SEARCH ===
// 🔥 reactive input search
const search = ref("");

// === REALTIME CLOCK (UI tambahan) ===
const now = ref(dayjs().format("HH:mm:ss"));
let interval;

// 🔥 Lifecycle mount
onMounted(() => {
    fetchHariLibur(); // ambil data awal

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
    { title: 'No', key: 'no', sortable: false },
    // { title: 'ID', key: 'id' }, tidak ditampilkan agar ga tabrakan dengan penomoran row
    { title: 'Tanggal', key: 'tanggal' },
    { title: 'Keterangan', key: 'keterangan' },
    { title: 'Actions', key: 'actions', sortable: false },
];

// === PAGINATION ===
// 🔥 v-data-table options (controlled pagination)
const options = ref({
    page: 1,
    itemsPerPage: 5,
});

// === FILTER SEARCH (CLIENT SIDE) ===
// 🔥 computed → reactive & efisien
const filteredHariLibur = computed(() => {
    if (!search.value) return hariLiburList.value;

    const keyword = search.value.toLowerCase();

    return hariLiburList.value.filter(hl => {
        const keterangan = (hl.keterangan || '').toLowerCase();

        // 🔥 penting: convert tanggal ke string untuk bisa di-search
        const tanggal = String(hl.tanggal || '');

        return (
            keterangan.includes(keyword) ||
            tanggal.includes(keyword)
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
                <v-icon size="28">mdi-calendar-remove</v-icon>
                Data Hari Libur
            </div>

            <!-- BARIS 2: ACTION -->
            <div class="d-flex align-center gap-2">

            <!-- SEARCH -->
            <!-- 🔥 flex-grow → otomatis melebar -->
            <!-- clearable → UX lebih baik -->
            <v-text-field
                v-model="search"
                placeholder="Cari Tanggal / Nama Hari Libur..."
                prepend-inner-icon="mdi-magnify"
                density="compact"
                variant="outlined"
                hide-details
                clearable
                class="flex-grow-1"
                autocomplete="off"
            />

            <!-- BUTTON TAMBAH -->
            <!-- 🔥 icon + text → lebih mudah dikenali -->
            <v-btn color="primary" @click="createModal = true">
                <v-icon start>mdi-plus</v-icon>
                Tambah
            </v-btn>

            </div>
        </v-card-title>

        <v-divider />

        <!-- === TABLE === -->
        <v-data-table
            v-model:options="options"
            :headers="headers"
            :items="filteredHariLibur"
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

            <!-- ACTIONS -->
            <!-- 🔥 tombol dengan icon (UX lebih baik) -->
            <template #item.actions="{ item }">
            <div class="d-flex gap-2">
                <v-btn size="small" color="primary" @click="openEditModal(item)">
                <v-icon start>mdi-square-edit-outline</v-icon>
                Edit
                </v-btn>

                <v-btn
                size="small"
                color="red"
                variant="tonal"
                @click="openDeleteModal(item)"
                >
                <v-icon start>mdi-delete-outline</v-icon>
                    Delete
                </v-btn>
            </div>
            </template>

        </v-data-table>

        </v-card>

        <!-- === MODAL COMPONENTS === -->
        <!-- 🔥 reusable modal (clean architecture) -->

        <!-- CREATE -->
        <HariLiburModal
            v-model="createModal"
            :mode="'create'"
            :hariLibur="null"
            @save="handleCreate"
        />

        <!-- EDIT -->
        <HariLiburModal
            v-model="editModal"
            :mode="'edit'"
            :hariLibur="selectedHariLibur"
            @save="handleEdit"
        />

        <!-- DELETE -->
        <HariLiburModal
            v-model="deleteModal"
            :mode="'delete'"
            :hariLibur="selectedHariLibur"
            @delete="handleDelete"
        />

    </v-container>
</template>