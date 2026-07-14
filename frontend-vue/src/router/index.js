import { createRouter, createWebHistory } from "vue-router";
import MainLayout from "@/layouts/MainLayout.vue";
import KioskLayout from "@/layouts/KioskLayout.vue";

// views
import PresensiView from "@/views/PresensiView.vue";
import KaryawanView from "@/views/KaryawanView.vue";
import HariLiburView from "@/views/HariLiburView.vue";
import IjinView from "@/views/IjinView.vue";
import UserView from "@/views/UserView.vue";
import PresensiScanView from "@/views/PresensiScanView.vue";
import LoginView from "@/views/LoginView.vue";

// view laporan
import LaporanPresensiPeriodeView from "@/views/LaporanPresensiPeriodeView.vue"
import LaporanPresensiKaryawanView from "@/views/LaporanPresensiKaryawanView.vue"


const routes = [
  {
    path: "/",

    // 🔥 MainLayout jadi parent (wrapper semua halaman)
    component: MainLayout,

    // 🔥 children → akan dirender di <router-view /> dalam MainLayout
    children: [
      {
        path: "", 
        name: "presensi",

        // 🔥 default route ("/")
        // saat buka localhost:5173 → langsung ke Presensi
        component: PresensiView,
        meta: {
          requiresAuth: true,
        },
      },

      {
        path: "karyawan",
        name: "karyawan",

        // 🔥 hasil URL → /karyawan
        component: KaryawanView,
      },

      {
        path: "hariLibur",
        name: "hariLibur",

        // 🔥 hasil URL → /karyawan
        component: HariLiburView,
      },

      {
        path: "ijin",
        name: "ijin",

        // 🔥 hasil URL → /ijin
        component: IjinView,
      },

      {
        path: "presensi",
        name: "presensi-page",

        // 🔥 optional (kalau mau /presensi juga bisa diakses)
        component: PresensiView,
      },

      {
        path: "user",
        name: "user",

        // 🔥 hasil URL → /user
        component: UserView,
      },

      // ==============================
      // 🔥 LAPORAN (MASUK SINI!)
      // ==============================
      {
        path: "laporan/periode",
        name: "LaporanPeriode",
        component: LaporanPresensiPeriodeView,
      },
      {
        path: "laporan/karyawan",
        name: "LaporanKaryawan",
        component: LaporanPresensiKaryawanView,
      },
    ],
  },

  // halaman login (TANPA menu)
  {
    path: "/login",

    // 🔥 hasil URL → /login
    component: LoginView,
  },

  // 🔥 layout kiosk (TANPA menu)
  {
    path: "/scan",
    component: KioskLayout,
    children: [
      {
        path: "",
        name: "PresensiScan",
        component: PresensiScanView,
      },
    ],
  },

];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 🔥 ROUTE GUARD
router.beforeEach((to, from) => {

  // 🔥 ambil token login
  const token = localStorage.getItem("token");

  // 🔥 jika route butuh auth tapi belum login
  if (to.meta.requiresAuth && !token) {
    // redirect ke login
    return "/login";
  }

  // 🔥 jika sudah login tapi buka halaman login
  if (to.path === "/login" && token) {
    return "/";
  }

  // 🔥 izinkan lanjut
  return true;
});

export default router;