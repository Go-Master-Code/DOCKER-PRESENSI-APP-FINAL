import { useTheme } from "vuetify";

// composable untuk dark mode
export const useAppTheme = () => {

    // akses theme vuetify
    const theme = useTheme();

    // =========================
    // 🔥 CHECK DARK MODE
    // =========================
    // return true jika theme sekarang dark
    const isDark = () => {
        return theme.global.current.value.dark;
    };

    // =========================
    // 🔥 TOGGLE THEME
    // =========================
    const toggleTheme = () => {
        // 🔥 cek theme aktif sekarang
        const newTheme =
            theme.global.current.value.dark
                ? "light"
                : "dark";

        // 🔥 gunakan API terbaru Vuetify
        theme.change(newTheme);

        // 🔥 simpan ke localStorage
        localStorage.setItem(
            "theme",
            newTheme
        );
    };

    // =========================
    // 🔥 LOAD THEME
    // =========================
    const loadTheme = () => {
        // ambil theme tersimpan
        const savedTheme =
            localStorage.getItem("theme");

        // apply jika ada
        if (savedTheme) {
            theme.change(savedTheme);
        }
    };

    return {
        isDark,
        toggleTheme,
        loadTheme,
    };
};