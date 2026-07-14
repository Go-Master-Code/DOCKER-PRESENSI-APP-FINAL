import { ref } from "vue";

// 🔥 STATE GLOBAL (shared di seluruh app)
const show = ref(false);
const text = ref("");
const color = ref("success"); // success | error | warning | info
const timeout = ref(3000);

// 🔥 FUNCTION UTAMA
const openSnackbar = (message, type = "success") => {
    text.value = message;
    color.value = type;
    show.value = true;
};

// 🔥 CLOSE manual (optional)
const closeSnackbar = () => {
    show.value = false;
};

// 🔥 EXPORT sebagai composable
export function useSnackbar() {
    return {
        show,
        text,
        color,
        timeout,
        openSnackbar,
        closeSnackbar,
    };
}