import { ref, computed } from "vue";

// 🔥 ambil token dari localStorage saat pertama kali app dibuka
const token = ref(localStorage.getItem("token") || null);

export function useAuth() {

  // 🔥 cek apakah user login
  const isAuthenticated = computed(() => !!token.value);

  // simpan token
  const setToken = (newToken) => {
    token.value = newToken;
    // simpan ke browser
    localStorage.setItem("token", newToken);
  };

  // simpan role
  const setRole = (role) => {
    localStorage.setItem("role", role);
  };

  // ambil role
  const getRole = () => {
    return localStorage.getItem("role");
  };

  // helper
  const isAdmin = () => {
    return getRole() === "admin";
  };

  // 🔥 LOGOUT
  const logout = () => {
    // hapus token dari state
    token.value = null;
    // hapus token dari browser
    localStorage.removeItem("token");
  };

  return {
    token,
    isAuthenticated,
    setToken,
    setRole,
    getRole,
    isAdmin,
    logout,
  };
}