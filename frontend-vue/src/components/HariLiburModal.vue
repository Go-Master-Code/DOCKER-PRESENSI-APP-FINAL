<script setup>
import { ref, watch, nextTick } from "vue";

// 🔥 Ref untuk akses langsung ke komponen input (dipakai untuk autofocus)
const inputTanggalRef = ref(null);
const inputKeteranganRef = ref(null);

// 🔥 Props dari parent (HariLiburView)
// modelValue → untuk buka/tutup modal (v-model)
// hariLibur → data yang dipilih (null saat create)
// mode → menentukan behavior (create | edit | delete)
const props = defineProps({
    modelValue: Boolean,
    hariLibur: { type: Object, default: () => null },
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

// 🔥 Utility: hanya angka (sanitize input)
// penting karena user bisa paste string
// const onlyNumber = (value) => {
//   return value.replace(/\D/g, ""); // hapus semua selain angka
// };

// === Sync props ke form saat modal dibuka ===
// 🔥 Kenapa watch modelValue?
// karena props.hariLibur tidak selalu berubah (misal create = null terus)
// jadi trigger terbaik adalah saat modal dibuka
watch(
  () => props.modelValue,
  async (isOpen) => {
    if (!isOpen) return;

    // 🔥 RESET ERROR SETIAP MODAL DIBUKA
    errors.value = {
        tanggal: "",
        keterangan: "",
    };

    if (props.mode === "create") {
      // 🔥 Reset form saat create
      formTanggal.value = "";
      formKeterangan.value = "";

      // 🔥 Autofocus ke Tanggal (UX lebih cepat)
      await nextTick(); // tunggu DOM render
      inputTanggalRef.value?.focus();
    } else {
      // 🔥 Isi form saat edit
      formTanggal.value = props.hariLibur?.tanggal || "";
      formKeterangan.value = props.hariLibur?.keterangan || "";

      // 🔥 Autofocus ke Nama (karena ID tidak boleh diubah)
      await nextTick();
      inputTanggalRef.value?.focus();
    }
  }
);

// 🔥 UX penting:
// Error hilang otomatis saat user mulai memperbaiki input

watch(formTanggal, () => {
    errors.value.tanggal = "";
});

watch(formKeterangan, () => {
    errors.value.keterangan = "";
});

// === Close modal ===
// 🔥 gunakan emit agar parent yang kontrol state
const closeModal = () => emit("update:modelValue", false);

// === Save handler (create & edit) ===
const save = async () => {
  // reset error
  errors.value = { tanggal: "", keterangan: "" };

  emit("save", {
    payload: {
      ...props.hariLibur,
      tanggal: formTanggal.value,
      keterangan: formKeterangan.value,
    }, // 🔥 CALLBACK ERROR (dari parent)
    // Modal tidak tahu API → hanya terima hasil
    onError: (errMsg) => {
      // 🔥 tangkap error dari parent DI SINI
      if (errMsg.toLowerCase().includes("tanggal")) { //👉 cek apakah error mengandung kata "tanggal"
        errors.value.tanggal = errMsg; //👉 kirim error ke field tanggal
      } else if (errMsg.toLowerCase().includes("keterangan")) { //👉 kalau error terkait keterangan → masuk ke field keterangan
        errors.value.keterangan = errMsg;
      } else { //👉 kalau tidak jelas: tetap tampilkan error (biar user tidak blank), default ke field tanggal
        errors.value.tanggal = errMsg;
      }
    }, // 🔥 CALLBACK SUCCESS
    onSuccess: () => {
      closeModal();
    },
  });
};

// === Delete handler ===
const remove = () => {
  emit("delete", props.hariLibur);
  closeModal();
};

// === ERROR STATE (INLINE VALIDATION) ===
// 🔥 simpan error per field (biar bisa tampil di bawah input)
const errors = ref({
  tanggal: "",
  keterangan: "",
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
            ? "Tambah Hari Libur"
            : props.mode === "edit"
            ? "Edit Hari Libur"
            : "Hapus Hari Libur"
        }}
      </v-card-title>

      <!-- ================================================= -->
      <!-- FORM -->
      <!-- ================================================= -->
      <v-form
        v-if="props.mode !== 'delete'"
        @submit.prevent="save"
      >

        <!-- CONTENT -->
        <v-card-text>

          <!-- FORM (CREATE & EDIT) -->
          <div v-if="props.mode !== 'delete'">

            <!-- 🔥 ID -->
            <!-- NOTE:
                - disabled saat edit (ID immutable)
                - maxlength + counter → UX guidance
                - inputmode → keyboard angka di mobile
            -->

            <!-- 🔥 DATE INPUT -->
              <!-- simple & clean (native date input) -->
              <v-text-field
                  label="Tanggal"
                  v-model="formTanggal"
                  type="date"
                  class="mb-5"
                  density="compact"
                  variant="outlined"
                  hide-details="auto"
                  style="max-width: 180px"

                  :error="!!errors.tanggal"
                  :error-messages="errors.tanggal"
              />

              <!-- 🔥 Keterangan -->
              <!-- NOTE:
                  - autofocus saat edit
                  - bisa ditambahkan rules kalau mau wajib -->
              <v-text-field
                  label="Keterangan"
                  ref="inputKeteranganRef"
                  v-model="formKeterangan"
                  density="compact"
                  variant="outlined"
                  autocomplete="off"
                  maxlength="100"
                  counter

                  :error="!!errors.keterangan"
                  :error-messages="errors.keterangan"
              />
          </div>
        </v-card-text>

        <!-- ACTIONS -->
        <v-card-actions>
          <v-spacer />

          <!-- BATAL -->
          <v-btn
            color="red"
            variant="tonal"
            size="small"
            @click="closeModal"
          >
            Batal
          </v-btn>

          <!-- SAVE -->
          <!-- 🔥 disable tombol untuk prevent invalid submit -->
          <v-btn
            v-if="props.mode !== 'delete'"
            color="primary"
            variant="elevated"
            size="small"
            :disabled="!formTanggal || formKeterangan.length > 100 || !formKeterangan"
            @click="save"
          >
            Simpan
          </v-btn>
        </v-card-actions>
      </v-form>

      <!-- ================================================= -->
      <!-- DELETE MODE -->
      <!-- ================================================= -->
      <!-- 🔥 gunakan optional chaining untuk safety -->
      <template v-else>
        <v-card-text>
          Apakah yakin ingin menghapus hari libur
          <strong>{{ props.hariLibur?.keterangan }}</strong>?
        </v-card-text>

        <v-card-actions>
          <v-spacer />
          <!-- DELETE -->
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