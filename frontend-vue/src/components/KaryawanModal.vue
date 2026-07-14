<script setup>
import { ref, watch, nextTick } from "vue";

// 🔥 Ref untuk akses langsung ke komponen input (dipakai untuk autofocus)
const inputIDRef = ref(null);
const inputNamaRef = ref(null);

// 🔥 Props dari parent (KaryawanView)
// modelValue → untuk buka/tutup modal (v-model)
// karyawan → data yang dipilih (null saat create)
// mode → menentukan behavior (create | edit | delete)
const props = defineProps({
  modelValue: Boolean,
  karyawan: { type: Object, default: () => null },
  mode: { type: String, default: "edit" },
});

// 🔥 Emit event ke parent
// update:modelValue → untuk close modal
// save → create/update
// delete → hapus data
const emit = defineEmits(["update:modelValue", "save", "delete"]);

// === Reactive form (state lokal, tidak langsung mengubah props) ===
const formID = ref("");
const formNama = ref("");
const formStatus = ref(false);

// 🔥 Utility: hanya angka (sanitize input)
// penting karena user bisa paste string
const onlyNumber = (value) => {
  return value.replace(/\D/g, ""); // hapus semua selain angka
};

// 🔥 Watch formID → auto filter setiap perubahan
// NOTE:
// - lebih aman dari keypress karena handle paste juga
// - slice(0,10) memastikan max 10 digit
watch(formID, (val) => {
  formID.value = onlyNumber(val).slice(0, 10);
});

// === Sync props ke form saat modal dibuka ===
// 🔥 Kenapa watch modelValue?
// karena props.karyawan tidak selalu berubah (misal create = null terus)
// jadi trigger terbaik adalah saat modal dibuka
watch(
  () => props.modelValue,
  async (isOpen) => {
    if (!isOpen) return;

    // 🔥 RESET ERROR SETIAP MODAL DIBUKA
    errors.value = {
      id: "",
      nama: "",
    };

    if (props.mode === "create") {
      // 🔥 Reset form saat create
      formID.value = "";
      formNama.value = "";
      formStatus.value = true;

      // 🔥 Autofocus ke ID (UX lebih cepat)
      await nextTick(); // tunggu DOM render
      inputIDRef.value?.focus();
    } else {
      // 🔥 Isi form saat edit
      formID.value = props.karyawan?.id || "";
      formNama.value = props.karyawan?.nama || "";
      formStatus.value = props.karyawan?.aktif || false;

      // 🔥 Autofocus ke Nama (karena ID tidak boleh diubah)
      await nextTick();
      inputNamaRef.value?.focus();
    }
  }
);

// 🔥 UX penting:
// Error hilang otomatis saat user mulai memperbaiki input

watch(formID, () => {
    errors.value.id = "";
});

watch(formNama, () => {
    errors.value.nama = "";
});


// === Close modal ===
// 🔥 gunakan emit agar parent yang kontrol state
const closeModal = () => emit("update:modelValue", false);

// === Enter untuk submit form ===
// 🔥 UX improvement → user tidak perlu klik tombol
// const handleEnter = () => {
//   if (props.mode === "delete") return;

//   // 🔥 Validasi sederhana (frontend guard)
//   if (!formID.value || formID.value.length !== 10 || !formNama.value) return;

//   save();
// };

// === Save handler (create & edit) ===
const save = async () => {
  // reset error
  errors.value = { id: "", nama: "" };

  emit("save", {
    payload: {
      ...props.karyawan,
      id: formID.value,
      nama: formNama.value,
      aktif: formStatus.value,
    }, // 🔥 CALLBACK ERROR (dari parent)
    // Modal tidak tahu API → hanya terima hasil
    onError: (errMsg) => {
      // 🔥 tangkap error dari parent DI SINI
      if (errMsg.toLowerCase().includes("id")) {
        errors.value.id = errMsg;
      } else if (errMsg.toLowerCase().includes("nama")) {
        errors.value.nama = errMsg;
      } else {
        errors.value.nama = errMsg;
      }
    }, // 🔥 CALLBACK SUCCESS
    onSuccess: () => {
      closeModal();
    },
  });
};

// === Delete handler ===
const remove = () => {
  emit("delete", props.karyawan);
  closeModal();
};

// === ERROR STATE (INLINE VALIDATION) ===
// 🔥 simpan error per field (biar bisa tampil di bawah input)
const errors = ref({
  id: "",
  nama: "",
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
            ? "Tambah Karyawan"
            : props.mode === "edit"
            ? "Edit Karyawan"
            : "Hapus Karyawan"
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
      <v-form
        v-if="props.mode !== 'delete'"
        @submit.prevent="save"
      >

        <v-card-text>

          <!-- 🔥 ID -->
          <!-- NOTE:
               - disabled saat edit (ID immutable)
               - maxlength + counter → UX guidance
               - inputmode → keyboard angka di mobile
          -->
          <v-text-field
            ref="inputIDRef"
            label="ID (10 digit)"
            v-model="formID"
            density="compact"
            variant="outlined"
            maxlength="10"
            autocomplete="off"
            counter
            inputmode="numeric"
            :disabled="props.mode === 'edit'"

            :error="!!errors.id"
            :error-messages="errors.id"

            :rules="[
              v => !!v || 'ID wajib diisi',
              v => v.length === 10 || 'ID harus 10 digit'
            ]"
          />

          <!-- 🔥 Nama -->
          <!-- NOTE:
               - autofocus saat edit
               - bisa ditambahkan rules kalau mau wajib
          -->
          <v-text-field
            ref="inputNamaRef"
            label="Nama"
            v-model="formNama"
            density="compact"
            variant="outlined"
            autocomplete="off"

            :error="!!errors.nama"
            :error-messages="errors.nama"
          />

          <!-- 🔥 Status -->
          <!-- NOTE:
               - warna dinamis biar lebih visual
          -->
          <v-switch
            v-model="formStatus"
            label="Aktif"
            :color="formStatus ? 'green' : 'red'"
            density="compact"
            inset
          />

        </v-card-text>

        <!-- ACTION -->
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
              !formID ||
              formID.length !== 10 ||
              !formNama
            "
          >
            Simpan
          </v-btn>

        </v-card-actions>

      </v-form>

      <!-- ================================================= -->
      <!-- DELETE CONFIRM -->
      <!-- 🔥 gunakan optional chaining untuk safety -->
      <!-- ================================================= -->
      <template v-else>

        <v-card-text>
          Apakah yakin ingin menghapus karyawan
          <strong>{{ props.karyawan?.nama }}</strong>?
        </v-card-text>

        <v-card-actions>

          <v-spacer />

          <!-- BATAL -->
          <v-btn
            color="secondary"
            variant="tonal"
            size="small"
            @click="closeModal"
          >
            Batal
          </v-btn>

          <!-- DELETE -->
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