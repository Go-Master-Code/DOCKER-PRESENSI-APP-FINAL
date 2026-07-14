<script setup>
  import { ref, onMounted } from "vue";
  import { useRouter } from "vue-router";

  import * as authService from "@/services/authService";
  import { useAuth } from "@/composables/useAuth";

  // reset session setiap kali login ulang
  import { resetSessionRedirect } from "@/services/api";

  // import logo
  import logo from "@/assets/logo.png";

  const router = useRouter();

  const { setToken, setRole } = useAuth();

  const username = ref("");
  const password = ref("");

  // 🔥 ref component username untuk autofocus
  const usernameField = ref(null);

  const loading = ref(false);
  const error = ref("");

  // === Enter untuk submit form ===
// 🔥 UX improvement → user tidak perlu klik tombol
  const handleEnter = () => {
    // 🔥 Validasi sederhana (frontend guard)
    if (!username.value || !password.value) return;

    handleLogin();
  };

  // 🔥 autofocus saat halaman dibuka
  onMounted(() => {
    usernameField.value?.focus();
  });
  
  const handleLogin = async () => {

    error.value = "";

    try {

      loading.value = true;

      const res = await authService.login({
        username: username.value,
        password: password.value,
      });

      // 🔥 DEBUG
      // console.log("LOGIN RESPONSE:", res.data);

      // ambil token dari backend
      const token = res.data.token; // hati-hati ambil token, perhatikan response json di struct backend, dia berada di luar slice data

      // setelah login sukses, simpan token dan role
      setToken(token);

      // 🔥 reset interceptor state (session login)
      resetSessionRedirect();

      setRole(res.data.data.role_nama) // ambil sesuai response api Login (cek postman)

      // redirect
      router.push("/");

    } catch (err) {

      console.error(err);

      error.value =
        err.response?.data?.error ||
        "Login gagal, silakan tunggu beberapa saat lagi!";

    } finally {
      loading.value = false;
    }
  };
</script>

<template>
  <v-container
    fluid
    class="fill-height bg-grey-lighten-4"
  >
    <v-row
      justify="center"
      align="center"
    >

      <!-- LOGIN CARD -->
      <v-col
        cols="12"
        sm="10"
        md="6"
        lg="4"
      >

        <v-card @keyup.enter="handleEnter"
          class="pa-8 rounded-xl"
          elevation="8"
        >

          <!-- HEADER -->
          <div class="text-center mb-6">
            <v-img
              :src="logo"
              width="110"
              class="mx-auto mb-1"
            />

            <div style="font-size: clamp(28px, 1.5vw, 20px);font-weight: 600;">
              Login
            </div>

            <div style="font-size: clamp(20px, 1.5vw, 20px); color: #666;">
              Sistem Presensi Karyawan
            </div>

          </div>

          <!-- USERNAME -->
          <v-text-field
            v-model="username"
            ref="usernameField"
            autocomplete="off"
            label="Username"
            prepend-inner-icon="mdi-account"
            variant="outlined"
            density="comfortable"
            class="mb-2"
          />

          <!-- PASSWORD -->
          <v-text-field
            v-model="password"
            autocomplete="off"
            label="Password"
            type="password"
            prepend-inner-icon="mdi-lock"
            variant="outlined"
            density="comfortable"
            class="mb-2"
          />

          <!-- ERROR -->
          <v-alert
            v-if="error"
            type="error"
            variant="tonal"
            density="compact"
            class="mb-4"
          >
            {{ error }}
          </v-alert>

          <!-- BUTTON -->
          <v-btn
            block
            color="primary"
            size="large"
            :disabled="!username || !password"
            :loading="loading"
            @click="handleLogin"
          >
            <v-icon start>
              mdi-login
            </v-icon>
            Login
          </v-btn>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>