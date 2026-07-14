<script setup>
import { computed, ref, onMounted, onUnmounted } from "vue";
import dayjs from "dayjs";
import IjinModal from "@/components/IjinModal.vue"; // ⬅ pastikan path benar
import * as ijinService from "@/services/ijinService";
import { watch } from "vue";
import { useIjin } from "@/composables/useIjin";
import { useSnackbar } from "@/composables/useSnackbar";

// =========================
// STATE DARI COMPOSABLE
// =========================
const { ijinList, loading, fetchIjin, selectedDate } = useIjin(); // selectedDate ambil dari useIjin.js (sumber tunggal)

// 🔥 set default saat awal load
selectedDate.value = dayjs().format("YYYY-MM-DD");

const selectedIjin = ref(null);

// state modal dan MODAL CONTROL
const editModal = ref(false);
const deleteModal = ref(false);
const createModal = ref(false);

// deklarasi snackbar
const { openSnackbar } = useSnackbar();

// === HANDLER CRUD ===

// 🔥 UPDATE
// NOTE:
// - await penting untuk memastikan request selesai sebelum refresh
// - fetch ulang untuk sinkronisasi data backend
const handleEdit = async ({ payload, onError, onSuccess }) => {
  // console log untuk debug payload
  // console.log("PAYLOAD FROM MODAL: ", payload)

  try {
    await ijinService.updateIjin(payload.id, payload); // param endpoint go: id, model

    await fetchIjin();
    // 🔥 tidak reset filter → UX konsisten

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
        await ijinService.deleteIjin(data.id);

        await fetchIjin();
        // 🔥 tetap sesuai tanggal yang dipilih

        openSnackbar("Data berhasil dihapus", "success");
    } catch (err) {
        console.error(err);
        openSnackbar("Gagal hapus data", "error");
    }
};

// 🔥 CREATE
const handleCreate = async ({ payload, onError, onSuccess }) => {
  // console log untuk debug payload
  // console.log("PAYLOAD FROM MODAL: ", payload)

  try {
    // 🔥 kirim ke backend
    await ijinService.createIjin(payload);

    await fetchIjin(); // 🔥 refresh data table
    // 🔥 tetap pakai selectedDate → filter tetap aktif

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
    selectedIjin.value = { ...item }; // clone (best practice)
    editModal.value = true;
};

// 🔥 Delete → cukup clone juga untuk safety
const openDeleteModal = (item) => {
    selectedIjin.value = { ...item }; // clone (best practice)
    deleteModal.value = true;  // buka modal
};

// =========================
// STATE LOCAL
// =========================
const now = ref(dayjs().format("HH:mm:ss"));
let interval;

// 🔍 input search
const search = ref("");

// =========================
// COMPUTED: FILTER DATA
// =========================
const filteredIjin = computed(() => {
  return ijinList.value.filter((item) => {

    if (!search.value) return true;

    const keyword = search.value.toLowerCase();

    return (
      item.karyawan_nama?.toLowerCase().includes(keyword) ||
      item.keterangan?.toLowerCase().includes(keyword) ||
      String(item.karyawan_id).includes(keyword)
    );
  });
});

// =========================
// LIFECYCLE
// =========================
onMounted(() => {
  // 🔥 langsung load berdasarkan tanggal hari ini
  fetchIjin();
  // 🔥 tidak pakai parameter
  // 👉 karena fetchIjin sudah baca selectedDate
  interval = setInterval(() => now.value = dayjs().format("HH:mm:ss"), 1000);
});

onUnmounted(() => clearInterval(interval));

// 🔥 reactive → setiap user ganti tanggal, data langsung refresh
watch(selectedDate, () => {
    fetchIjin();
    // 🔥 setiap tanggal berubah → auto reload
});

// =========================
// HEADER TABLE
// =========================
const headers = [
  { title: "No", key: "no", sortable: false },
  { title: "Tanggal", key: "tanggal" },
  { title: "ID Karyawan", key: "karyawan_id" },
  { title: "Nama", key: "karyawan_nama" },
  { title: "Jenis Ijin", key: "jenis_ijin_nama" },
  { title: "Keterangan", key: "keterangan" },
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
            <v-icon size="28">mdi-account-clock</v-icon>
            Ijin Karyawan
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
                  placeholder="Cari ID / Nama Karyawan / Keterangan..."
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
          :items="filteredIjin"
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

          <!-- Kolom (No) AUTO NUMBERING -->
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
      <IjinModal
          v-model="createModal"
          :mode="'create'"
          :Ijin="null"
          @save="handleCreate"
      />

      <!-- EDIT -->
      <IjinModal
          v-model="editModal"
          :mode="'edit'"
          :Ijin="selectedIjin"
          @save="handleEdit"
      />

      <!-- DELETE -->
      <IjinModal
          v-model="deleteModal"
          :mode="'delete'"
          :Ijin="selectedIjin"
          @delete="handleDelete"
      />
  </v-container>
</template>