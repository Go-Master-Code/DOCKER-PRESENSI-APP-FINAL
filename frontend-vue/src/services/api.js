import axios from "axios";

// 🔥 import router untuk redirect
import router from "@/router";

// import untuk session expired dari composables
import { useSessionExpired } from "@/composables/useSessionExpired";

// ==============================
// 🔥 CREATE AXIOS INSTANCE
// ==============================
const api = axios.create({

    // base URL backend
    // baseURL: "http://localhost:8080",

    // versi implementasi docker
    baseURL: import.meta.env.VITE_API_URL,

    // timeout request
    timeout: 10000,
});

// ==============================
// 🔥 AUTO ATTACH JWT TOKEN
// ==============================
api.interceptors.request.use((config) => {
  // ambil token dari localStorage
  const token = localStorage.getItem("token");

  // jika ada token
  if (token) {
    // attach Authorization header
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

// ==============================
// 🔥 PREVENT MULTIPLE REDIRECT
// ==============================
export let isRedirecting = false; // harus direset setelah login ulang

// func reset session setiap user selesai login ulang
export const resetSessionRedirect = () => {
  isRedirecting = false; // reset kembali jadi false
}

// ==============================
// 🔥 GLOBAL RESPONSE HANDLER
// ==============================
api.interceptors.response.use(
    // response sukses
    (response) => response,
    // response error
    (error) => {
        // cek unauthorized (status 401)
        if (
            error.response?.status === 401 &&
            !isRedirecting
        ) {
            isRedirecting = true;

            // 🔥 hapus session login
            localStorage.removeItem("token");
            localStorage.removeItem("role");

            // 🔥 tampilkan modal session expired
            const {
                openSessionExpired
            } = useSessionExpired();

            openSessionExpired();

            // 🔥 redirect login
            // 🔥 beri waktu snackbar tampil dulu
            setTimeout(() => {
                // redirect login
                router.push("/login");
            }, 1200);
        }

        return Promise.reject(error);
    }
);

export default api;

// Axios akan menjadi jembatan frontend ↔ backend