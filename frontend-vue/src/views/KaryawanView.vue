<script setup>
import { computed, ref, onMounted, onUnmounted } from "vue";
import dayjs from "dayjs";
import KaryawanModal from "@/components/KaryawanModal.vue"; // ⬅ pastikan path benar
import ImportKaryawanModal from "@/components/ImportKaryawanModal.vue"; // ⬅ pastikan path benar
import * as karyawanService from "@/services/karyawanService";
import { useKaryawan } from "@/composables/useKaryawan";
import { useSnackbar } from "@/composables/useSnackbar";
// import useAuth untuk Role Based Access
import { useAuth } from "@/composables/useAuth";

// === STATE MODAL & DATA TERPILIH ===
// 🔥 selectedKaryawan → data yang diklik (edit/delete)
// NOTE: selalu di-clone supaya tidak langsung mutate table (best practice)
const selectedKaryawan = ref(null);

const editModal = ref(false);
const deleteModal = ref(false);
const createModal = ref(false);
const importModal = ref(false);

const { openSnackbar } = useSnackbar();

// === COMPOSABLE (DATA SOURCE) ===
// 🔥 useKaryawan → centralized state (clean architecture)
const { karyawanList, loading, fetchKaryawan } = useKaryawan();

// ambil helper
const { isAdmin } = useAuth(); // ambil role user yang login

// === HANDLER CRUD ===
// 🔥 UPDATE
// NOTE:
// - await penting untuk memastikan request selesai sebelum refresh
// - fetch ulang untuk sinkronisasi data backend
const handleEdit = async ({ payload, onError, onSuccess }) => {
  try {
    await karyawanService.updateKaryawan(payload.id, payload);
    await fetchKaryawan();

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
        await karyawanService.deleteKaryawan(data.id);
        await fetchKaryawan();

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
    await karyawanService.createKaryawan(payload);
    await fetchKaryawan(); // 🔥 refresh data table

    onSuccess(); // 🔥 kasih tahu modal: sukses
  } catch (err) {
    // 🔥 ambil error dari backend (gin response)
    const message =
      err.response?.data?.error || // 🔥 spesifik (id sudah ada)
      err.response?.data?.message || // 🔥 fallback
      "Terjadi kesalahan";

    onError(message); // 🔥 kirim balik ke modal
  }
};

// handle sukses import
const handleImportSuccess = async () => {
  await fetchKaryawan(); // 🔥 refresh table setelah import
};

// === MODAL CONTROL ===

// 🔥 Edit → clone object agar tidak langsung update UI sebelum save
const openEditModal = (item) => {
    selectedKaryawan.value = { ...item };
    editModal.value = true;
};

// 🔥 Delete → cukup clone juga untuk safety
const openDeleteModal = (item) => {
    selectedKaryawan.value = { ...item };
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
    fetchKaryawan(); // ambil data awal

    // 🔥 update jam tiap detik
    interval = setInterval(() => {
        now.value = dayjs().format("HH:mm:ss");
    }, 1000);
});

// 🔥 Cleanup interval (penting untuk mencegah memory leak)
onUnmounted(() => clearInterval(interval));

// === TABLE HEADERS ===
// const headers = [
//     { title: 'No', key: 'no', sortable: false },
//     { title: 'ID', key: 'id' },
//     { title: 'Nama', key: 'nama' },
//     { title: 'Status', key: 'aktif' },
//     { title: 'Actions', key: 'actions', sortable: false },
// ];
// 🔥 key harus sesuai dengan data item
// 🔥 dynamic headers berdasarkan role
const headers = computed(() => {

    // 🔥 header default
    const baseHeaders = [
        { title: 'No', key: 'no', sortable: false },
        { title: 'ID', key: 'id' },
        { title: 'Nama', key: 'nama' },
        { title: 'Status', key: 'aktif' },
    ];

    // 🔥 hanya admin bisa lihat kolom actions
    if (isAdmin()) { // isAdmin() adalah pemanggilan function isAdmin yang return value nya true or false
        baseHeaders.push({
            title:"Actions",
            key:"actions",
            sortable:false,
        });
    }

    return baseHeaders;
});

// === PAGINATION ===
// 🔥 v-data-table options (controlled pagination)
const options = ref({
    page: 1,
    itemsPerPage: 5,
});

// === FILTER SEARCH (CLIENT SIDE) ===
// 🔥 computed → reactive & efisien
const filteredKaryawan = computed(() => {
    if (!search.value) return karyawanList.value;

    const keyword = search.value.toLowerCase();

    return karyawanList.value.filter(k => {
        const nama = (k.nama || '').toLowerCase();

        // 🔥 penting: convert id ke string untuk bisa di-search
        const id = String(k.id || '');

        return (
        nama.includes(keyword) ||
        id.includes(keyword)
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
                    <v-icon size="28">mdi-account-group</v-icon>
                    Data Karyawan
                </div>

                <!-- BARIS 2: ACTION -->
                <div class="d-flex align-center gap-2">
                    <!-- SEARCH -->
                    <!-- 🔥 flex-grow → otomatis melebar -->
                    <!-- clearable → UX lebih baik -->
                    <v-text-field
                        v-model="search"
                        placeholder="Cari ID / Nama Karyawan..."
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
                    <!-- tampilkan button jika user role adalah admin -->
                    <v-btn
                        color="primary"
                        @click="createModal = true">
                        <v-icon start>mdi-plus</v-icon>
                        Tambah
                    </v-btn>

                    <!-- BUTTON IMPORT FROM EXCEL -->
                    <v-btn color="success" @click="importModal = true">
                        <v-icon start>mdi-download</v-icon>
                        Import Data
                    </v-btn>
                </div>
            </v-card-title>

            <v-divider />

            <!-- === TABLE === -->
            <v-data-table
                v-model:options="options"
                :headers="headers"
                :items="filteredKaryawan"
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
                <!-- tidak perlu pakai v-if lagi karena rba untuk kolom actions sudah diatur di headers (lihat bagian headers) -->
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

                <!-- STATUS -->
                <!-- 🔥 visual indicator -->
                <template #item.aktif="{ item }">
                <v-chip
                    :color="item.aktif ? 'green' : 'red'"
                    size="small"
                    variant="flat"
                >
                    {{ item.aktif ? 'Aktif' : 'Nonaktif' }}
                </v-chip>
                </template>

            </v-data-table>

            </v-card>

            <!-- === MODAL COMPONENTS === -->
            <!-- 🔥 reusable modal (clean architecture) -->

            <!-- CREATE -->
            <KaryawanModal
                v-model="createModal"
                :mode="'create'"
                :karyawan="null"
                @save="handleCreate"
            />

            <!-- EDIT -->
            <KaryawanModal
                v-model="editModal"
                :mode="'edit'"
                :karyawan="selectedKaryawan"
                @save="handleEdit"
            />

            <!-- DELETE -->
            <KaryawanModal
                v-model="deleteModal"
                :mode="'delete'"
                :karyawan="selectedKaryawan"
                @delete="handleDelete"
            />

            <!-- IMPORT -->
            <ImportKaryawanModal
                v-model="importModal"
                @success="handleImportSuccess"
            />
    </v-container>
</template>