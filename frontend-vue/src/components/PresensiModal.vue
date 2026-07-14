<script setup>
import { ref, watch, nextTick, computed } from "vue";

// import service karyawan dan jenis Presensi untuk get data dan dimasukkan ke dalam combobox
import * as karyawanService from "@/services/karyawanService";

// 🔥 LIST DATA UNTUK DROPDOWN
const karyawanOptions = ref([]); // 🔥 state untuk menampung data dropdown karyawan

// 🔥 LOADING STATE (UX penting)
// 🔥 loading indicator (UX penting saat fetch API)
const loadingKaryawan = ref(false);

// 🔥 Ref untuk akses langsung ke komponen input (dipakai untuk autofocus)
const inputTanggalRef = ref(null);
// ref untuk akses langsung (autofocus) ke datepicker pulang
const inputWaktuPulangRef = ref(null);

// state untuk edit mode (tanggal, combobox karyawan, waktu masuk di disable)
const isEditMode = computed(() => props.mode === "edit");

// state untuk create mode (waktu pulang di disable, karena akan sama valuenya dengan waktu masuk)
const isCreateMode = computed(() => props.mode === "create")

// function datepicker pada modal untuk select today's date langsung
const getToday = () => {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0");
  const day = String(now.getDate()).padStart(2, "0");

  return `${year}-${month}-${day}`;
};

// TIMEPICKER waktu masuk dan
// 🔥 STATE WAKTU (UI: HH:mm, backend: HH:mm:ss)
const formWaktuMasuk = ref("");
const formWaktuPulang = ref("");
// 🔥 kontrol buka/tutup dropdown
const menuMasuk = ref(false);
const menuPulang = ref(false);

// cek / validasi tiap value dari komponen input
const isFormInvalid = computed(() => {
  return (
    !formTanggal.value ||
    !formKaryawanID.value ||
    !formWaktuMasuk.value ||
    (props.mode === "edit" && !formWaktuPulang.value) ||
    isInvalidTime.value
  );
});

// di isFormInvalid => jika mode nya edit dan formWaktuPulang kosong, tidak bisa disubmit form nya

// 🔥 VALIDASI: pulang tidak boleh lebih kecil dari masuk
const isInvalidTime = computed(() => {
  if (!formWaktuMasuk.value || !formWaktuPulang.value) return false;

  return formWaktuPulang.value < formWaktuMasuk.value;
});

// 🔥 FORMAT KE HH:mm:ss untuk timepicker
const formatToHHmmss = (time) => {
  if (!time) return null;
  if (time.length === 8) return time;
  return `${time}:00`;
};

// 🔥 Props dari parent (PresensiView)
// modelValue → untuk buka/tutup modal (v-model)
// Presensi → data yang dipilih (null saat create)
// mode → menentukan behavior (create | edit | delete)
const props = defineProps({
    modelValue: Boolean,
    Presensi: { type: Object, default: () => null }, // var Presensi harus sama dengan yang di passing dari view
    mode: { type: String, default: "edit" },
});

// 🔥 Emit event ke parent
// update:modelValue → untuk close modal
// save → create/update

const emit = defineEmits(["update:modelValue", "save"]);

// === Reactive form (state lokal, tidak langsung mengubah props) ===
const formTanggal = ref("");
const formKaryawanID = ref(null);

// === Sync props ke form saat modal dibuka ===
// 🔥 Kenapa watch modelValue?
// karena props.Presensi tidak selalu berubah (misal create = null terus)
// jadi trigger terbaik adalah saat modal dibuka
watch(
  () => props.modelValue, // 🔥 trigger saat modal open/close
  async (isOpen) => {
    if (!isOpen) return; // 🔥 hanya jalan saat modal dibuka

    // 🔥 RESET ERROR SETIAP MODAL DIBUKA
    errors.value = {
      waktu_masuk: "",
      waktu_pulang: "",
    };

    if (props.mode === "create") {
      // 🔥 Reset form saat create
      // field tanggal di set ke today
      formTanggal.value = getToday(); // ini akan trigger watcher, await fetchDropdownData() di atas tidak perlu dilakukan (double api call)
      formKaryawanID.value = null;
      formWaktuMasuk.value = ""; // time picker
      formWaktuPulang.value = ""; // time picker

      // fetchDropdownData memfilter karyawan yang sudah ada di tabel presensi di tgl hari ini
      fetchDropdownData();

      // 🔥 Autofocus ke Tanggal (UX lebih cepat)
      await nextTick(); // tunggu DOM render
      inputTanggalRef.value?.focus();
    } else {
      // 1. set tanggal dulu
      formTanggal.value = props.Presensi?.tanggal || "";

      // 2. fetch dropdown
      await fetchDropdownData();

      // 3. pastikan karyawan yang diedit ada di options
      const currentKaryawanId = props.Presensi?.karyawan_id;

      const isExist = karyawanOptions.value.some(
        k => k.value === currentKaryawanId
      );

      if (!isExist && props.Presensi) {
        karyawanOptions.value.push({
          title: props.Presensi.karyawan_nama,
          value: currentKaryawanId,
        });
      }

      // 4. set value
      formKaryawanID.value = currentKaryawanId;
      
      // 🔥 PREFILL DARI BACKEND
      // backend: "08:30:00"
      // UI butuh: "08:30"
      formWaktuMasuk.value = props.Presensi?.waktu_masuk
        ? props.Presensi.waktu_masuk.slice(0, 5)
        : "";

      // formWaktuPulang.value = props.Presensi?.waktu_pulang -> picker sesuai waktu pulang row data
      //   ? props.Presensi.waktu_pulang.slice(0, 5)
      //   : "";

      // picker waktu pulang dikosongkan agar tinggal diedit
      formWaktuPulang.value="";

      // 🔥 Autofocus KE TANGGAL
      await nextTick();
      // 🔥 buka time picker otomatis
      menuPulang.value = true;

      // 🔥 optional: tetap fokus ke input
      inputWaktuPulangRef.value?.focus();
    }
  }
);

const fetchDropdownData = async () => {
  if (!formTanggal.value) return; // guard

  try {
    // 🔥 Aktifkan loading state
    // Berguna untuk UX → bisa tampilkan spinner di autocomplete
    loadingKaryawan.value = true;

    // 🔥 Panggil API
    // Mengambil daftar karyawan yang BELUM absen pada tanggal yang dipilih
    // formTanggal.value harus format "YYYY-MM-DD"
    const resKaryawan = await karyawanService.getAllKaryawanBelumAbsen(formTanggal.value);

    // 🔥 SAFE ACCESS (anti error null / undefined)
    // - resKaryawan?.data?.data → optional chaining (hindari crash jika undefined)
    // - ?? [] → jika hasilnya null / undefined, fallback ke array kosong
    // Kenapa penting?
    // Karena backend kadang return:
    // { data: null } → ini akan bikin .map() error kalau tidak di-handle
    const list = resKaryawan?.data?.data ?? [];

    // 🔥 Mapping ke format Vuetify
    // v-autocomplete BUTUH struktur:
    // { title: "yang ditampilkan", value: "yang disimpan" }
    karyawanOptions.value = list.map(k => ({
      title: k.nama,
      value: k.id,
    }));

  } catch (err) {
    console.error(err);
    // 🔥 Fallback UI
    karyawanOptions.value = []; // 🔥 fallback biar UI tidak rusak
    // Set ke array kosong supaya:
    // - tidak crash
    // - autocomplete tetap render normal (meskipun kosong)
  } finally {
    // 🔥 matikan loading
    loadingKaryawan.value = false;
  }
};

// 🔥 UX penting:
// Error hilang otomatis saat user mulai memperbaiki input

watch(formTanggal, async (newTanggal) => { // reactive mengubah dropdown karyawan sesuai tanggal yang dipilih
  if (!newTanggal) return;

  await fetchDropdownData(); // single source of truth
});

watch(formWaktuMasuk, () => {
    errors.value.waktu_masuk = "";
});

watch(formWaktuPulang, () => {
    errors.value.waktu_pulang = "";
});

watch(formKaryawanID, () => {
    errors.value.karyawan_id = "";
});

// === Close modal ===
// 🔥 gunakan emit agar parent yang kontrol state
const closeModal = () => emit("update:modelValue", false);

// === Save handler (create & edit) ===
const save = () => {
  if (isInvalidTime.value) return; // jika waktu invalid (waktu pulang<waktu masuk)
  
  emit("save", {
    payload: {
      tanggal: formTanggal.value, // 🔥 dari date picker
      // 🔥 FORMAT FINAL KE BACKEND
      waktu_masuk: isEditMode.value
        ? props.Presensi?.waktu_masuk // 🔥 lock dari backend
        : formatToHHmmss(formWaktuMasuk.value),

      waktu_pulang: isEditMode.value // jika mode edit, ambil value waktu pulang dari time picker
        ? formatToHHmmss(formWaktuPulang.value)
        : formatToHHmmss(formWaktuMasuk.value), // 🔥 create = waktu masuk (jika bukan mode edit, value = waktu masuk)

      karyawan_id: isEditMode.value
        ? props.Presensi?.karyawan_id // 🔥 lock agar tidak bisa diubah via devtools
        : formKaryawanID.value,
    },

    // 🔥 callback error dari parent
    onError: (errMsg) => {
       // 🔥 mapping error ke field yang sesuai
      const msg = errMsg.toLowerCase();

      if (msg.includes("karyawan")) {
        errors.value.karyawan_id = errMsg;
      } else if (msg.includes("masuk")) {
        errors.value.waktu_masuk = errMsg;
      } else if (msg.includes("pulang")) {
        errors.value.waktu_pulang = errMsg;
      } else if (msg.includes("tanggal")) {
        errors.value.tanggal = errMsg;
      }
    },

    // 🔥 callback success → tutup modal
    onSuccess: () => closeModal(),
  });
};

// === ERROR STATE (INLINE VALIDATION) ===
// 🔥 simpan error per field (biar bisa tampil di bawah input)
const errors = ref({
  tanggal: "",
  karyawan_id: "",
  waktu_masuk: "",
  waktu_pulang: "",
});
</script>

<template>
  <!-- 🔥 v-model ke props.modelValue → modal dikontrol parent -->
  <v-dialog v-model="props.modelValue" max-width="400">

    <v-card>

      <!-- TITLE -->
      <v-card-title class="text-h6">
        {{
          props.mode === "create"
            ? "Tambah Presensi Karyawan"
            : "Edit Presensi Karyawan"
        }}
      </v-card-title>

      <!-- ================================================= -->
      <!-- FORM (CREATE & EDIT) -->
      <!-- NOTE:
           - gunakan v-form agar Enter dan klik tombol submit
             mengikuti behaviour standar HTML
           - @submit.prevent mencegah reload halaman
      -->
      <!-- ================================================= -->
      <v-form @submit.prevent="save">

        <v-card-text>

          <!-- 🔥 DATE INPUT -->
          <!-- simple & clean (native date input) -->
          <v-text-field
            ref="inputTanggalRef"
            label="Tanggal"
            v-model="formTanggal"
            type="date"
            :disabled="isEditMode"
            class="mb-5"
            density="compact"
            variant="outlined"
            hide-details="auto"
            style="max-width: 180px"

            :error="!!errors.tanggal"
            :error-messages="errors.tanggal"
          />

          <!-- karyawan bisa banyak, jadi pakai autocomplete agar tidak makan space -->
          <v-autocomplete
            v-model="formKaryawanID"
            :items="karyawanOptions"
            :disabled="isEditMode"
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

          <!-- 🔥 WAKTU MASUK -->
          <v-menu
            v-model="menuMasuk"
            :close-on-content-click="false"
            transition="scale-transition"
            :disabled="isEditMode"
          >
            <template #activator="{ props }">
              <v-text-field
                v-model="formWaktuMasuk"
                label="Waktu Masuk"
                placeholder="HH:mm"
                autocomplete="off"
                v-bind="props"
                :disabled="isEditMode"
                density="compact"
                variant="outlined"
                prepend-inner-icon="mdi-clock-start"
              />
            </template>

            <v-time-picker
              v-model="formWaktuMasuk"
              format="24hr"
              @update:modelValue="menuMasuk = false"
            />
          </v-menu>

          <!-- 🔥 WAKTU PULANG -->
          <v-menu
            v-model="menuPulang"
            :close-on-content-click="false"
            transition="scale-transition"
            :disabled="isCreateMode"
          >
            <template #activator="{ props }">
              <v-text-field
                ref="inputWaktuPulangRef"
                v-model="formWaktuPulang"
                label="Waktu Pulang"
                placeholder="HH:mm"
                autocomplete="off"
                v-bind="props"
                :disabled="isCreateMode"
                density="compact"
                variant="outlined"
                prepend-inner-icon="mdi-clock-end"
              />
            </template>

            <v-time-picker
              v-model="formWaktuPulang"
              format="24hr"
              @update:modelValue="menuPulang = false"
            />

            <!-- 🔥 ERROR VALIDASI WAKTU -->
            <v-alert
              v-if="isInvalidTime"
              type="error"
              density="compact"
              class="mb-2"
            >
              Waktu pulang tidak boleh lebih awal dari waktu masuk
            </v-alert>
          </v-menu>

        </v-card-text>

        <!-- ACTION -->
        <v-card-actions>

          <v-spacer />

          <!-- BATAL -->
          <v-btn
            color="red"
            variant="tonal"
            size="small"
            type="button"
            @click="closeModal"
          >
            Batal
          </v-btn>

          <!-- SAVE -->
          <!-- NOTE:
               - type="submit" mengikuti mekanisme submit bawaan form
               - tidak lagi memakai @click agar tidak terjadi double submit
               - disable tombol untuk prevent invalid submit
          -->
          <v-btn
            type="submit"
            color="primary"
            variant="elevated"
            size="small"
            :disabled="
              !formTanggal ||
              !formKaryawanID ||
              !formWaktuMasuk ||
              (props.mode === 'edit' && !formWaktuPulang) ||
              isInvalidTime
            "
          >
            Simpan
          </v-btn>

        </v-card-actions>

      </v-form>

      <!-- TIDAK ADA FITUR DELETE PRESENSI -->

    </v-card>

  </v-dialog>
</template>