<script setup>
import { ref, watch, nextTick, computed } from "vue";
import * as roleService from "@/services/roleService";

// 🔥 Ref untuk akses langsung ke komponen input (dipakai untuk autofocus)
const inputUsernameRef = ref(null);

// tombol show/hide 👁️ untuk password
const showPassword = ref(false);

// 🔥 LIST DATA UNTUK DROPDOWN
const roleOptions = ref([]); // 🔥 state untuk menampung data dropdown jenis role

// 🔥 LOADING STATE (UX penting)
// 🔥 loading indicator (UX penting saat fetch API)
const loadingRole = ref(false);

// 🔥 Props dari parent (UserView)
// modelValue → untuk buka/tutup modal (v-model)
// user → data yang dipilih (null saat create)
// mode → menentukan behavior (create | edit | delete)
const props = defineProps({
  modelValue: Boolean,
  user: { type: Object, default: () => null },
  mode: { type: String, default: "edit" },
});

// 🔥 Emit event ke parent
// update:modelValue → untuk close modal
// save → create/update
// delete → hapus data
const emit = defineEmits(["update:modelValue", "save", "delete"]);

// === Reactive form (state lokal, tidak langsung mengubah props) ===
const formUsername = ref("");
const formPassword = ref("");
const formConfirmPassword = ref("");
const formRole = ref(null);

// memastikan password yang diinput dan confirm password sama
const isPasswordMismatch = computed(() => {
  if (!formPassword.value || !formConfirmPassword.value) return false; // jika formPassword ATAU confirm password salah satunya masih kosong (belum diinput), jangan anggap error
  return formPassword.value !== formConfirmPassword.value; // kalau tidak sama -> true, kalau sama -> false
});

// watcher agar error hilang saat user mengetik kembali confirm password
watch([formPassword, formConfirmPassword], () => {
  if (formPassword.value === formConfirmPassword.value) {
    errors.value.password = ""; // hilangkan error saat password dan confirm sudah sama
  }
});

// === Sync props ke form saat modal dibuka ===
// 🔥 Kenapa watch modelValue?
// karena props.user tidak selalu berubah (misal create = null terus)
// jadi trigger terbaik adalah saat modal dibuka
watch(
  () => props.modelValue,
  async (isOpen) => {
    if (!isOpen) return;

    // get master data role masukkan ke dalam combobox
    await fetchDropdownData();

    // 🔥 RESET ERROR SETIAP MODAL DIBUKA
    errors.value = {
      username: "",
      password: "",
      role_id: "",
    };

    if (props.mode === "create") {
      // 🔥 Reset form saat create
      formUsername.value = "";
      formPassword.value = "";
      formConfirmPassword.value = "";
      formRole.value = null;

      // 🔥 Autofocus ke ID (UX lebih cepat)
      await nextTick(); // tunggu DOM render
      inputUsernameRef.value?.focus();
    } else {
      // 🔥 Isi form saat edit
      formUsername.value = props.user?.username || "";
      formPassword.value = props.user?.password || "";
      formConfirmPassword.value = props.user?.password || "";
      formRole.value = props.user?.role_id || null;

      // 🔥 Autofocus ke Nama (karena ID tidak boleh diubah)
      await nextTick();
      inputUsernameRef.value?.focus();
    }
  }
);

// load dropdown data jenis role user
const fetchDropdownData = async () => {
  try {
    // 🔥 aktifkan loading
    loadingRole.value = true;

    // 🔥 ambil 1 API role
    const [resRole] = await Promise.all([
      roleService.getAllRole(),
    ]);

    // 🔥 mapping data ke format Vuetify (WAJIB)
    // Vuetify butuh: { title, value }
    roleOptions.value = resRole.data.data.map(r => ({
      title: r.nama,   // 🔥 teks yang ditampilkan ke user
      value: r.id,     // 🔥 nilai yang dikirim ke backend
    }));

  } catch (err) {
    console.error(err); // 🔥 debug kalau API error
  } finally {
    // 🔥 matikan loading
    loadingRole.value = false;
  }
};

// 🔥 UX penting:
// Error hilang otomatis saat user mulai memperbaiki input

watch(formUsername, () => {
    errors.value.username = "";
});

watch(formPassword, () => {
    errors.value.password = "";
});


// === Close modal ===
// 🔥 gunakan emit agar parent yang kontrol state
const closeModal = () => emit("update:modelValue", false);

// === Enter untuk submit form ===
// 🔥 UX improvement → user tidak perlu klik tombol
// const handleEnter = () => {
//   if (props.mode === "delete") return;

//   // 🔥 Validasi sederhana (frontend guard)
//   if (!formUsername.value || formUsername.value.length > 50 || !formPassword.value || formPassword.value.length > 50) return;

//   save();
// };

// === Save handler (create & edit) ===
const save = async () => {
  // reset error
  errors.value = { username: "", password: "", role_id: "" };

  // double safety jika password dan confirm password tidak sama
  if (formPassword.value !== formConfirmPassword.value) {
    errors.value.password = "Password tidak sama";
    return;
  }

  const payload = {
    ...(props.mode === "edit" ? props.user : {}),
    username: formUsername.value,
    password: formPassword.value,
    role_id: formRole.value,
  };

  // 🔥 DEBUG DI SINI (PASTI MUNCUL)
  // console.log("=== PAYLOAD DARI MODAL ===");
  // console.log(payload);

  emit("save", {
    payload, // 🔥 CALLBACK ERROR (dari parent)

    // Modal tidak tahu API → hanya terima hasil
    onError: (errMsg) => {
      // 🔥 tangkap error dari parent DI SINI
      if (errMsg.toLowerCase().includes("user")) { // error mengandung kata user (username)
        errors.value.username = errMsg;
      } else if (errMsg.toLowerCase().includes("password")) { // error mengandung kata password
        errors.value.password = errMsg;
      } else { // error lain
        errors.value.role_id = errMsg;
      }
    }, // 🔥 CALLBACK SUCCESS
    onSuccess: () => {
      closeModal();
    },
  });
};

// === Delete handler ===
const remove = () => {
  emit("delete", props.user);
  closeModal();
};

// === ERROR STATE (INLINE VALIDATION) ===
// 🔥 simpan error per field (biar bisa tampil di bawah input)
const errors = ref({
  username: "",
  password: "",
  role_id: "",
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
            ? "Tambah User"
            : props.mode === "edit"
            ? "Edit User"
            : "Hapus User"
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

          <!-- 🔥 Username -->
          <!-- NOTE:
               - maxlength + counter → UX guidance
          -->
          <v-text-field
            ref="inputUsernameRef"
            label="Username (Max. 50 karakter)"
            v-model="formUsername"
            density="compact"
            variant="outlined"
            maxlength="50"
            autocomplete="off"
            counter

            :error="!!errors.username"
            :error-messages="errors.username"
          />

          <!-- 🔥 Password -->
          <v-text-field
            ref="inputPasswordRef"
            label="Password"
            v-model="formPassword"
            :type="showPassword ? 'text' : 'password'"
            density="compact"
            variant="outlined"
            autocomplete="off"

            :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
            @click:append-inner="showPassword = !showPassword"

            :error="!!errors.password"
            :error-messages="errors.password"
          />

          <!-- 🔥 Confirm Password -->
          <v-text-field
            label="Confirm Password"
            v-model="formConfirmPassword"
            :type="showPassword ? 'text' : 'password'"
            density="compact"
            variant="outlined"
            autocomplete="off"

            :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
            @click:append-inner="showPassword = !showPassword"

            :error="isPasswordMismatch"
            :error-messages="isPasswordMismatch ? 'Password tidak sama' : ''"
          />

          <!-- jenis role pakai combobox karena relatif sedikit datanya -->
          <v-select
            v-model="formRole"
            :items="roleOptions"
            item-title="title"
            item-value="value"
            label="Role"
            density="compact"
            variant="outlined"
            :loading="loadingRole"
            clearable

            :error="!!errors.role_id"
            :error-messages="errors.role_id"
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
              !formUsername ||
              formUsername.length > 50 ||
              !formRole ||
              !formPassword ||
              formPassword.length > 50 ||
              !formConfirmPassword ||
              formConfirmPassword.length > 50 ||
              isPasswordMismatch
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
          Apakah yakin ingin menghapus user
          <strong>{{ props.user?.username }}</strong>?
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