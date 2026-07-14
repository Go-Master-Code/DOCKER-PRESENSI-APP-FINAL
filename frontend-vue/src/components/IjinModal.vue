<script setup>
import { ref, watch, nextTick } from "vue";

// import service karyawan dan jenis ijin untuk get data dan dimasukkan ke dalam combobox
import * as karyawanService from "@/services/karyawanService";
import * as jenisIjinService from "@/services/jenisIjinService";

// 🔥 LIST DATA UNTUK DROPDOWN
const karyawanOptions = ref([]); // 🔥 state untuk menampung data dropdown karyawan
const jenisIjinOptions = ref([]); // 🔥 state untuk menampung data dropdown jenis ijin

// 🔥 LOADING STATE (UX penting)
// 🔥 loading indicator (UX penting saat fetch API)
const loadingKaryawan = ref(false);
const loadingJenisIjin = ref(false);

// 🔥 Ref untuk akses langsung ke komponen input (dipakai untuk autofocus)
const inputTanggalRef = ref(null);
const inputKeteranganRef = ref(null);

// function datepicker pada modal untuk select today's date langsung
const getToday = () => {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0");
  const day = String(now.getDate()).padStart(2, "0");

  return `${year}-${month}-${day}`;
};

// 🔥 Props dari parent (IjinView)
// modelValue → untuk buka/tutup modal (v-model)
// Ijin → data yang dipilih (null saat create)
// mode → menentukan behavior (create | edit | delete)
const props = defineProps({
    modelValue: Boolean,
    Ijin: { type: Object, default: () => null },
    mode: { type: String, default: "edit" },
});

// 🔥 Emit event ke parent
// update:modelValue → untuk close modal
// save → create/update
// delete → hapus data
const emit = defineEmits(["update:modelValue", "save", "delete"]);

// === Reactive form (state lokal, tidak langsung mengubah props) ===
const formTanggal = ref("");
const formKeterangan = ref("");
const formKaryawanID = ref(null);
const formJenisIjinID = ref(null);

// 🔥 Utility: hanya angka (sanitize input)
// penting karena user bisa paste string
// const onlyNumber = (value) => {
//   return value.replace(/\D/g, ""); // hapus semua selain angka
// };

// === Sync props ke form saat modal dibuka ===
// 🔥 Kenapa watch modelValue?
// karena props.Ijin tidak selalu berubah (misal create = null terus)
// jadi trigger terbaik adalah saat modal dibuka
watch(
  () => props.modelValue, // 🔥 trigger saat modal open/close
  async (isOpen) => {
    if (!isOpen) return; // 🔥 hanya jalan saat modal dibuka

    // get master data dan masukkan ke dalam combobox
    await fetchDropdownData();

    // 🔥 RESET ERROR SETIAP MODAL DIBUKA
    errors.value = {
        tanggal: "",
        keterangan: "",
        karyawan_id: "",
        jenis_ijin_id: "",
    };

    if (props.mode === "create") {
      // 🔥 Reset form saat create
      formTanggal.value = getToday(); // ini akan trigger watcher, await fetchDropdownData() di atas tidak perlu dilakukan (double api call)
      formKeterangan.value = "";
      formKaryawanID.value = null; // karena bentuknya autocomplete
      formJenisIjinID.value = null; // karena bentuknya select

      // fetchDropdownData memfilter karyawan yang sudah ada di tabel ijin di tgl hari ini
      fetchDropdownData();

      // 🔥 Autofocus ke Tanggal (UX lebih cepat)
      await nextTick(); // tunggu DOM render
      inputTanggalRef.value?.focus();
    } else {
      // 🔥 Isi form saat edit
      formTanggal.value = props.Ijin?.tanggal || "";
      formKeterangan.value = props.Ijin?.keterangan || "";
      formKaryawanID.value = props.Ijin?.karyawan_id || null;
      formJenisIjinID.value = props.Ijin?.jenis_ijin_id || null;

      // 2. fetch dropdown
      await fetchDropdownData();

      // 🔥 Autofocus ke Nama (karena ID tidak boleh diubah)
      await nextTick();
      inputTanggalRef.value?.focus();
    }
  }
);

const fetchDropdownData = async () => {
  if (!formTanggal.value) return; // guard

  try {
    // 🔥 aktifkan loading
    loadingKaryawan.value = true;
    loadingJenisIjin.value = true;

    // 🔥 ambil 2 API sekaligus (lebih cepat dari sequential)
    const [resKaryawan, resJenis] = await Promise.all([
      karyawanService.getAllKaryawanBelumIjin(formTanggal.value),
      jenisIjinService.getAllJenisIjin(),
    ]);

    // 🔥 mapping data ke format Vuetify (WAJIB)
    // Vuetify butuh: { title, value }
    karyawanOptions.value = resKaryawan.data.data.map(k => ({
      title: k.nama,   // 🔥 teks yang ditampilkan ke user
      value: k.id,     // 🔥 nilai yang dikirim ke backend
    }));

    jenisIjinOptions.value = resJenis.data.data.map(j => ({
      title: j.nama,   // 🔥 teks yang ditampilkan ke user
      value: j.id,     // 🔥 nilai yang dikirim ke backend
    }));

  } catch (err) {
    console.error(err); // 🔥 debug kalau API error
  } finally {
    // 🔥 matikan loading
    loadingKaryawan.value = false;
    loadingJenisIjin.value = false;
  }
};

// 🔥 UX penting:
// Error hilang otomatis saat user mulai memperbaiki input

watch(formTanggal, async (newTanggal) => { // reactive mengubah dropdown karyawan sesuai tanggal yang dipilih
  if (!newTanggal) return;

  await fetchDropdownData(); // single source of truth
});

watch(formKeterangan, () => {
    errors.value.keterangan = "";
});

watch(formKaryawanID, () => {
    errors.value.karyawan_id = "";
});

watch(formJenisIjinID, () => {
    errors.value.jenis_ijin_id = "";
});

// === Close modal ===
// 🔥 gunakan emit agar parent yang kontrol state
const closeModal = () => emit("update:modelValue", false);

// === Save handler (create & edit) ===
const save = () => {
  emit("save", {
    payload: {
      id: props.Ijin?.id, // 🔥 id harus ada pada param di endpoint, hasilnya misal: localhost:8080/api/ijin_karyawan/6, 6 diambil dari id ini
      tanggal: formTanggal.value, // 🔥 dari date picker
      karyawan_id: formKaryawanID.value, // 🔥 dari autocomplete
      jenis_ijin_id: formJenisIjinID.value, // 🔥 dari select
      keterangan: formKeterangan.value,
    },

    // 🔥 callback error dari parent
    onError: (errMsg) => {
       // 🔥 mapping error ke field yang sesuai
      const msg = errMsg.toLowerCase();

      if (msg.includes("karyawan")) {
        errors.value.karyawan_id = errMsg;
      } else if (msg.includes("jenis")) {
        errors.value.jenis_ijin_id = errMsg;
      } else if (msg.includes("tanggal")) {
        errors.value.tanggal = errMsg;
      } else {
        errors.value.keterangan = errMsg;
      }
    },

    // 🔥 callback success → tutup modal
    onSuccess: () => closeModal(),
  });
};

// === Delete handler ===
const remove = () => {
  emit("delete", props.Ijin);
  closeModal();
};

// === ERROR STATE (INLINE VALIDATION) ===
// 🔥 simpan error per field (biar bisa tampil di bawah input)
const errors = ref({
  tanggal: "",
  keterangan: "",
  karyawan_id: "",
  jenis_ijin_id: "",
});
</script>

<template>
  <!-- 🔥 v-model ke props.modelValue → modal dikontrol parent -->
  <v-dialog v-model="props.modelValue" max-width="400">
    
    <!-- 🔥 Tangkap enter global di dalam modal -->
    <v-card>

      <!-- TITLE -->
      <v-card-title class="text-h6">
        {{
          props.mode === "create"
            ? "Tambah Ijin Karyawan"
            : props.mode === "edit"
            ? "Edit Ijin Karyawan"
            : "Hapus Ijin Karyawan"
        }}
      </v-card-title>

      <!-- ================================================= -->
      <!-- FORM -->
      <!-- ================================================= -->
      <v-form
        v-if="props.mode !== 'delete'"
        @submit.prevent="save"
      >
      
        <v-card-text>

          <!-- DATE -->
          <v-text-field
            ref="inputTanggalRef"
            label="Tanggal"
            v-model="formTanggal"
            type="date"
            class="mb-5"
            density="compact"
            variant="outlined"
            hide-details="auto"
            style="max-width:180px"

            :error="!!errors.tanggal"
            :error-messages="errors.tanggal"
          />

          <!-- KARYAWAN -->
          <v-autocomplete
            v-model="formKaryawanID"
            :items="karyawanOptions"
            item-title="title"
            item-value="value"
            label="Pilih Karyawan"
            placeholder="Cari nama karyawan..."
            density="compact"
            variant="outlined"
            :loading="loadingKaryawan"
            autocomplete="off"
            clearable

            :error="!!errors.karyawan_id"
            :error-messages="errors.karyawan_id"
          />

          <!-- JENIS IJIN -->
          <v-select
            v-model="formJenisIjinID"
            :items="jenisIjinOptions"
            item-title="title"
            item-value="value"
            label="Jenis Ijin"
            density="compact"
            variant="outlined"
            :loading="loadingJenisIjin"
            clearable

            :error="!!errors.jenis_ijin_id"
            :error-messages="errors.jenis_ijin_id"
          />

          <!-- KETERANGAN -->
          <v-text-field
            ref="inputKeteranganRef"
            label="Keterangan"
            v-model="formKeterangan"
            density="compact"
            variant="outlined"
            autocomplete="off"
            maxlength="100"
            counter

            :error="!!errors.keterangan"
            :error-messages="errors.keterangan"
          />

        </v-card-text>

        <!-- ACTION -->
        <v-card-actions>

          <v-spacer />

          <v-btn
            color="red"
            variant="tonal"
            size="small"
            @click="closeModal"
          >
            Batal
          </v-btn>

          <v-btn
            type="submit"
            color="primary"
            variant="elevated"
            size="small"
            :disabled="
              !formTanggal ||
              !formKaryawanID ||
              !formJenisIjinID ||
              !formKeterangan ||
              formKeterangan.length > 100
            "
          >
            Simpan
          </v-btn>

        </v-card-actions>

      </v-form>

      <!-- ================================================= -->
      <!-- DELETE MODE -->
      <!-- ================================================= -->
      <template v-else>

        <v-card-text>
          Apakah yakin ingin menghapus ijin
          <strong>"{{ props.Ijin?.keterangan }}"</strong>
          dari
          <strong>{{ props.Ijin?.karyawan_nama }}</strong>
          tanggal
          <strong>{{ props.Ijin?.tanggal }}</strong>?
        </v-card-text>

        <v-card-actions>

          <v-spacer />

          <v-btn
            color="secondary"
            variant="tonal"
            size="small"
            @click="closeModal"
          >
            Batal
          </v-btn>

          <v-btn
            color="red"
            size="small"
            variant="elevated"
            @click="remove"
          >
            Hapus
          </v-btn>

        </v-card-actions>

      </template>

    </v-card>
  </v-dialog>
</template>