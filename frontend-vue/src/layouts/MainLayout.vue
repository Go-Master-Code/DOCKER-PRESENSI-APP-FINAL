<script setup>
    import { ref } from "vue";
    import { useDisplay } from "vuetify";
    import { watch } from "vue";
    import { useSnackbar } from "@/composables/useSnackbar";
    // 🔥 router
    import { useRouter } from "vue-router";

    // import backupService untuk backup DB by single click
    import * as backupService from "@/services/backupService";

    // 🔥 auth composable
    import { useAuth } from "@/composables/useAuth";

    // toggle dark / light mode
    import { useAppTheme } from "@/composables/useTheme";

    // ==============================
    // 🔥 HANDLE BACKUP DATABASE
    // ==============================
    // Fungsi untuk:
    // 1. Request backup ke backend
    // 2. Menerima file .sql
    // 3. Download otomatis ke browser
    const handleBackup = async () => {
        try {
            // 🔥 call API backup
            const res =
                await backupService.backupDatabase();

            // 🔥 ubah response file menjadi URL sementara browser
            const url =
                window.URL.createObjectURL(
                    new Blob([res.data])
                );

            // 🔥 buat element <a> sementara
            const link =
                document.createElement("a");

            // 🔥 arahkan link ke file blob
            link.href = url;

            // 🔥 nama file download
            // 🔥 ambil tanggal & waktu sekarang
            const now = new Date();

            // format:
            // 2026-05-11_14-30-25
            const formattedDate =
                `${now.getFullYear()}-` +
                `${String(now.getMonth() + 1).padStart(2, "0")}-` +
                `${String(now.getDate()).padStart(2, "0")}_` +
                `${String(now.getHours()).padStart(2, "0")}-` +
                `${String(now.getMinutes()).padStart(2, "0")}-` +
                `${String(now.getSeconds()).padStart(2, "0")}`;

            // 🔥 set nama file dinamis
            link.setAttribute(
                "download",
                `backup_DB_${formattedDate}.sql`
            );

            // 🔥 masukkan ke DOM
            document.body.appendChild(link);

            // 🔥 trigger download otomatis
            link.click();

            // 🔥 cleanup memory
            window.URL.revokeObjectURL(url);

        } catch (err) {
            // 🔥 tampilkan error
            console.error(err);
        }
    };
    // =========================
    // 🔥 DARK MODE
    // =========================
    const {
        isDark,
        toggleTheme
    } = useAppTheme();

    // 🔥 router
    const router = useRouter();

    // 🔥 auth logout + setToken + setRole untuk RBA
    const { logout, isAdmin } = useAuth();
    
    // =====================
    // 🔥 LOGOUT
    // =====================
    const handleLogout = () => {

        // hapus token
        logout();

        // redirect ke login
        router.push("/login");
    };

    // 🔥 ambil state global snackbar
    const { show, text, color, timeout } = useSnackbar();

    // 🔥 ambil info device (reactive dari Vuetify)
    const { mobile } = useDisplay();

    // ⚠️ ini hanya dieksekusi sekali saat mount, jadi pakai watch di atas saja
    const drawer = ref(!mobile.value);

    // 🔥 drawer state
    // NOTE:
    // - desktop → true (terbuka)
    // - mobile → false (tertutup)
    watch(mobile, (isMobile) => {
        drawer.value = !isMobile;
    });
</script>

<template>
  <!-- 🔥 ROOT VUETIFY (WAJIB) -->
  <v-app>

    <!-- ===================== -->
    <!-- 🔥 APP BAR -->
    <!-- ===================== -->
    <v-app-bar elevation="1" border="b">

        <!-- 🔥 tombol toggle drawer -->
        <v-app-bar-nav-icon @click="drawer = !drawer" />

        <!-- 🔥 title -->
        <!-- NOTE: gunakan typography di sini (bukan di child div) -->
        <v-toolbar-title class="font-weight-bold">
                Dashboard
        </v-toolbar-title>

        <!-- ========================= -->
        <!-- 🔥 DARK MODE TOGGLE -->
        <!-- ========================= -->
        <v-btn
            icon
            variant="text"
            @click="toggleTheme"
        >
            <v-icon>
                {{
                    isDark()
                        ? "mdi-weather-sunny"
                        : "mdi-weather-night"
                }}
            </v-icon>
        </v-btn>
    </v-app-bar>

    <!-- ===================== -->
    <!-- 🔥 NAVIGATION DRAWER -->
    <!-- ===================== -->
    <v-navigation-drawer
        v-model="drawer"
        :permanent="!mobile"  
        :temporary="mobile" 
        width="260"
        elevation="2"
    >

    <!-- :rail="!mobile" adalah mini side bar yang hanya akan tampil di desktop, bentuknya icon kecil -->

      <!-- 🔥 MENU LIST -->
      <v-list nav>

        <!-- MENU: KARYAWAN -->
        <v-list-item
            prepend-icon="mdi-account-group"
            title="Karyawan"
            to="/karyawan"
            active-class="bg-primary"
            @click="mobile && (drawer = false)"
        />

        <!-- MENU: HARI LIBUR -->
        <v-list-item
            prepend-icon="mdi-calendar-remove"
            title="Hari Libur"
            to="/hariLibur"
            active-class="bg-primary"
            @click="mobile && (drawer = false)"
        />

        <!-- MENU: PRESENSI -->
        <v-list-item
            prepend-icon="mdi-calendar-check"
            title="Presensi"
            to="/presensi"
            active-class="bg-primary"
            @click="mobile && (drawer = false)"
        />

        <!-- MENU: IJIN KARYAWAN -->
        <v-list-item
            prepend-icon="mdi-account-clock"
            title="Ijin Karyawan"
            to="/ijin"
            active-class="bg-primary"
            @click="mobile && (drawer = false)"
        />

        <!-- MENU: USER MANAGEMENT -->
        <v-list-item
            v-if="isAdmin()"
            prepend-icon="mdi-account-cog"
            title="Manajemen User"
            to="/user"
            active-class="bg-primary"
            @click="mobile && (drawer = false)"
        />

        <!-- ========================= -->
        <!-- 🔥 DROPDOWN LAPORAN -->
        <!-- ========================= -->
        <v-list-group value="laporan">

            <!-- HEADER -->
            <template #activator="{ props }">
            <v-list-item
                v-bind="props"
                title="Laporan Presensi"
                prepend-icon="mdi-file-chart"
            />
            </template>

            <!-- CHILD MENU -->
            <v-list-item
                title="Per Periode"
                to="/laporan/periode"
                prepend-icon="mdi-calendar-range"
                active-class="bg-primary"
            />

            <v-list-item
                title="Per Karyawan"
                to="/laporan/karyawan"
                prepend-icon="mdi-account"
                active-class="bg-primary"
            />
        </v-list-group>

        <!-- ========================= -->
        <!-- 🔥 BACKUP DB SINGLE CLICK -->
        <!-- ========================= -->
        <!-- <v-list-item
            v-if="isAdmin()"
            prepend-icon="mdi-database-export"
            title="Backup DB"
            base-color="green"
            @click="handleBackup"
        /> -->

        <!-- ========================= -->
        <!-- 🔥 LOGOUT -->
        <!-- ========================= -->
        <v-divider class="my-2" />

        <v-list-item
            prepend-icon="mdi-logout"
            title="Logout"
            base-color="red"
            @click="handleLogout"
        />
      </v-list>
    </v-navigation-drawer>

    <!-- ===================== -->
    <!-- 🔥 MAIN CONTENT -->
    <!-- ===================== -->
    <v-main>

      <!-- 🔥 halaman akan dirender di sini -->
      <router-view />

    </v-main>

    <!-- GLOBAL SNACKBAR -->
    <v-snackbar
        v-model="show"
        :color="color"
        :timeout="timeout"
        location="top right"
    >
    
        <div class="d-flex align-center gap-2">
            <v-icon>
            {{
                color === "success" ? "mdi-check-circle" :
                color === "error" ? "mdi-alert-circle" :
                "mdi-information"
            }}
            </v-icon>

            {{ text }}
        </div>

        <!-- tombol close -->
        <template #actions>
            <v-btn variant="text" @click="show = false">
                Tutup
            </v-btn>
        </template>
    </v-snackbar>
  </v-app>
</template>