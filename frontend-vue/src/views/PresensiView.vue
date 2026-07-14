<script setup>
import { computed, ref, onMounted, onUnmounted } from "vue";
import { watch } from "vue";
import dayjs from "dayjs";
import PresensiModal from "@/components/PresensiModal.vue"; // ⬅ pastikan path benar
import * as presensiService from "@/services/presensiService";
import { usePresensi } from "@/composables/usePresensi";

// =========================
// STATE DARI COMPOSABLE
// =========================
const { presensiList, loading, fetchPresensi, selectedDate } = usePresensi();

// 🔥 set default saat awal load
selectedDate.value = dayjs().format("YYYY-MM-DD");

// state untuk edit data di modal
const selectedPresensi = ref(null);

// state modal dan MODAL CONTROL
const createModal = ref(false);
const editModal = ref(false);
// =========================

// state CRUD
// 🔥 CREATE
const handleCreate = async ({ payload, onError, onSuccess }) => {
  try {
    // 🔥 kirim ke backend
    await presensiService.createPresensi(payload);
    await fetchPresensi(); // 🔥 refresh data table

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

// 🔥 UPDATE
// NOTE:
// - await penting untuk memastikan request selesai sebelum refresh
// - fetch ulang untuk sinkronisasi data backend
const handleEdit = async ({ payload, onError, onSuccess }) => {
  // console log untuk debug payload
  // console.log("PAYLOAD FROM MODAL: ", payload)

  try {
    await presensiService.createPresensi(payload); // param endpoint go: id, model

    await fetchPresensi();
    // 🔥 tidak reset filter → UX konsisten

    onSuccess();
  } catch (err) {
    const message =
      err.response?.data?.error ||
      "Gagal update data";

    onError(message);
  }
};

// === MODAL CONTROL ===
// 🔥 Edit → clone object agar tidak langsung update UI sebelum save
const openEditModal = (item) => {
    selectedPresensi.value = { ...item }; // clone (best practice)
    editModal.value = true;
};

// STATE LOCAL
// =========================
const now = ref(dayjs().format("HH:mm:ss"));
let interval;

// 🔍 input search
const search = ref("");

// =========================
// COMPUTED: FILTER DATA
// =========================
const filteredPresensi = computed(() => {
  return presensiList.value.filter((item) => {

    if (!search.value) return true;

    const keyword = search.value.toLowerCase();

    return (
      item.karyawan_nama?.toLowerCase().includes(keyword) ||
      String(item.karyawan_id).includes(keyword)
    );
  });
});

// =========================
// LIFECYCLE
// =========================
onMounted(() => {
  // 🔥 langsung load berdasarkan tanggal hari ini
  fetchPresensi();
  // 🔥 tidak pakai parameter
  // 👉 karena fetchPresensi sudah baca selectedDate
  interval = setInterval(() => now.value = dayjs().format("HH:mm:ss"), 1000);
});

onUnmounted(() => clearInterval(interval));

// 🔥 reactive → setiap user ganti tanggal, data langsung refresh
watch(selectedDate, () => {
    fetchPresensi();
    // 🔥 setiap tanggal berubah → auto reload
});

// =========================
// HEADER TABLE
// =========================
const headers = [
  { title: "No", key: "no", sortable: false },
  { title: "ID Karyawan", key: "karyawan_id" },
  { title: "Nama", key: "karyawan_nama" },
  { title: "Tanggal", key: "tanggal" },
  { title: "Waktu Masuk", key: "waktu_masuk" },
  { title: "Status", key: "terlambat" },
  { title: "Waktu Pulang", key: "waktu_pulang" },
  { title: 'Actions', key: 'actions', sortable: false },
];

// =========================
// PAGINATION (biar konsisten)
// =========================
const options = ref({
  page: 1,
  itemsPerPage: 5,
});
</script>

<template>
  <v-container fluid>
    <v-card elevation="6" rounded="xl">

      <!-- === HEADER (SAMA SEPERTI KARYAWAN) === -->
      <v-card-title class="d-flex flex-column gap-2">

        <!-- TITLE -->
        <div class="text-h4 font-weight-bold d-flex align-center gap-2">
          <v-icon size="28">mdi-calendar-check</v-icon>
          Presensi Harian
        </div>

        <!-- ACTION -->
        <div class="d-flex align-center gap-2">

            <!-- 🔥 DATE FILTER -->
            <!-- simple & clean (native date input) -->
            <v-text-field
                v-model="selectedDate"
                label="Filter Tanggal"
                type="date"
                density="compact"
                variant="outlined"
                hide-details
                style="max-width: 180px"
            />
            
            <!-- SEARCH -->
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
            <v-btn color="primary" @click="createModal = true">
                <v-icon start>mdi-plus</v-icon>
                Tambah
            </v-btn>

        </div>
      </v-card-title>

      <v-divider class="mb-2" />

      <!-- === TABLE === -->
      <v-data-table
        v-model:options="options"
        :headers="headers"
        :items="filteredPresensi"
        :loading="loading"
        :items-per-page-options="[5, 10, 25, 50, { title: 'All', value: -1 }]"
        class="modern-table"
        hover
        density="comfortable"
        rounded="lg"
      >

        <!-- Loading -->
        <template #loading>
          <v-skeleton-loader type="table-row@5" />
        </template>

        <!-- NOMOR -->
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
          </div>
        </template>

        <!-- HANDLE NULL -->
        <template #item.waktu_masuk="{ item }">
          {{ item.waktu_masuk || '-' }}
        </template>

        <template #item.waktu_pulang="{ item }">
          {{ item.waktu_pulang || '-' }}
        </template>
        
        <template #item.terlambat="{ item }">
        <v-chip
            :color="item.terlambat ? 'red' : 'green'"
            size="small"
            variant="flat"
        >
            {{ item.terlambat ? 'Terlambat' : 'OK' }}
        </v-chip>
        
        <!-- <v-chip
            :color="item.terlambat ? 'red' : 'green'"
            size="small"
            variant="tonal"
        >
            <v-icon start size="16">
            {{ item.terlambat ? 'mdi-alert-circle' : 'mdi-check-circle' }}
            </v-icon>

            {{ item.terlambat ? 'Terlambat' : 'Tepat Waktu' }}
        </v-chip> -->

        </template>
      </v-data-table>

    </v-card>

    <!-- === MODAL COMPONENTS === -->
    <!-- 🔥 reusable modal (clean architecture) -->

    <!-- CREATE -->
    <PresensiModal
        v-model="createModal"
        :mode="'create'"
        :Presensi="null"
        @save="handleCreate"
    />

    <!-- EDIT -->
    <PresensiModal
        v-model="editModal"
        :mode="'edit'"
        :Presensi="selectedPresensi"
        @save="handleEdit"
    />

  </v-container>
</template>